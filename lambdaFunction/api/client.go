package api

import (
	"net/http"
	"os"
)

type DatadogClient struct {
	APIKey  string
	AppKey  string
	BaseURL string
	Client  *http.Client
}

func NewDatadogClient() *DatadogClient {
	// Datadog API credentials
	datadogAPIKey := os.Getenv("DATADOG_API_KEY")
	datadogAppKey := os.Getenv("DATADOG_APP_KEY")

	return &DatadogClient{
		APIKey:  datadogAPIKey,
		AppKey:  datadogAppKey,
		BaseURL: "https://api.datadoghq.com/api/v1/monitor",
		Client:  &http.Client{},
	}
}
