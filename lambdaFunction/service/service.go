package service

import (
	"fmt"
	"lambda_function/api"
	"time"
)

type SlackEvent struct {
	Challenge string `json:"challenge,omitempty"` // slack api bot sends a request with a challenge parameter, and lambda api endpoint must respond with the challenge value. This all is done to verify the Request URL by the slack api bot.
	Event     struct {
		User      string `json:"user"`
		Type      string `json:"type"`
		Timestamp string `json:"ts"`
		Text      string `json:"text"`
		ThreadTS  string `json:"thread_ts"`
		Metadata  struct {
			EventType    string `json:"event_type"`
			EventPayload struct {
				MonitorID int   `json:"monitor_id"`
				EventTS   int64 `json:"event_ts"`
			} `json:"event_payload"`
		} `json:"metadata"`
		EventTimestamp string `json:"event_time"`
	} `json:"event"`
	Type      string `json:"type"`
	EventID   string `json:"event_id"`
	EventTime int64  `json:"event_time"`
}

func ProcessMessage(slackEvent SlackEvent) error {
	fmt.Println("Processing message: ")

	datadogClient := api.NewDatadogClient()
	/*
		ToDo: Retrieve datadog monitor_id from datadog alert received in slack channel
	*/
	monitorID := 164327389 // for testing, hardcoding the datadog monitor_id
	if slackEvent.Event.Text == "acknowledged" {
		monitorData, err := datadogClient.MonitorCurrentState(monitorID)
		if err != nil {
			return err
		}
		// ToDo: After testing, later update if condition, to mute the monitor only if the monitor is in "Alert" state.
		if monitorData.MonitorOverAllState == "Alert" || monitorData.MonitorOverAllState == "OK" || monitorData.MonitorOverAllState == "No Data" {
			err = datadogClient.MuteMonitor(monitorID, 6*time.Hour)
			if err != nil {
				return err
			}
		}

	} else {
		err := datadogClient.UnmuteMonitor(monitorID)
		if err != nil {
			return err
		}
	}

	fmt.Println("Message processed successfully")
	return nil
}
