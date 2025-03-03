package main

import (
	"fmt"
	"log"
	"net"
	"sort"
)

func main() {
	domain := "flamapp.com"
	mxRecords, err := net.LookupMX(domain)
	if err != nil {
		log.Fatal(err)
	}

	sort.Slice(mxRecords, func(i, j int) bool {
		return mxRecords[i].Pref < mxRecords[j].Pref
	})

	fmt.Println("MX records for", domain)
	for _, mx := range mxRecords {
		fmt.Printf("Hostname: %s, Preference: %d\n", mx.Host, mx.Pref)
	}

	if mxRecords[0].Host == "smtp.google.com" {
		fmt.Println("You can also login via Google")
	}

}
