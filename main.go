package main

import (
	"fmt"
	"github.com/Piszmog/eurekaclient/eureka"
	"github.com/Piszmog/httpclient"
	"os"
	"time"
)

func main() {
	hostname, _ := os.Hostname()
	instance := eureka.CreateInstance("testapp", 8080)
	httpClient := httpclient.CreateDefaultHttpClient()
	fmt.Println("registering app")
	client := eureka.CreateLocalClient("http://localhost:8761", "testapp", hostname, httpClient)
	err := client.Register(instance)
	if err != nil {
		panic(err)
	}
	fmt.Println("successful register")
	defer deleteApp(client)
	time.Sleep(2 * time.Second)
	fmt.Println("sending heartbeat")
	err = client.Heartbeat()
	if err != nil {
		panic(err)
	}
	fmt.Println("sent heartbeat")
	time.Sleep(2 * time.Second)
	fmt.Println("retrieving all apps")
	applications, err := client.GetAllApps()
	if err != nil {
		panic(err)
	}
	fmt.Printf("Response %+v\n", applications)
	time.Sleep(2 * time.Second)
	fmt.Println("retrieving instances")
	application, err := client.GetAppInstances("testapp")
	if err != nil {
		panic(err)
	}
	fmt.Printf("Response %+v\n", application)
}

func deleteApp(client eureka.Client) {
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
