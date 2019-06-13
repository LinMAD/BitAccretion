package client

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"time"
)

const (
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
		httpClient:      &http.Client{Timeout: 10 * time.Second},
		key:             APIKey,
		relicErrorCodes: errorCodes,
	}
}

// Authenticate pass authentication and check if API available
func (c *NRelicClient) Authenticate() (bool, error) {
	authReq, _ := http.NewRequest("GET", fmt.Sprintf("%s/applications.json", baseURI), nil)
	authReq.Header.Add("X-Api-Key", c.key)

	authResp, authRespErr := c.httpClient.Do(authReq)
	if authRespErr != nil {
		return false, nil
	}

	defer authResp.Body.Close()

	if errMsg, isFound := c.relicErrorCodes[authResp.StatusCode]; isFound {
		return false, fmt.Errorf("new relic API err %s", errMsg)
	}

	return true, nil
}

// isReachableHost check if still have connection (helps to avoid null pointer exceptions)
func (c *NRelicClient) isReachableHost(host, port string) error {
	endpoint := host + ":" + port
	conn, err := net.DialTimeout("tcp", endpoint, 1 * time.Second)
	if err != nil {
		return fmt.Errorf("%s is unreachable, error: %s", endpoint, err.Error())
	}
	defer conn.Close()

	return nil
}

// GetApplicationHost returns stats by application Id
func (c *NRelicClient) GetApplicationHost(appID string) (appHost ApplicationHost) {
	if connectionErr := c.isReachableHost(domain, "80"); connectionErr != nil {
		return c.GetApplicationHost(appID)
	}

	req, err := http.NewRequest("GET", fmt.Sprintf("%s/applications/%s/hosts.json", baseURI, appID), nil)
	if err != nil {
		return
	}
	req.Header.Add("X-Api-Key", c.key)

	resp, respErr := c.httpClient.Do(req)
	if respErr != nil {
		return
	}

	defer resp.Body.Close()

	if _, isFound := c.relicErrorCodes[resp.StatusCode]; isFound {
		return
	}

	respBody, readErr := ioutil.ReadAll(resp.Body)
	if readErr != nil {
		return
	}

	unMarshErr := json.Unmarshal(respBody, &appHost)
	if unMarshErr != nil {
		return ApplicationHost{}
	}

	return
}

// GetHostMetricData metrics for application and metrics
func (c *NRelicClient) GetHostMetricData(appID string, hostID int, metricsNames []string) (m MetricsData) {
	if connectionErr := c.isReachableHost(domain, "80"); connectionErr != nil {
		return c.GetHostMetricData(appID, hostID, metricsNames)
	}

	now := time.Now()
	from := now.Add(time.Duration(-1) * time.Minute)

	req, err := http.NewRequest("GET", fmt.Sprintf("%s/applications/%s/hosts/%d/metrics/data.json", baseURI, appID, hostID), nil)
	if err != nil {
		return
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
		return
	}

	defer resp.Body.Close()

	if _, isFound := c.relicErrorCodes[resp.StatusCode]; isFound {
		return
	}

	respBody, readErr := ioutil.ReadAll(resp.Body)
	if readErr != nil {
		return
	}

	unMarshErr := json.Unmarshal(respBody, &m)
	if unMarshErr != nil {
		return MetricsData{}
	}

	return
}
