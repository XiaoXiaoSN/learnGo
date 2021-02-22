package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"

	"github.com/oschwald/geoip2-golang"
)

func main() {
	// db, err := geoip2.Open("GeoLite2-ASN.mmdb")
	// db, err := geoip2.Open("GeoLite2-Country.mmdb")
	db, err := geoip2.Open("GeoLite2-City.mmdb")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	// If you are using strings that may be invalid, check that ip is not nil
	//ip := net.ParseIP("81.2.69.142")
	ip := net.ParseIP("211.22.98.180")
	record, err := db.City(ip)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(prettyPrint(record))

	fmt.Printf("Portuguese (BR) city name: %v\n", record.City.Names["pt-BR"])
	if len(record.Subdivisions) > 0 {
		fmt.Printf("English subdivision name: %v\n", record.Subdivisions[0].Names["en"])
		fmt.Printf("Chinese subdivision name: %v\n", record.Subdivisions[0].Names["zh-CN"])
	}
	fmt.Printf("Russian country name: %v\n", record.Country.Names["ru"])
	fmt.Printf("English country name: %v\n", record.Country.Names["en"])
	fmt.Printf("Chinese country name: %v\n", record.Country.Names["zh-CN"])
	fmt.Printf("ISO country code: %v\n", record.Country.IsoCode)
	fmt.Printf("Time zone: %v\n", record.Location.TimeZone)
	fmt.Printf("Coordinates: %v, %v\n", record.Location.Latitude, record.Location.Longitude)
	// Output:
	// Portuguese (BR) city name: Londres
	// English subdivision name: England
	// Russian country name: Великобритания
	// ISO country code: GB
	// Time zone: Europe/London
	// Coordinates: 51.5142, -0.0931
}

func prettyPrint(data interface{}) string {
	jsonByte, err := json.MarshalIndent(data, "", "    ")
	if err != nil {
		return err.Error()

	}
	return fmt.Sprintf("%s\n", jsonByte)
}
