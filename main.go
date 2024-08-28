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

// global variables
var (
	ddns_apiusername, ddns_apiToken, ddns_domain, ddns_id, ddns_url string
)

// initializing global variables from .env
func init() {

	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file:", err)
		return
	}

	ddns_domain = os.Getenv("DOMAIN")
	ddns_id = os.Getenv("ID")
	ddns_apiusername = os.Getenv("APIUSERNAME")
	ddns_apiToken = os.Getenv("APITOKEN")

	ddns_url = fmt.Sprintf("https://api.name.com/v4/domains/%s/records/%s", ddns_domain, ddns_id)
}

// json struct for domain records (name.com)
type DNSRecord struct {
	ID         int    `json:"id"`
	DomainName string `json:"domainName"`
	Host       string `json:"host"`
	FQDN       string `json:"fqdn"`
	Type       string `json:"type"`
	Answer     string `json:"answer"`
	TTL        int    `json:"ttl"`
}

// json struct for ip records (aws.mugund10.top/ip)
type Ip struct {
	Addr string `json:"addr"`
}

//gets ip from url (mostly url = "aws.mugund10.top/ip" ) 
func Getip(url string) string {
	var msg Ip
	resp, err := http.Get(url)
	if err != nil {
		log.Fatalln(err)
	}
	defer resp.Body.Close()

	body := bodyreader(resp.Body)
	err = json.Unmarshal(body, &msg)
	if err != nil {
		fmt.Println(err)
	}

	return msg.Addr
}

// reads http response body
func bodyreader(reader io.Reader) []byte {
	body, err := io.ReadAll(reader)
	if err != nil {
		log.Fatalln(err)
	}
	return body
}

// reads dns value of 'A' record of "homeserver" subdomain
func Getdns(record *DNSRecord) {

	req, err := http.NewRequest("GET", ddns_url, nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}

	ReqHeaders(req, ddns_apiusername, ddns_apiToken)

	resp, err := Request(req)
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

}

// update ip to A record of "homeserver" subdomain
func Putdns(ip string) {

	jsonPayload := []byte(fmt.Sprintf(`{
		"host": "homeserver",
		"type": "A",
		"answer": "%s",
		"ttl": 600
	}`, ip))

	req, err := http.NewRequest("PUT", ddns_url, bytes.NewBuffer(jsonPayload))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}

	ReqHeaders(req, ddns_apiusername, ddns_apiToken)

	resp, err := Request(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
	}
	defer resp.Body.Close()

	if resp.Status != "200 OK" {
		log.Println(resp.Status)
	}

}

// makes http request
func Request(req *http.Request) (*http.Response, error) {
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// sets required Headers
func ReqHeaders(req *http.Request, apiUsername, apiToken string) {
	req.SetBasicAuth(apiUsername, apiToken)
	req.Header.Set("Content-Type", "application/json")
}

func main() {

	record := DNSRecord{}

	for {
		//using a service which returns ip
		ip := Getip("https://aws.mugund10.top/ip")
		Getdns(&record)
		if ip == record.Answer {
			// Do nothing
			log.Println("[ current ip : ", ip, "]", "[ old ip : ", record.Answer, "]", " -> DNS : not changed")
		} else {
			log.Println("[ current ip : ", ip, "]", "[ old ip : ", record.Answer, "]", " -> DNS : changed")
			Putdns(ip)
		}
		time.Sleep(time.Minute * 5)
	}

}
