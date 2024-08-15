package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
)

var (
	apiUsername, apiToken, domain, id string
)

func init() {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file:", err)
		return
	}

	// Initialize global variables with environment variables
	domain = os.Getenv("DOMAIN")
	id = os.Getenv("ID")
	apiUsername = os.Getenv("APIUSERNAME")
	apiToken = os.Getenv("APITOKEN")
}

type DNSRecord struct {
	ID         int    `json:"id"`
	DomainName string `json:"domainName"`
	Host       string `json:"host"`
	FQDN       string `json:"fqdn"`
	Type       string `json:"type"`
	Answer     string `json:"answer"`
	TTL        int    `json:"ttl"`
}

var url string

func init() {
	url = fmt.Sprintf("https://api.name.com/v4/domains/%s/records/%s", domain, id)
}

func Getip(url string) string {
	resp, err := http.Get(url)
	if err != nil {
		log.Fatalln(err)
	}
	defer resp.Body.Close()

	body := bodyreader(resp.Body)
	return string(body[8 : len(body)-2])
}

func bodyreader(reader io.Reader) []byte {
	body, err := io.ReadAll(reader)
	if err != nil {
		log.Fatalln(err)
	}
	return body
}

func Getdns(record *DNSRecord) {

	client := &http.Client{}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}

	req.SetBasicAuth(apiUsername, apiToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error making request:", err)
		return
	}
	defer resp.Body.Close()

	body := bodyreader(resp.Body)

	err = json.Unmarshal(body, &record)
	if err != nil {
		fmt.Println("Error parsing JSON:", err)
		return
	}
	// fmt.Println("Response Body:", string(body))
	// fmt.Printf("DNS Record: %+v\n", record)

}

func Putdns(ip string) {

	// JSON payload as a byte slice
	jsonPayload := []byte(fmt.Sprintf(`{
		"host": "homeserver",
		"type": "A",
		"answer": "%s",
		"ttl": 600
	}`, ip))

	// Create a new request using http.NewRequest
	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(jsonPayload))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}

	// Set headers
	req.SetBasicAuth(apiUsername, apiToken)
	req.Header.Set("Content-Type", "application/json")

	// Create a new HTTP client and send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return
	}
	defer resp.Body.Close()

	// Read and print the response
	//body := bodyreader(resp.Body)

	// fmt.Println("Response Status:", resp.Status)
	// fmt.Println("Response Body:", string(body))

	if resp.Status != "200 OK" {
		log.Println(resp.Status)
	}

}

func main() {

	ip := Getip("https://mugund10.openwaves.in/ip")
	record := DNSRecord{}

	for {

		Getdns(&record)
		if ip == record.Answer {
			// Do nothing
			log.Println("current ip = ", ip, "old ip", record.Answer, " //not changed")
		} else {
			log.Println("current ip = ", ip, "old ip", record.Answer, " //changed")
			Putdns(ip)
		}
		time.Sleep(time.Minute * 5)
	}

}
