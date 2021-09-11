package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

type IP struct {
	Query string
}

func main() {
	publicIp, err := getPublicIp()
	if err != nil {
		log.Fatal(err.Error())
	}

	err = updateDynuDomain(*publicIp)
	if err != nil {
		log.Fatal(err.Error())
	}

	log.Print("domain updated successfully")
}

func getPublicIp() (*string, error){
	req, err := http.Get("http://ip-api.com/json/")
	if err != nil {
		return nil, err
	}
	defer req.Body.Close()

	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return nil, err
	}

	var ip IP
	json.Unmarshal(body, &ip)

	return &ip.Query, nil
}

func updateDynuDomain(publicIp string) error {
	client := &http.Client{
		Timeout: time.Second * 10,
	}

	domainId := os.Getenv("DOMAIN_ID")
	domainName := os.Getenv("DOMAIN_NAME")
	apiKey := os.Getenv("API_KEY")

	url := fmt.Sprintf("https://api.dynu.com/v2/dns/%s", domainId)

	values := map[string]string{"id": domainId, "name": domainName, "ipv4Address": publicIp}
	json_data, err := json.Marshal(values)

	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(json_data))
	if err != nil {
		return err
	}
	req.Header.Set("API-Key", apiKey)

	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	return nil
}