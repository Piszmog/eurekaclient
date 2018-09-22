package net

import (
	"github.com/Piszmog/cfservices"
	"github.com/Piszmog/cloudconfigclient/net"
	"github.com/pkg/errors"
	"net/http"
)

const (
	defaultServiceRegistryName = "p-service-registry"
)

// CreateCloudClient creates a http client to access Service Registries running in the cloud (specifically Cloud Foundry).
//
// The environment variables 'VCAP_SERVICES' provides a JSON that contains an entry with the key 'p-service-registry'. This
// entry and used to build an OAuth2 client.
func CreateCloudClient() ([]*http.Client, error) {
	return CreateCloudClientForService(defaultServiceRegistryName)
}

// CreateCloudClientForService creates a ConfigClient to access Service Registries running in the cloud (specifically Cloud Foundry).
//
// The environment variables 'VCAP_SERVICES' provides a JSON. The JSON should contain the entry matching the specified name. This
// entry and used to build an OAuth2 client.
func CreateCloudClientForService(name string) ([]*http.Client, error) {
	serviceCredentials, err := GetCloudCredentials(name)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create cloud client")
	}
	return CreateOAuth2Client(serviceCredentials.Credentials)
}

// CreateOAuth2Client creates a http client to access Service Registries from an array of credentials.
func CreateOAuth2Client(creds []cfservices.Credentials) ([]*http.Client, error) {
	httpClients := make([]*http.Client, len(creds))
	for index, cred := range creds {
		uri := cred.Uri
		client, err := net.CreateOAuth2Client(&cred)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to create oauth2 client for %s", uri)
		}
		httpClients[index] = client
	}
	return httpClients, nil
}

// GetCloudCredentials retrieves the Service Registry's credentials so an OAuth2 client can be created.
func GetCloudCredentials(name string) (*cfservices.ServiceCredentials, error) {
	serviceCreds, err := cfservices.GetServiceCredentialsFromEnvironment(name)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get credentials for the Config Server service %s", name)
	}
	return serviceCreds, nil
}
