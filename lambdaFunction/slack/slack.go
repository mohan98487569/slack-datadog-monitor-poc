package slack

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
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

	var slackMsg struct {
		Messages []Message `json:"messages"`
	}

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
