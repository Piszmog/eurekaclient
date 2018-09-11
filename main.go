package main

import (
	"fmt"
	"github.com/Piszmog/eurekaclient/client"
	"github.com/Piszmog/eurekaclient/eureka"
	"github.com/Piszmog/eurekaclient/service"
	"github.com/Piszmog/httpclient"
	"time"
)

func main() {
	httpClient := httpclient.CreateDefaultHttpClient()
	fmt.Println("registering app")
	appInstance, _ := service.CreateApplicationInstance("http://localhost:8761", "testapp", "", 8080, httpClient)
	eurekaClient := client.CreateLocalClient("http://localhost:8761", httpClient)
	err := appInstance.Register()
	if err != nil {
		panic(err)
	}
	fmt.Println("successful register")
	defer deleteApp(appInstance)
	time.Sleep(3 * time.Second)
	fmt.Println("sending heartbeat")
	//err = appInstance.Heartbeat()
	appInstance.StartHeartbeats(10)
	if err != nil {
		panic(err)
	}
	fmt.Println("sent heartbeat")
	time.Sleep(90 * time.Second)
	fmt.Println("retrieving all apps")
	applications, err := eurekaClient.GetAllApps()
	if err != nil {
		panic(err)
	}
	fmt.Printf("Response %+v\n", applications)
	time.Sleep(2 * time.Second)
	fmt.Println("retrieving instances")
	application, err := eurekaClient.GetAppInstances("testapp")
	if err != nil {
		panic(err)
	}
	fmt.Printf("Response %+v\n", application)
}

func deleteApp(client service.EurekaInstance) {
	time.Sleep(2 * time.Second)
	fmt.Println("updating status to DOWN")
	err := client.UpdateStatus(eureka.Down)
	if err != nil {
		panic(err)
	}
	fmt.Println("updated status")
	time.Sleep(10 * time.Second)
	fmt.Println("canceling instance")
	err = client.CancelInstance()
	if err != nil {
		panic(err)
	}
	fmt.Println("cancelled instance")
}
