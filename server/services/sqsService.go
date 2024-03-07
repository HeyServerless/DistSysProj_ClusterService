package services

import (
	"log"

	"encoding/json"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/gin-gonic/gin"
	//"github.com/canhlinh/sqsworker"
)

type Inbound_SQS_Message struct {
	UUID       string
	Expression string
}

type Outbound_SQS_Message struct {
	UUID   string
	Result float32
}

func CreateSession() *session.Session {
	// create a new session with your AWS credentials
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-east-1"), // replace with your desired region
		Credentials: credentials.NewStaticCredentials(
			"AKIAZD66RPV3STVFMGWT",                     // replace with your access key ID
			"zGa2gtrJs788r97JI6A3mG+fF8ED2x5Wx6yf0zqy", // replace with your secret access key
			""), // replace with your session token
	})
	if err != nil {
		log.Fatalf("Failed to create session: %s", err)
	}
	return sess
}

func EnqueueRequestToInboundSqs(c *gin.Context, expression string, uuid string) (string, error) {
	// create a new session with your AWS credentials
	sess := CreateSession()

	// create an SQS client
	svc := sqs.New(sess)

	// define the message payload
	message := &Inbound_SQS_Message{
		UUID:       uuid,
		Expression: expression,
	}
	// stringify the message payload
	//message := fmt.Sprintf("Hello World! %d", i)

	jsonBytes, err := json.Marshal(message)
	if err != nil {
		fmt.Println("Error encoding message:", err)
		return "", err
	}

	jsonString := string(jsonBytes)
	// define the SQS queue URL
	queueURL := "https://sqs.us-east-1.amazonaws.com/626995068279/InboundQueue"

	// create the SQS message input object
	input := &sqs.SendMessageInput{
		MessageBody: aws.String(jsonString),
		QueueUrl:    aws.String(queueURL),
	}

	// send the message to the SQS queue
	success, err := svc.SendMessage(input)
	if err != nil {
		fmt.Println("Error sending message:", err)
		return "", err
	}

	fmt.Println("Message sent to queue with ID:", *success.MessageId)
	return *success.MessageId, nil

}

func EnqueueRequestToOutboundSqs(c *gin.Context, result float32, uuid string) (err error) {
	// create a new session with your AWS credentials
	sess := CreateSession()

	// create an SQS client
	svc := sqs.New(sess)

	// define the message payload
	message := &Outbound_SQS_Message{
		UUID:   uuid,
		Result: result,
	}
	// stringify the message payload
	//message := fmt.Sprintf("Hello World! %d", i)

	jsonBytes, err := json.Marshal(message)
	if err != nil {
		fmt.Println("Error encoding message:", err)
		return
	}

	jsonString := string(jsonBytes)
	// define the SQS queue URL
	queueURL := "https://sqs.us-east-1.amazonaws.com/626995068279/OutboundQueue"

	// create the SQS message input object
	input := &sqs.SendMessageInput{
		MessageBody: aws.String(jsonString),
		QueueUrl:    aws.String(queueURL),
	}

	// send the message to the SQS queue
	_, err = svc.SendMessage(input)
	if err != nil {
		log.Fatalf("Failed to send message: %s", err)
	}

	log.Println("Message sent to SQS queue")
	return err
}

// the below function are one time setup for the queue and trigger

func CreateQueue(queueName string, endpointURL string) error {
	// Create a new session in the us-west-2 region.
	sess := CreateSession()

	// Create a new SQS client.
	svc := sqs.New(sess)
	// check if queue exists
	_, err := svc.GetQueueUrl(&sqs.GetQueueUrlInput{
		QueueName: aws.String(queueName),
	})
	if err == nil {
		fmt.Println("queue already exists")
		return nil
	}

	// Create the queue.
	result, err := svc.CreateQueue(&sqs.CreateQueueInput{
		QueueName: aws.String(queueName),
		Attributes: map[string]*string{
			"ReceiveMessageWaitTimeSeconds": aws.String("20"),
			"VisibilityTimeout":             aws.String("60"),
		},
	})
	if err != nil {
		if err.Error() == "QueueAlreadyExists" {
			fmt.Println("queue already exists")
			return nil
		}
		return err

	}
	log.Println("Queue Created:", *result.QueueUrl)

	queueUrl := *result.QueueUrl

	// add the http endpoint to the queue as trigger by queue on message received
	err = AddQueueTrigger(queueUrl, endpointURL)

	if err != nil {
		fmt.Println("failed to create queue trigger,", err)
		return err
	}

	fmt.Println("queue created with trigger:", queueUrl)

	return nil
}

func AddQueueTrigger(queueUrl string, endpointURL string) error {
	// Create a new session in the us-west-2 region.
	sess := CreateSession()

	// Create a new SQS client.
	svc := sqs.New(sess)

	// Configure the HTTP endpoint for dequeuing the messages.
	_, err := svc.SetQueueAttributes(&sqs.SetQueueAttributesInput{
		QueueUrl: aws.String(queueUrl),
		Attributes: map[string]*string{
			"ReceiveMessageWaitTimeSeconds": aws.String("20"),
			"VisibilityTimeout":             aws.String("60"),
			"Policy": aws.String(fmt.Sprintf(`{
				"Version": "2012-10-17",
				"Id": "Policy%s",
				"Statement": [
					{
						"Sid": "Stmt%s",
						"Effect": "Allow",
						"Principal": "*",
						"Action": "sqs:SendMessage",
						"Resource": "%s",
						"Condition": {
							"ArnEquals": {
								"aws:SourceArn": "%s"
							}
						}
					}
				]
			}`, queueUrl, queueUrl, queueUrl, endpointURL)),
		},
	})
	if err != nil {
		return err
	}

	return nil
}
