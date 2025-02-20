package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type MonitorCurrentStateResp struct {
	MonitorName         string `json:"name"`
	MonitorOverAllState string `json:"overall_state"` // can be "No Data" or "OK" or "Alert" or "Warn"
	MonitorPriority     int    `json:"priority"`
}

func (d *DatadogClient) sendRequest(method, url string, payload interface{}) (*http.Response, error) {
	var body io.Reader
	if payload != nil {
		jsonData, err := json.Marshal(payload)
		if err != nil {
			return nil, fmt.Errorf("error marshalling JSON: %v", err)
		}
		body = bytes.NewBuffer(jsonData)
	}

	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("DD-API-KEY", d.APIKey)
	req.Header.Set("DD-APPLICATION-KEY", d.AppKey)

	resp, err := d.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error sending request: %v", err)
	}

	return resp, nil
}

func (d *DatadogClient) MuteMonitor(monitorID int, duration time.Duration) error {
	muteEndTime := time.Now().Add(duration).Unix()
	url := fmt.Sprintf("%s/%d/mute", d.BaseURL, monitorID)
	payload := map[string]int64{"end": muteEndTime}

	resp, err := d.sendRequest("POST", url, payload)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to mute monitor, status code: %d", resp.StatusCode)
	}

	fmt.Printf("Monitor %d muted untill %v.\n", monitorID, time.Unix(muteEndTime, 0).Local())
	return nil
}

func (d *DatadogClient) UnmuteMonitor(monitorID int) error {
	url := fmt.Sprintf("%s/%d/unmute", d.BaseURL, monitorID)

	resp, err := d.sendRequest("POST", url, nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to unmute monitor, status code: %d", resp.StatusCode)
	}

	fmt.Printf("Monitor %d unmuted successfully.\n", monitorID)
	return nil
}

func (d *DatadogClient) MonitorCurrentState(monitorID int) (*MonitorCurrentStateResp, error) {
	var monitorDetails *MonitorCurrentStateResp
	url := fmt.Sprintf("%s/%d", d.BaseURL, monitorID)

	resp, err := d.sendRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch monitor state, status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %v", err)
	}

	err = json.Unmarshal(body, &monitorDetails)
	if err != nil {
		return nil, fmt.Errorf("could not unmarshal monitor current state response : %w", err)
	}

	//fmt.Println("Monitor Details:", string(body))
	fmt.Printf("MonitorDetails: %+v\n", monitorDetails)
	return monitorDetails, nil
}

func (d *DatadogClient) UpdateMonitor(monitorID int) error {
	url := fmt.Sprintf("%s/%d", d.BaseURL, monitorID)
	// For testing, hardcoding monitor thresholds to be updated as below
	criticalThreshold := -1
	warningThreshold := 1
	criticalRecoveryThreshold := -0.009
	warningRecoveryThreshold := 1.001

	// Define the request payload
	payload := map[string]interface{}{
		"query": "max(last_5m):max:datadog_demo.Count{type:counter1} by {type} < -1",
		"options": map[string]interface{}{
			"thresholds": map[string]interface{}{
				"critical":          criticalThreshold,
				"critical_recovery": warningThreshold,
				"warning":           criticalRecoveryThreshold,
				"warning_recovery":  warningRecoveryThreshold,
			},
		},
	}

	resp, err := d.sendRequest("PUT", url, payload)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to update monitor, status code: %d", resp.StatusCode)
	}

	fmt.Printf("Monitor %d Updated successfully.\n", monitorID)
	return nil
}
