package eureka

import (
	"encoding/xml"
	"strconv"
	"strings"
)

type Instance struct {
	XMLName          xml.Name       `xml:"instance"`
	HostName         string         `xml:"hostName"`
	Application      string         `xml:"app"`
	InstanceId       string         `xml:"instanceId"`
	IPAddress        string         `xml:"ipAddr"`
	VIPAddress       string         `xml:"vipAddress"`
	SecureVIPAddress string         `xml:"secureVipAddress"`
	Status           StatusType     `xml:"status"`
	OverriddenStatus StatusType     `xml:"overriddenStatus"`
	Port             int            `xml:"port"`
	SecurePort       int            `xml:"securePort"`
	HomePageURL      string         `xml:"homePageUrl"`
	StatusPageURL    string         `xml:"statusPageUrl"`
	HealthCheckURL   string         `xml:"healthCheckUrl"`
	DataCenterInfo   DataCenterInfo `xml:"dataCenterInfo"`
	LeaseInfo        LeaseInfo      `xml:"leaseInfo"`
}

type DataCenterInfo struct {
	Name     DCNameType     `xml:"name"`
	Metadata AmazonMetadata `xml:"metadata"`
}

type DCNameType string

const (
	MyOwn  DCNameType = "MyOwn"
	Amazon DCNameType = "Amazon"
)

type StatusType string

const (
	Up           StatusType = "UP"
	Down         StatusType = "DOWN"
	Starting     StatusType = "STARTING"
	OutOfService StatusType = "OUT_OF_SERVICE"
	Unknown      StatusType = "UNKNOWN"
)

type AmazonMetadata struct {
	AMILaunchIndex   string `xml:"ami-launch-index"`
	LocalHostName    string `xml:"local-hostname"`
	AvailabilityZone string `xml:"availability-zone"`
	InstanceId       string `xml:"instance-id"`
	PublicIPV4       string `xml:"public-ipv4"`
	PublicHostName   string `xml:"public-hostname"`
	AMIManifestPath  string `xml:"ami-manifest-path"`
	LocalIPV4        string `xml:"local-ipv4"`
	HostName         string `xml:"hostname"`
	AMIId            string `xml:"ami-id"`
	InstanceType     string `xml:"instance-type"`
}

type LeaseInfo struct {
	DurationInSecs        int `xml:"durationInSecs"`
	RenewalIntervalInSecs int `xml:"renewalIntervalInSecs"`
	RegistrationTimestamp int `xml:"registrationTimestamp"`
	LastRenewalTimestamp  int `xml:"lastRenewalTimestamp"`
	EvictionTimestamp     int `xml:"evictionTimestamp"`
	ServiceUpTimestamp    int `xml:"serviceUpTimestamp"`
}

func CreateInstance(appName, hostname, instanceId string, port int) Instance {
	return Instance{
		HostName:         hostname,
		Port:             port,
		Status:           Up,
		IPAddress:        "127.0.0.1",
		VIPAddress:       appName,
		SecureVIPAddress: appName,
		Application:      strings.ToUpper(appName),
		InstanceId:       instanceId,
		HomePageURL:      "http://" + hostname + ":" + strconv.Itoa(8080),
		DataCenterInfo: DataCenterInfo{
			Name: MyOwn,
		},
		LeaseInfo: LeaseInfo{
			RenewalIntervalInSecs: 30,
			DurationInSecs:        90,
		},
	}
}
