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
	GetAllApps() (*Applications, error)
	GetAppInstances(appName string) (*Application, error)
}

type Client struct {
	baseUrl    string
	appName    string
	instanceId string
	httpClient *http.Client
}

func CreateLocalClient(baseUrl, appName, hostname string, httpClient *http.Client) Client {
	return Client{
		baseUrl:    baseUrl,
		appName:    strings.ToUpper(appName),
		instanceId: hostname + ":" + appName,
		httpClient: httpClient,
	}
}

func (client Client) Register(instance Instance) error {
	xmlString, err := xml.Marshal(instance)
	if err != nil {
		return errors.Wrap(err, "failed to convert instance object to xml")
	}
	resp, err := client.httpClient.Post(client.baseUrl+"/eureka/apps/"+client.appName,
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
	request, err := http.NewRequest("PUT",
		client.baseUrl+"/eureka/apps/"+client.appName+"/"+client.instanceId,
		nil)
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
	request, err := http.NewRequest("PUT",
		client.baseUrl+"/eureka/apps/"+client.appName+"/"+client.instanceId+"/status?value="+string(statusType),
		nil)
	if err != nil {
		return errors.Wrapf(err, "failed to create a request to update status to %v", statusType)
	}
	resp, err := client.httpClient.Do(request)
	if err != nil {
		return errors.Wrapf(err, "failed to update status to %v", statusType)
	}
	if resp.StatusCode != 200 {
		return errors.Errorf("failed to update status to %v. Received status %d", statusType, resp.StatusCode)
	}
	return nil
}

func (client Client) CancelInstance() error {
	request, err := http.NewRequest("DELETE", client.baseUrl+"/eureka/apps/"+client.appName+"/"+client.instanceId, nil)
	if err != nil {
		return errors.Wrap(err, "failed to create request to cancel instance")
	}
	resp, err := client.httpClient.Do(request)
	if err != nil {
		return errors.Wrap(err, "failed to cancel instance")
	}
	if resp.StatusCode != 200 {
		return errors.Errorf("failed to cancel instance. Received status %d", resp.StatusCode)
	}
	return nil
}

func (client Client) GetAllApps() (*Applications, error) {
	resp, err := client.httpClient.Get(client.baseUrl + "/eureka/apps")
	if err != nil {
		return nil, errors.Wrap(err, "failed to retrieve all apps")
	}
	if resp.StatusCode != 200 {
		return nil, errors.Errorf("failed to get all apps. Received status %d", resp.StatusCode)
	}
	defer resp.Body.Close()
	applications := &Applications{}
	decoder := xml.NewDecoder(resp.Body)
	err = decoder.Decode(applications)
	if err != nil {
		return nil, errors.Wrap(err, "failed to decode response")
	}
	return applications, nil
}

func (client Client) GetAppInstances(appName string) (*Application, error) {
	resp, err := client.httpClient.Get(client.baseUrl + "/eureka/apps/" + strings.ToUpper(appName))
	if err != nil {
		return nil, errors.Wrap(err, "failed to retrieve app instances")
	}
	if resp.StatusCode != 200 {
		return nil, errors.Errorf("failed to get all apps. Received status %d", resp.StatusCode)
	}
	defer resp.Body.Close()
	application := &Application{}
	decoder := xml.NewDecoder(resp.Body)
	err = decoder.Decode(application)
	if err != nil {
		return nil, errors.Wrap(err, "failed to decode response")
	}
	return application, nil
}
