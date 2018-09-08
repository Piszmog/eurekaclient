package main

import (
	"encoding/xml"
	"fmt"
	"github.com/Piszmog/eurekaclient/eureka"
	"github.com/Piszmog/httpclient"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"
)

func main() {
	hostname, _ := os.Hostname()
	instance := eureka.CreateInstance("testapp", 8080)
	client := httpclient.CreateDefaultHttpClient()
	const baseUrl = "http://localhost:8761"
	request, _ := http.NewRequest("PUT", baseUrl+"/eureka/apps/TESTAPP/"+hostname+":testapp", nil)
	resp, _ := client.Do(request)
	if resp.StatusCode != 404 {
		panic("sent a successful heartbeat")
	}
	xmlString, _ := xml.Marshal(instance)
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
