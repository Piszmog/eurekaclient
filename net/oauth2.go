package net

import (
	"github.com/Piszmog/cfservices"
	"github.com/Piszmog/httpclient"
	"github.com/pkg/errors"
	"net/http"
)

const (
	defaultServiceRegistryName = "p-service-registry"
)

// CreateCloudClient creates a ConfigClient to access Config Servers running in the cloud (specifically Cloud Foundry).
//
// The environment variables 'VCAP_SERVICES' provides a JSON that contains an entry with the key 'p-config-server'. This
// entry and used to build an OAuth2 client.
func CreateCloudClient() (*http.Client, error) {
	return CreateCloudClientForService(defaultServiceRegistryName)
}

// CreateCloudClientForService creates a ConfigClient to access Config Servers running in the cloud (specifically Cloud Foundry).
//
// The environment variables 'VCAP_SERVICES' provides a JSON. The JSON should contain the entry matching the specified name. This
// entry and used to build an OAuth2 client.
func CreateCloudClientForService(name string) (*http.Client, error) {
	serviceCredentials, err := GetCloudCredentials(name)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create cloud client")
	}
	creds := serviceCredentials.Credentials
	if len(creds) == 0 {
		return nil, errors.New("failed to obtain credentials for the service registry")
	}
	return CreateOAuth2Client(creds[0])
}

// CreateOAuth2Client creates a ConfigClient to access Config Servers from an array of credentials.
func CreateOAuth2Client(cred cfservices.Credentials) (*http.Client, error) {
	configUri := cred.Uri
	client, err := httpclient.CreateOAuth2Client(cred.ClientId, cred.ClientSecret, cred.AccessTokenUri)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to create oauth2 client for %s", configUri)
	}
	return client, nil
}

// GetCloudCredentials retrieves the Config Server's credentials so an OAuth2 client can be created.
func GetCloudCredentials(name string) (*cfservices.ServiceCredentials, error) {
	serviceCreds, err := cfservices.GetServiceCredentialsFromEnvironment(name)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get credentials for the Config Server service %s", name)
	}
	return serviceCreds, nil
}
