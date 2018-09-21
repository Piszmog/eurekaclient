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
	Heartbeat() error
	UpdateStatus(statusType eureka.StatusType) error
	CancelInstance() error
}

type ApplicationInstance struct {
	eurkeaURL  string
	appName    string
	instanceId string
	httpClient *http.Client
}

func Register(baseURL, eurekaAPIPath string, registryInstance eureka.RegistryInstance, httpClient *http.Client) (*ApplicationInstance, error) {
	hostname, err := os.Hostname()
	if err != nil {
		return nil, errors.Wrapf(err, "failed to retrieve hostname of application %s", registryInstance.AppName)
	}
	appName := strings.ToLower(registryInstance.AppName)
	instanceId := hostname + ":" + appName
	instance, err := eureka.CreateInstance(instanceId, registryInstance)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to create an instance of %s to register with Eureka", registryInstance.AppName)
	}
	xmlString, err := xml.Marshal(instance)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to convert instance object to xml. Instance: %+v", instance)
	}
	fullURL := baseURL + eurekaAPIPath
	resp, err := httpClient.Post(fullURL+"/"+appName, XmlHeader, strings.NewReader(string(xmlString)))
	if err != nil {
		return nil, errors.Wrapf(err, "failed to register an instance of %s with eureka", registryInstance.AppName)
	}
	if resp.StatusCode != 204 {
		return nil, errors.Errorf("failed to register an instance of %s with eureka. Status code %d", registryInstance.AppName, resp.StatusCode)
	}
	return &ApplicationInstance{
		eurkeaURL:  fullURL,
		appName:    strings.ToUpper(appName),
		instanceId: instanceId,
		httpClient: httpClient,
	}, nil
}

func (eurekaInstance ApplicationInstance) Heartbeat() error {
	request, err := http.NewRequest(Put,
		eurekaInstance.eurkeaURL+"/"+eurekaInstance.appName+"/"+eurekaInstance.instanceId,
		nil)
	if err != nil {
		return errors.Wrapf(err, "failed to create heartbeat request for application %s", eurekaInstance.instanceId)
	}
	resp, err := eurekaInstance.httpClient.Do(request)
	if err != nil {
		return errors.Wrapf(err, "failed to send heartbeat for %s to eureka", eurekaInstance.instanceId)
	}
	if resp.StatusCode != 200 {
		return errors.Errorf("failed to send heartbeat for %s. Status code: %d", eurekaInstance.instanceId, resp.StatusCode)
	}
	return nil
}

func (eurekaInstance ApplicationInstance) UpdateStatus(statusType eureka.StatusType) error {
	request, err := http.NewRequest(Put,
		eurekaInstance.eurkeaURL+"/"+eurekaInstance.appName+"/"+eurekaInstance.instanceId+"/status?value="+string(statusType),
		nil)
	if err != nil {
		return errors.Wrapf(err, "failed to create a request to update status of %s to %v", eurekaInstance.instanceId, statusType)
	}
	resp, err := eurekaInstance.httpClient.Do(request)
	if err != nil {
		return errors.Wrapf(err, "failed to update status if to %v", eurekaInstance.instanceId, statusType)
	}
	if resp.StatusCode != 200 {
		return errors.Errorf("failed to update status of %s to %v. Received status %d", eurekaInstance.instanceId, statusType, resp.StatusCode)
	}
	return nil
}

func (eurekaInstance ApplicationInstance) CancelInstance() error {
	request, err := http.NewRequest(Delete,
		eurekaInstance.eurkeaURL+"/"+eurekaInstance.appName+"/"+eurekaInstance.instanceId,
		nil)
	if err != nil {
		return errors.Wrapf(err, "failed to create request to cancel instance %s", eurekaInstance.instanceId)
	}
	resp, err := eurekaInstance.httpClient.Do(request)
	if err != nil {
		return errors.Wrapf(err, "failed to cancel instance %s", eurekaInstance.instanceId)
	}
	if resp.StatusCode != 200 {
		return errors.Errorf("failed to cancel %s. Received status %d", eurekaInstance.instanceId, resp.StatusCode)
	}
	return nil
}

func (eurekaInstance ApplicationInstance) StartHeartbeats(intervalInSeconds time.Duration, errs chan error) {
	go func() {
		for {
			err := eurekaInstance.Heartbeat()
			if err != nil {
				errs <- errors.Wrapf(err, "failed to send a heartbeat for instance %s", eurekaInstance.instanceId)
			}
			<-time.After(intervalInSeconds * time.Second)
		}
	}()
}
