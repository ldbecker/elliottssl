package main

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type ApiEndpointsResponse struct {
	IP                string `json:"ipAddress"`
	StatusMsg         string `json:"statusMessage"`
	Grade             string `json:"grade"`
	GradeTrustIgnored string `json:"gradeTrustIgnored"`
	HasWarnings       bool   `json:"hasWarnings"`
	Progress          int    `json:"progress"`
}

type ApiResponse struct {
	Hostname  string                 `json:"host"`
	StatusMsg string                 `json:"status"`
	StartTime int64                  `json:"startTime,omitempty"`
	TestTime  int64                  `json:"testTime,omitempty"`
	Endpoints []ApiEndpointsResponse `json:"endpoints,omitempty"`
}

var APIENDPOINT = "https://api.ssllabs.com/api/v2/"
var INFOPATH = APIENDPOINT + "info"
var ANALYZEPATHNEW = APIENDPOINT + "analyze?startNew=on&host="
var ANALYZEPATH = APIENDPOINT + "analyze?host="
var httpClient *http.Client

func analyzeEndpoint(hostname string) ([]ApiEndpointsResponse, error) {
	results := make([]ApiEndpointsResponse, 0)
	newResp, err := http.Get(ANALYZEPATHNEW + hostname)
	if err != nil {
		return results, fmt.Errorf("error initiating scan: %v", err)
	}
	//defer newResp.Body.Close()
	newRespBytes, err := io.ReadAll(newResp.Body)
	if err != nil {
		return results, fmt.Errorf("error parsing init scan response: %v", err)
	}

	var apiResp ApiResponse
	fmt.Println(string(newRespBytes))
	err = json.Unmarshal(newRespBytes, &apiResp)
	if err != nil {
		return results, fmt.Errorf("error unmarshaling init scan response: %v", err)
	}
	if newResp.StatusCode != http.StatusOK {
		return results, fmt.Errorf("init scan failed with code %v", newResp.StatusCode)
	}

	for apiResp.StatusMsg != "READY" {
		fmt.Printf("Status: %v\n", apiResp.StatusMsg)
		fmt.Printf("Sleeping 5 seconds...\n")
		time.Sleep(time.Second * time.Duration(5))
		newResp, err = http.Get(ANALYZEPATH + hostname)
		if err != nil {
			return results, fmt.Errorf("scan update failed: %v", err)
		}
		//defer newResp.Body.Close()
		newRespBytes, err = io.ReadAll(newResp.Body)
		if err != nil {
			return results, fmt.Errorf("error reading update body: %v", err)
		}
		err = json.Unmarshal(newRespBytes, &apiResp)
		if err != nil {
			return results, fmt.Errorf("error unmarshaling update body: %v", err)

		}
	}

	return apiResp.Endpoints, nil

}

func main() {
	transport := &http.Transport{
		TLSClientConfig:   &tls.Config{InsecureSkipVerify: false},
		DisableKeepAlives: false,
		Proxy:             http.ProxyFromEnvironment,
	}

	httpClient = &http.Client{Transport: transport}

	fmt.Println("hi")
	//We could have more hostnames, but for now lets just do one
	hostnames := []string{
		"elliottmgmt.com",
	}
	//first see if we can get the info path from the api
	infoResp, err := httpClient.Get(INFOPATH)
	if err != nil {
		panic(fmt.Errorf("error getting %v: %v", INFOPATH, err))
	}
	if infoResp.StatusCode != http.StatusOK {
		//cannot communicate with api. fail
		panic(fmt.Errorf("received non-200 status code '%v' from api", infoResp.StatusCode))
	}

	for _, curHost := range hostnames {
		//start the scan
		results, err := analyzeEndpoint(curHost)
		if err != nil {
			panic(err)
		}
		fmt.Printf("Results: \n%v", results)
	}

}
