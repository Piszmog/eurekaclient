package service

import (
	"encoding/xml"
	"github.com/Piszmog/eurekaclient/eureka"
	"github.com/pkg/errors"
	"net/http"
	"os"
	"strings"
	"time"
)

const (
	Put       = "PUT"
	Delete    = "DELETE"
	XmlHeader = "application/xml"
)

type EurekaInstance interface {
	Register() error
	Heartbeat() error
	UpdateStatus(statusType eureka.StatusType) error
	CancelInstance() error
}

type ApplicationInstance struct {
	baseUrl    string
	appName    string
	port       int
	instanceId string
	httpClient *http.Client
}

func CreateApplicationInstance(baseUrl, appName, urlPath string, port int, httpClient *http.Client) (*ApplicationInstance, error) {
	hostname, err := os.Hostname()
	if err != nil {
		return nil, errors.Wrap(err, "failed to retrieve hostname of application")
	}
	eurekaURLPath := urlPath
	if len(urlPath) == 0 {
		eurekaURLPath = "/eureka/apps"
	}
	return &ApplicationInstance{
		baseUrl:    baseUrl + eurekaURLPath,
		appName:    strings.ToUpper(appName),
		port:       port,
		instanceId: hostname + ":" + appName,
		httpClient: httpClient,
	}, nil
}

func (eurekaInstance ApplicationInstance) Register() error {
	instance, err := eureka.CreateInstance(eurekaInstance.appName, getHostname(eurekaInstance.instanceId), eurekaInstance.instanceId, eurekaInstance.port)
	if err != nil {
		return errors.Wrap(err, "failed to create an Eureka instance")
	}
	xmlString, err := xml.Marshal(instance)
	if err != nil {
		return errors.Wrap(err, "failed to convert instance object to xml")
	}
	resp, err := eurekaInstance.httpClient.Post(eurekaInstance.baseUrl+"/"+eurekaInstance.appName,
		XmlHeader,
		strings.NewReader(string(xmlString)))
	if err != nil {
		return errors.Wrap(err, "failed to register instance with eureka")
	}
	if resp.StatusCode != 204 {
		return errors.Errorf("failed to register instance with eureka. Status code %d", resp.StatusCode)
	}
	return nil
}

func (eurekaInstance ApplicationInstance) Heartbeat() error {
	request, err := http.NewRequest(Put,
		eurekaInstance.baseUrl+"/"+eurekaInstance.appName+"/"+eurekaInstance.instanceId,
		nil)
	if err != nil {
		return errors.Wrap(err, "failed to create heartbeat request")
	}
	resp, err := eurekaInstance.httpClient.Do(request)
	if err != nil {
		return errors.Wrap(err, "failed to send heartbeat to eureka")
	}
	if resp.StatusCode != 200 {
		return errors.Errorf("failed to send heartbeat. Status code: %d", resp.StatusCode)
	}
	return nil
}

func (eurekaInstance ApplicationInstance) UpdateStatus(statusType eureka.StatusType) error {
	request, err := http.NewRequest(Put,
		eurekaInstance.baseUrl+"/"+eurekaInstance.appName+"/"+eurekaInstance.instanceId+"/status?value="+string(statusType),
		nil)
	if err != nil {
		return errors.Wrapf(err, "failed to create a request to update status to %v", statusType)
	}
	resp, err := eurekaInstance.httpClient.Do(request)
	if err != nil {
		return errors.Wrapf(err, "failed to update status to %v", statusType)
	}
	if resp.StatusCode != 200 {
		return errors.Errorf("failed to update status to %v. Received status %d", statusType, resp.StatusCode)
	}
	return nil
}

func (eurekaInstance ApplicationInstance) CancelInstance() error {
	request, err := http.NewRequest(Delete,
		eurekaInstance.baseUrl+"/"+eurekaInstance.appName+"/"+eurekaInstance.instanceId,
		nil)
	if err != nil {
		return errors.Wrap(err, "failed to create request to cancel instance")
	}
	resp, err := eurekaInstance.httpClient.Do(request)
	if err != nil {
		return errors.Wrap(err, "failed to cancel instance")
	}
	if resp.StatusCode != 200 {
		return errors.Errorf("failed to cancel  Received status %d", resp.StatusCode)
	}
	return nil
}

func (eurekaInstance ApplicationInstance) StartHeartbeats(intervalInSeconds time.Duration) {
	go func() {
		for {
			eurekaInstance.Heartbeat() //todo do something with the error from the heartbeat
			<-time.After(intervalInSeconds * time.Second)
		}
	}()
}

func getHostname(instanceId string) string {
	return instanceId[:strings.IndexByte(instanceId, ':')]
}
