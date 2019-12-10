package main

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs"
	"os"
	"time"
)

func GenerateCredentialChain() (*credentials.Credentials, error) {
	var providers []credentials.Provider

	// Add the environment credential provider
	providers = append(providers, &credentials.EnvProvider{})

	// Create the credentials required to access the API.
	creds := credentials.NewChainCredentials(providers)
	if creds == nil {
		return nil, fmt.Errorf("could not compile valid credential providers from static config, environment, shared, or instance metadata")
	}

	return creds, nil
}

func GetAwsCredentials() (*credentials.Credentials, error) {
	creds, err := GenerateCredentialChain()
	if err != nil {
		return nil, err
	}
	if creds == nil {
		return nil, fmt.Errorf("could not compile valid credential providers from static config, environment, shared, or instance metadata")
	}

	_, err = creds.Get()
	if err != nil {
		return nil, err
	}
	return creds, nil
}

func writeLog(message string) {

	creds, err := GetAwsCredentials()
	if err != nil {
		return
	}

	// Use the credentials we've found to construct an STS session
	awsConfig := aws.Config{Credentials: creds,}
	awsConfig.Endpoint = aws.String(os.Getenv("AWS_LOG_ENDPOINT"))
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		Config: awsConfig,
	}))
	svc := cloudwatchlogs.New(sess)

	logEvent := cloudwatchlogs.InputLogEvent{
		Timestamp: aws.Int64(aws.TimeUnixMilli(time.Now())),
		Message:   aws.String(string(message)),
	}
	streamId := "stream-id"
	createLogStreamInput := cloudwatchlogs.CreateLogStreamInput{
		LogGroupName:  aws.String("app"),
		LogStreamName: aws.String(streamId),
	}

	_, err = svc.CreateLogStream(&createLogStreamInput)

	putLogEventInput := cloudwatchlogs.PutLogEventsInput{
		LogGroupName:  aws.String("app"),
		LogStreamName: aws.String(streamId),
		LogEvents:     []*cloudwatchlogs.InputLogEvent{&logEvent,},
	}

	_, err = svc.PutLogEvents(&putLogEventInput)
	if err != nil {
		fmt.Println(err)
		return
	}
}

func main() {
	writeLog("日本語")
}
