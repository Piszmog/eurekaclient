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
	instance := eureka.RegistryInstance{
		AppName:    "demo-client",
		Port:       8080,
		SecurePort: 443,
	}
	appInstance, err := service.Register("http://localhost:8761/eureka/apps", instance, httpClient)
	eurekaClient := client.CreateLocalClient("http://localhost:8761", httpClient)
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
	time.Sleep(30 * time.Second)
	fmt.Println("retrieving all apps")
	applications, err := eurekaClient.GetAllApps()
	if err != nil {
		panic(err)
	}
	fmt.Printf("Response %+v\n", applications)
	time.Sleep(2 * time.Second)
	fmt.Println("retrieving instances")
	application, err := eurekaClient.GetAppInstances("demo-client")
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
