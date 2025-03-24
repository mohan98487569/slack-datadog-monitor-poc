package slack

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"time"
)

// SlackService handles interactions with the Slack API
type SlackClient struct {
	BaseURL string
	Client  *http.Client
	Token   string
}

func NewSlackClient() *SlackClient {
	slackToken := os.Getenv("SLACK_TOKEN")
	return &SlackClient{
		BaseURL: "https://slack.com/api",
		Client:  &http.Client{},
		Token:   slackToken,
	}
}

// Message represents a single Slack message
type Message struct {
	Text        string       `json:"text"`
	User        string       `json:"user"`
	Timestamp   string       `json:"ts"`
	Attachments []Attachment `json:"attachments,omitempty"`
}

type Attachment struct {
	Text      string `json:"text"`
	Title     string `json:"title"`
	TitleLink string `json:"title_link"`
}

// SlackMessages represents messages from Slack API
type SlackMessages struct {
	Messages []Message `json:"messages"`
}

// FetchThreadsFirstMessage retrieves the first message in a Slack thread
func (s *SlackClient) FetchThreadsFirstMessage(channel, threadTS string) (*Message, error) {
	url := fmt.Sprintf("%s/conversations.replies?channel=%s&ts=%s", s.BaseURL, channel, threadTS)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+s.Token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch Slack thread messages, status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var slackMsg SlackMessages

	if err := json.Unmarshal(body, &slackMsg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal Slack thread messages: %w", err)
	}

	if len(slackMsg.Messages) == 0 {
		return nil, fmt.Errorf("no messages found in thread")
	}

	return &slackMsg.Messages[0], nil
}

// ExtractMonitorID extracts the monitor ID from a given URL
func ExtractMonitorID(url string) (string, error) {
	re := regexp.MustCompile(`monitors/(\d+)`)
	matches := re.FindStringSubmatch(url)

	if len(matches) > 1 {
		return matches[1], nil
	}
	return "", fmt.Errorf("monitor ID not found in URL")
}

type PostMessageResponse struct {
	OK        bool   `json:"ok"`
	Timestamp string `json:"ts"`
	Error     string `json:"error,omitempty"`
}

// SendMessage sends a message to a Slack channel
func (s *SlackClient) SendMessage(channelID, message, threadTS string) (*PostMessageResponse, error) {
	fmt.Println("Sending message to slack: ", message)
	url := fmt.Sprintf("%s/chat.postMessage", s.BaseURL)
	var payload map[string]interface{}

	if threadTS != "" {
		payload = map[string]interface{}{
			"channel":   channelID,
			"text":      message,
			"thread_ts": threadTS, // Ensure it is a reply to the thread
		}
	} else {
		payload = map[string]interface{}{
			"channel": channelID,
			"text":    message,
		}
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+s.Token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.Client.Do(req)
	if err != nil {
		return nil, err
	}
	//body, err := io.ReadAll(resp.Body)
	//if err != nil {
	//	fmt.Println("error reading response body: ", err)
	//}
	//fmt.Println("body: ", string(body))
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to send message to Slack, status code: %d", resp.StatusCode)
	}
	fmt.Println("Successfully sent slack new message: ", message)

	var postResp PostMessageResponse
	if err := json.NewDecoder(resp.Body).Decode(&postResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal slack post message response: %v", err)
	}

	if !postResp.OK {
		return nil, fmt.Errorf("slack API error: %s", postResp.Error)
	}

	return &postResp, nil
}

func (s *SlackClient) GetBotMessageTimestamp(channelID, userID, text string) (string, error) {

	// Get start of the day timestamp (Unix format)
	startOfDay := time.Now().Truncate(24 * time.Hour).Unix()

	// Build API request URL
	url := fmt.Sprintf("%s/conversations.history?channel=%s&oldest=%d&limit=100", s.BaseURL, channelID, startOfDay)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("Authorization", "Bearer "+s.Token)

	resp, err := s.Client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("error reading response body: ", err)
		return "", err
	}

	var slackMsg SlackMessages
	err = json.Unmarshal(body, &slackMsg)
	if err != nil {
		fmt.Println("error unmarshalling response body: ", err)
		return "", err
	}

	// Find first message from the user with text "Hi"
	for _, msg := range slackMsg.Messages {
		// slack-bot-to-mute-datadog-monitor Member ID: U08ESQU7G9H
		if msg.User == "U08ESQU7G9H" && msg.Text == text {
			return msg.Timestamp, nil
		}
	}

	// No matching message found
	postMesgResp, err := s.SendMessage(channelID, text, "")
	if err != nil {
		fmt.Errorf("failed to send message to Slack: %w\n", err)
	}
	fmt.Println("postMesgResp.Timestamp: ", postMesgResp.Timestamp)

	return postMesgResp.Timestamp, nil
}
