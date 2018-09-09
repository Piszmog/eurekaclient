package eureka

import "encoding/xml"

type Applications struct {
	XMLName              xml.Name    `xml:"applications"`
	VersionsDelta        int         `xml:"versions__delta"`
	ApplicationsHashCode string      `xml:"apps__hashcode"`
	Application          Application `xml:"application"`
}

type Application struct {
	XMLName  xml.Name `xml:"application"`
	Name     string   `xml:"name"`
	Instance Instance `xml:"instance"`
}
