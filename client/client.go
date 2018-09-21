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
	BaseUrl       string
	EurekaAPIPath string
	HttpClient    *http.Client
}

func (client EurekaClient) GetAllApps() (*eureka.Applications, error) {
	resp, err := client.HttpClient.Get(client.BaseUrl + client.EurekaAPIPath)
	if err != nil {
		return nil, errors.Wrap(err, "failed to retrieve all apps from eureka")
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
	resp, err := client.HttpClient.Get(client.BaseUrl + client.EurekaAPIPath + "/" + strings.ToUpper(appName))
	if err != nil {
		return nil, errors.Wrapf(err, "failed to retrieve instances of application %s", appName)
	}
	if resp.StatusCode != 200 {
		return nil, errors.Errorf("failed to get all instances of %s. Received status %d", appName, resp.StatusCode)
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
