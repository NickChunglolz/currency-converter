package client

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

type EnvKey string

const (
	API_ADDRESS EnvKey = "API_ADDRESS"
	API_TOKEN   EnvKey = "API_TOKEN"

	RATE_PATH     string = "rates"
	CURRENCY_PATH string = "currencies"

	GET_RATE_PATTERN     string = "%s/%s?source=%s&target=%s"
	GET_CURRENCY_PATTERN string = "%s/%s"
)

type RateReply struct {
	Rate   float64 `json:"rate"`
	Source string  `json:"source"`
	Target string  `json:"target"`
	Time   string  `json:"time"`
}

type CurrencyReply struct {
	Code             string   `json:"code"`
	Symbol           string   `json:"symbol"`
	Name             string   `json:"name"`
	CountryKeywords  []string `json:"countryKeywords"`
	SupportsDecimals bool     `json:"supportsDecimals"`
}

type Client struct {
	baseURL    string
	httpClient *http.Client
}

func NewClient() *Client {
	return &Client{
		baseURL:    os.Getenv(string(API_ADDRESS)),
		httpClient: &http.Client{},
	}
}

func loadApiToken(req *http.Request) (*http.Request, error) {
	apiToken := os.Getenv(string(API_TOKEN))
	if apiToken == "" {
		return nil, fmt.Errorf("API_TOKEN environment variable is not set")
	}

	req.Header.Set("Authorization", "Bearer "+apiToken)

	return req, nil
}

func (client *Client) doHttpRequest(req *http.Request) (*[]byte, error) {
	resp, err := client.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make GET request: %w", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP request failed with status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	return &body, nil
}

func (client *Client) GetRates(source string, target string) ([]RateReply, error) {

	req, err := client.generateGetRatesRequest(source, target)
	if err != nil {
		return nil, err
	}

	body, err := client.doHttpRequest(req)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var apiResponse []RateReply
	if err := json.Unmarshal(*body, &apiResponse); err != nil {
		return nil, fmt.Errorf("failed to decode JSON response: %w", err)
	}

	return apiResponse, nil
}

func (client *Client) GetCurrencies() ([]CurrencyReply, error) {

	req, err := client.generateGetCurrenciesRequest()
	if err != nil {
		return nil, err
	}

	body, err := client.doHttpRequest(req)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var apiResponse []CurrencyReply
	if err := json.Unmarshal(*body, &apiResponse); err != nil {
		return nil, fmt.Errorf("failed to decode JSON response: %w", err)
	}

	return apiResponse, nil
}

func (client *Client) generateGetRatesRequest(source string, target string) (*http.Request, error) {
	// Construct the URL
	url := fmt.Sprintf(GET_RATE_PATTERN, client.baseURL, RATE_PATH, source, target)

	// Create a new request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	loadApiToken(req)

	return req, nil
}

func (client *Client) generateGetCurrenciesRequest() (*http.Request, error) {
	// Construct the URL
	url := fmt.Sprintf(GET_CURRENCY_PATTERN, client.baseURL, CURRENCY_PATH)

	// Create a new request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	loadApiToken(req)

	return req, nil
}
