package service

import (
	"encoding/xml"
	"fmt"
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
	Heartbeat() error
	UpdateStatus(statusType eureka.StatusType) error
	CancelInstance() error
}

type ApplicationInstance struct {
	baseUrl    string
	appName    string
	instanceId string
	httpClient *http.Client
}

func Register(eurekaURL string, registryInstance eureka.RegistryInstance, httpClient *http.Client) (*ApplicationInstance, error) {
	hostname, err := os.Hostname()
	if err != nil {
		return nil, errors.Wrap(err, "failed to retrieve hostname of application")
	}
	appName := strings.ToLower(registryInstance.AppName)
	instanceId := hostname + ":" + appName
	instance, err := eureka.CreateInstance(instanceId, registryInstance)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create an instance to register with Eureka")
	}
	xmlString, err := xml.Marshal(instance)
	if err != nil {
		return nil, errors.Wrap(err, "failed to convert instance object to xml")
	}
	resp, err := httpClient.Post(eurekaURL+"/"+appName, XmlHeader, strings.NewReader(string(xmlString)))
	if err != nil {
		return nil, errors.Wrap(err, "failed to register instance with eureka")
	}
	if resp.StatusCode != 204 {
		return nil, errors.Errorf("failed to register instance with eureka. Status code %d", resp.StatusCode)
	}
	return &ApplicationInstance{
		baseUrl:    eurekaURL,
		appName:    strings.ToUpper(appName),
		instanceId: instanceId,
		httpClient: httpClient,
	}, nil
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
			fmt.Printf("Sent heart beat at %s\n", time.Now())
			err := eurekaInstance.Heartbeat() //todo do something with the error from the heartbeat
			if err != nil {
				fmt.Println("failed to send a heartbeat")
			}
			<-time.After(intervalInSeconds * time.Second)
		}
	}()
}
