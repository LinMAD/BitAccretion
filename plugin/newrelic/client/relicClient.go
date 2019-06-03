package client

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"time"
)

const (
	// Log tag
	tag = "New Relic Client"
	// base endpoint API
	domain  = "api.newrelic.com"
	baseURI = "https://api.newrelic.com/v2"
)

// NRelicClient new relic API client
type NRelicClient struct {
	httpClient      *http.Client
	key             string
	relicErrorCodes map[int]string
}

// NewRelicClient to access to API
func NewRelicClient(APIKey string) *NRelicClient {
	errorCodes := make(map[int]string, 3)
	errorCodes[401] = "Invalid API Key or not given"
	errorCodes[403] = "Access not enabled via API"
	errorCodes[500] = "New Relic API unavailable"

	return &NRelicClient{
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
		key:             APIKey,
		relicErrorCodes: errorCodes,
	}
}

// Authenticate pass authentication and check if API available
func (c *NRelicClient) Authenticate() bool {
	authReq, authReqErr := http.NewRequest("GET", fmt.Sprintf("%s/applications.json", baseURI), nil)
	if authReqErr != nil {
		log.Printf("%s: Unable to create authenticate request", tag)
	}

	authReq.Header.Add("X-Api-Key", c.key)

	authResp, authRespErr := c.httpClient.Do(authReq)
	if authRespErr != nil {
		log.Printf("%s: Unreachable err: %v", tag, authRespErr.Error())
		return false
	}

	defer authResp.Body.Close()

	if errMsg, isFound := c.relicErrorCodes[authResp.StatusCode]; isFound {
		log.Panicf("%s: %s", tag, errMsg)
	}

	return true
}

// isReachableHost check if still have connection (helps to avoid null pointer exceptions)
func (c *NRelicClient) isReachableHost(host, port string) error {
	timeout := time.Duration(5 * time.Second)
	endpoint := host + ":" + port

	conn, err := net.DialTimeout("tcp", endpoint, timeout)
	if err != nil {
		return fmt.Errorf(endpoint+" is unreachable, error: %s", err.Error())
	}

	conn.Close()

	return nil
}

// GetApplicationsList return applications
func (c *NRelicClient) GetApplicationsList() Application {
	if connectionErr := c.isReachableHost(domain, "80"); connectionErr != nil {
		log.Printf("%s: %v", tag, connectionErr.Error())

		time.Sleep(5 * time.Second)
		return c.GetApplicationsList()
	}
	var apps Application

	req, err := http.NewRequest("GET", fmt.Sprintf("%s/applications.json", baseURI), nil)
	if err != nil {
		log.Printf("%s: %v", tag, err.Error())
		return apps
	}

	req.Header.Add("X-Api-Key", c.key)

	resp, respErr := c.httpClient.Do(req)
	if respErr != nil {
		log.Printf("%s: %v", tag, respErr.Error())
		return apps
	}

	defer resp.Body.Close()

	if errMsg, isFound := c.relicErrorCodes[resp.StatusCode]; isFound {
		log.Panicf("%s: %s", tag, errMsg)
		return apps
	}

	respBody, readErr := ioutil.ReadAll(resp.Body)
	if readErr != nil {
		log.Printf("%s: %v", tag, readErr.Error())
		return apps
	}

	json.Unmarshal(respBody, &apps)

	return apps
}

// GetApplicationHost returns stats by application Id
func (c *NRelicClient) GetApplicationHost(appID string) ApplicationHost {
	if connectionErr := c.isReachableHost(domain, "80"); connectionErr != nil {
		log.Printf("%s: %v", tag, connectionErr.Error())

		time.Sleep(5 * time.Second)
		return c.GetApplicationHost(appID)
	}

	var appsHost ApplicationHost

	req, err := http.NewRequest("GET", fmt.Sprintf("%s/applications/%s/hosts.json", baseURI, appID), nil)
	if err != nil {
		log.Printf("%s: %v", tag, err.Error())
		return appsHost
	}

	req.Header.Add("X-Api-Key", c.key)

	resp, respErr := c.httpClient.Do(req)
	if respErr != nil {
		log.Printf("%s: Unreachable New Relic err: %v", tag, respErr.Error())
		return appsHost
	}

	defer resp.Body.Close()

	if errMsg, isFound := c.relicErrorCodes[resp.StatusCode]; isFound {
		log.Printf("%s: Communication error with New Relic %s", tag, errMsg)
		return appsHost
	}

	respBody, readErr := ioutil.ReadAll(resp.Body)
	if readErr != nil {
		log.Printf("%s: %v", tag, readErr.Error())
		return appsHost
	}

	json.Unmarshal(respBody, &appsHost)

	return appsHost
}

// GetHostMetricData metrics for application and metrics
func (c *NRelicClient) GetHostMetricData(appID string, hostID int, metricsNames []string) MetricsData {
	if connectionErr := c.isReachableHost(domain, "80"); connectionErr != nil {
		log.Printf("%s: %v", tag, connectionErr.Error())

		time.Sleep(5 * time.Second)
		return c.GetHostMetricData(appID, hostID, metricsNames)
	}

	var hostMetrics MetricsData

	now := time.Now()
	from := now.Add(time.Duration(-1) * time.Minute)

	req, err := http.NewRequest("GET", fmt.Sprintf("%s/applications/%s/hosts/%d/metrics/data.json", baseURI, appID, hostID), nil)
	if err != nil {
		log.Printf("%s: %v", tag, err.Error())
		return hostMetrics
	}

	q := req.URL.Query()
	q.Add("from", from.Truncate(time.Minute).UTC().Format(time.RFC3339))
	q.Add("to", now.Truncate(time.Minute).UTC().Format(time.RFC3339))
	for _, name := range metricsNames {
		q.Add("names[]", name)
	}

	req.URL.RawQuery = q.Encode()
	req.Header.Add("X-Api-Key", c.key)

	resp, respErr := c.httpClient.Do(req)
	if respErr != nil {
		log.Printf("%s: %v", tag, respErr.Error())
		return hostMetrics
	}

	defer resp.Body.Close()

	if errMsg, isFound := c.relicErrorCodes[resp.StatusCode]; isFound {
		log.Printf("%s: Unreachable New Relic err: %s", tag, errMsg)
		return hostMetrics
	}

	respBody, readErr := ioutil.ReadAll(resp.Body)
	if readErr != nil {
		log.Printf("%s: %v", tag, readErr.Error())
		return hostMetrics
	}

	json.Unmarshal(respBody, &hostMetrics)

	return hostMetrics
}
