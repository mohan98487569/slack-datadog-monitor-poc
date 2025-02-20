package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"lambda_function/service"
	"net/http"
)

func main() {
	fmt.Println("Starting Slack to Datadog Monitor Automation requests.")
	lambda.Start(handler)
}

//

func handler(ctx context.Context, eventData events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	var slackEvent service.SlackEvent
	fmt.Println("Handling event")
	//fmt.Printf("request: %+v\n", eventData)
	fmt.Println("eventData.Body", eventData.Body)

	err := json.Unmarshal([]byte(eventData.Body), &slackEvent)
	if err != nil {
		fmt.Printf("Error unmarshaling request body: %v\n", err)
		return events.APIGatewayV2HTTPResponse{StatusCode: http.StatusBadRequest}, nil
	}

	fmt.Printf("slackEvent: %+v \n", slackEvent)

	if slackEvent.Event.ThreadTS != "" && (slackEvent.Event.Text == "acknowledged" || slackEvent.Event.Text == "resolved") {
		err := service.ProcessMessage(slackEvent)
		if err != nil {
			fmt.Errorf("Processing message failed %v", err)
		}
	}

	// Handle Slack Challenge (Needed for Slack to verify API)
	if slackEvent.Challenge != "" {
		return events.APIGatewayV2HTTPResponse{
			StatusCode: http.StatusOK,
			Body:       fmt.Sprintf(`{"challenge": "%s"}`, slackEvent.Challenge),
			Headers: map[string]string{
				"Content-Type": "application/json",
			},
		}, nil
	}

	//  Return Success Response
	fmt.Println("Successfully handled request.")
	return events.APIGatewayV2HTTPResponse{
		StatusCode: http.StatusOK,
		Body:       `{"message": "Acknowledgement received"}`,
	}, nil
}
