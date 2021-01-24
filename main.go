package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/sirupsen/logrus"
	"github.com/bcneng/twitter-contest/twitter"
)

// Version stores the git commit SHA from where the app got built. Injected when building.
var Version string

// ResponseBody represents the lambda response body
type ResponseBody struct {
	Winners []string  `json:"winners"`
	Version string    `json:"version"`
	Time    time.Time `json:"time"`
}

func handler(request events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	apiKey, err := getQueryParam(request, "api_key", true)
	if err != nil {
		return nil, err
	}

	apiKeySecret, err := getQueryParam(request, "api_key_secret", true)
	if err != nil {
		return nil, err
	}

	tweetIDStr, err := getQueryParam(request, "tweet_id", true)
	if err != nil {
		return nil, err
	}

	tweetID, err := strconv.Atoi(tweetIDStr)
	if err != nil {
		return nil, errors.New("tweet_id should be a valid integer")
	}

	pickStr, err := getQueryParam(request, "pick", true)
	if err != nil {
		return nil, err
	}

	pick, err := strconv.Atoi(pickStr)
	if err != nil {
		return nil, errors.New("pick should be a valid integer")
	}

	// If not set, will use the account author of the tweet
	account, err := getQueryParam(request, "account_to_follow", false)
	if err != nil {
		return nil, err
	}
	account = strings.TrimPrefix(account, "@")

	winners, err := twitter.Contest(twitter.Credentials{
		APIKey:       apiKey,
		APIKeySecret: apiKeySecret,
	}, tweetID, pick, account)

	if err != nil {
		return nil, err
	}

	if len(winners) > 0 {
		logrus.WithField("winners", winners).Infoln("found winners!")
	} else {
		logrus.Infoln("could not find winners")
	}

	encodedBody, err := json.Marshal(ResponseBody{
		Version: Version,
		Winners: winners,
		Time:    time.Now(),
	})
	if err != nil {
		return nil, err
	}

	return &events.APIGatewayProxyResponse{
		StatusCode:      200,
		Headers:         map[string]string{"Content-Type": "application/json"},
		Body:            string(encodedBody),
		IsBase64Encoded: false,
	}, nil
}

func main() {
	lambda.Start(handler)
}

func getQueryParam(request events.APIGatewayProxyRequest, name string, required bool) (string, error) {
	val, ok := request.QueryStringParameters[name]
	if required && !ok {
		return "", fmt.Errorf("%q query param is mandatory", name)
	}

	return val, nil
}
