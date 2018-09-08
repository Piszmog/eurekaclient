package main

import (
	"encoding/xml"
	"fmt"
	"github.com/Piszmog/httpclient"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
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

func main() {
	hostname, _ := os.Hostname()
	instance := Instance{
		HostName:         hostname,
		Port:             8080,
		Status:           Up,
		IPAddress:        "127.0.0.1",
		VIPAddress:       "testapp",
		SecureVIPAddress: "testapp",
		Application:      "TESTAPP",
		InstanceId:       hostname + ":testapp",
		HomePageURL:      "http://" + hostname + ":" + strconv.Itoa(8080),
		DataCenterInfo: DataCenterInfo{
			Name: MyOwn,
		},
		LeaseInfo: LeaseInfo{
			RenewalIntervalInSecs: 30,
			DurationInSecs:        90,
		},
	}
	xmlString, _ := xml.Marshal(instance)
	client := httpclient.CreateDefaultHttpClient()
	const baseUrl = "http://localhost:8761"
	request, _ := http.NewRequest("PUT", baseUrl+"/eureka/apps/TESTAPP/"+hostname+":testapp", nil)
	resp, _ := client.Do(request)
	if resp.StatusCode != 404 {
		panic("sent a successful heartbeat")
	}
	s := string(xmlString)
	if xmlString == nil {
		panic("fail")
	}
	resp, _ = client.Post(baseUrl+"/eureka/apps/TESTAPP", "application/xml", strings.NewReader(s))
	if resp.StatusCode != 204 {
		panic("failed to get 204")
	}
	defer deleteApp(nil, baseUrl, resp, client, hostname)
	time.Sleep(2 * time.Second)
	request, _ = http.NewRequest("PUT", baseUrl+"/eureka/apps/TESTAPP/"+hostname+":testapp", nil)
	resp, _ = client.Do(request)
	if resp.StatusCode != 200 {
		panic("failed to send a successful heartbeat")
	}
	resp, _ = client.Get(baseUrl + "/eureka/apps")
	if resp.StatusCode != 200 {
		panic("failed to get 200")
	}
	bytes, _ := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	fmt.Printf("Response %s\n", string(bytes))

	resp, _ = client.Get(baseUrl + "/eureka/apps/TESTAPP")
	if resp.StatusCode != 200 {
		panic("failed to get 200")
	}
	bytes, _ = ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	fmt.Printf("Response %s\n", string(bytes))
}

func deleteApp(request *http.Request, baseUrl string, resp *http.Response, client *http.Client, hostname string) {
	time.Sleep(2 * time.Second)
	request, _ = http.NewRequest("PUT", baseUrl+"/eureka/apps/TESTAPP/"+hostname+":testapp/status?value=DOWN", nil)
	resp, _ = client.Do(request)
	if resp.StatusCode != 200 {
		panic("failed to update the status to DOWN")
	}
	time.Sleep(10 * time.Second)
	request, _ = http.NewRequest("DELETE", baseUrl+"/eureka/apps/TESTAPP/"+hostname+":testapp", nil)
	resp, _ = client.Do(request)
	if resp.StatusCode != 200 {
		panic("failed to get 200, got " + resp.Status)
	}
}
