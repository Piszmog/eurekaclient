package eureka

import (
	"encoding/xml"
	"github.com/pkg/errors"
	"net/http"
	"strings"
)

type Eureka interface {
	Register(instance Instance) error
	Heartbeat() error
	UpdateStatus(statusType StatusType) error
	CancelInstance() error
	GetAllApps() error
	GetAppInstances(appName string) error
}

type Client struct {
	baseUrl    string
	appName    string
	instanceId string
	httpClient *http.Client
}

func (client Client) Register(instance Instance) error {
	xmlString, err := xml.Marshal(instance)
	if err != nil {
		return errors.Wrap(err, "failed to convert instance object to xml")
	}
	resp, err := client.httpClient.Post(client.baseUrl+"/eureka/apps/"+strings.ToUpper(client.appName),
		"application/xml",
		strings.NewReader(string(xmlString)))
	if err != nil {
		return errors.Wrap(err, "failed to register instance with eureka")
	}
	if resp.StatusCode != 204 {
		return errors.Errorf("failed to register instance with eureka. Status code %d", resp.StatusCode)
	}
	return nil
}

func (client Client) Heartbeat() error {
	request, err := http.NewRequest("PUT", client.baseUrl+"/eureka/apps/"+strings.ToUpper(client.appName)+"/"+client.instanceId, nil)
	if err != nil {
		return errors.Wrap(err, "failed to create heartbeat request")
	}
	resp, err := client.httpClient.Do(request)
	if err != nil {
		return errors.Wrap(err, "failed to send heartbeat to eureka")
	}
	if resp.StatusCode != 200 {
		return errors.Errorf("failed to send heartbeat. Status code: %d", resp.StatusCode)
	}
	return nil
}

func (client Client) UpdateStatus(statusType StatusType) error {

}

func (client Client) CancelInstance() error {

}

func (client Client) GetAllApps() error {

}

func (client Client) GetAppInstances(appName string) error {

}
