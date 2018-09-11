package client

import (
	"encoding/xml"
	"github.com/Piszmog/eurekaclient/eureka"
	"github.com/pkg/errors"
	"net/http"
	"strings"
)

type Eureka interface {
	GetAllApps() (*eureka.Applications, error)
	GetAppInstances(appName string) (*eureka.Application, error)
}

type EurekaClient struct {
	baseUrl    string
	httpClient *http.Client
}

func CreateLocalClient(baseUrl string, httpClient *http.Client) *EurekaClient {
	return &EurekaClient{
		baseUrl:    baseUrl,
		httpClient: httpClient,
	}
}

func (client EurekaClient) GetAllApps() (*eureka.Applications, error) {
	resp, err := client.httpClient.Get(client.baseUrl + "/eureka/apps")
	if err != nil {
		return nil, errors.Wrap(err, "failed to retrieve all apps")
	}
	if resp.StatusCode != 200 {
		return nil, errors.Errorf("failed to get all apps. Received status %d", resp.StatusCode)
	}
	defer resp.Body.Close()
	applications := &eureka.Applications{}
	decoder := xml.NewDecoder(resp.Body)
	err = decoder.Decode(applications)
	if err != nil {
		return nil, errors.Wrap(err, "failed to decode response")
	}
	return applications, nil
}

func (client EurekaClient) GetAppInstances(appName string) (*eureka.Application, error) {
	resp, err := client.httpClient.Get(client.baseUrl + "/eureka/apps/" + strings.ToUpper(appName))
	if err != nil {
		return nil, errors.Wrap(err, "failed to retrieve app instances")
	}
	if resp.StatusCode != 200 {
		return nil, errors.Errorf("failed to get all apps. Received status %d", resp.StatusCode)
	}
	defer resp.Body.Close()
	application := &eureka.Application{}
	decoder := xml.NewDecoder(resp.Body)
	err = decoder.Decode(application)
	if err != nil {
		return nil, errors.Wrap(err, "failed to decode response")
	}
	return application, nil
}
