package sqs

import (
	"encoding/json"
	"fmt"

	"github.com/guilhermealegre/go-clean-arch-infrastructure-lib/sqs/config"

	"github.com/aws/aws-sdk-go-v2/aws"

	"github.com/guilhermealegre/go-clean-arch-infrastructure-lib/domain/message"

	awsSdk "github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	sqsSdk "github.com/aws/aws-sdk-go/service/sqs"
	errorCodes "github.com/guilhermealegre/go-clean-arch-infrastructure-lib/errors"
)

const (
	messageCreatedQueue = "Queue created %s at %s"
)

func NewMigration(name string, config *config.Connection, conn *sqsSdk.SQS) (*Migration, error) {
	return &Migration{
		name:       fmt.Sprintf("%s :: Migration", name),
		config:     config,
		connection: conn,
	}, nil
}

func (m *Migration) canRun() bool {
	return !m.config.MigrationsDisabled
}

func (m *Migration) Run(queue string) error {
	if !m.canRun() {
		return nil
	}

	// create queue
	_, queueUrl, err := m.createQueue(queue)
	if err != nil {
		return err
	}

	// create dead letter queue
	_, _, err = m.createDeadLetterQueue(queue, queueUrl)
	if err != nil {
		return err
	}

	return nil
}

func (m *Migration) createQueue(queue string) (created bool, queueUrl *string, err error) {
	// get queue
	getQueueResult, err := m.connection.GetQueueUrl(
		&sqsSdk.GetQueueUrlInput{
			QueueName:              awsSdk.String(queue),
			QueueOwnerAWSAccountId: awsSdk.String(m.config.Credentials.IdAccount),
		},
	)

	if err == nil {
		// queue already exists
		return false, getQueueResult.QueueUrl, nil
	}

	awsErr, ok := err.(awserr.Error)
	if !ok || awsErr.Code() != sqsSdk.ErrCodeQueueDoesNotExist {
		message.ErrorMessage(m.name, errorCodes.ErrorGettingQueue().Formats(queue, err))
		return false, nil, err
	}

	// create queue
	createQueueResult, err := m.connection.CreateQueue(&sqsSdk.CreateQueueInput{
		QueueName: awsSdk.String(queue),
	})

	if err != nil || createQueueResult == nil {
		message.ErrorMessage(m.name, errorCodes.ErrorCreatingQueue().Formats(queue, err))
		return false, nil, err
	} else {
		message.Message(m.name, fmt.Sprintf(messageCreatedQueue, queue, *createQueueResult.QueueUrl))
		return true, createQueueResult.QueueUrl, nil
	}
}

func (m *Migration) createDeadLetterQueue(queue string, queueUrl *string) (created bool, dltQueueUrl *string, err error) {
	// create queue
	dltQueue := fmt.Sprintf("%s-dlq", queue)
	created, dltQueueUrl, err = m.createQueue(dltQueue)
	if err != nil {
		return created, dltQueueUrl, err
	}

	// configure dead letter queue
	err = m.configureDeadLetterQueue(queueUrl, dltQueue)

	return created, dltQueueUrl, err
}

func (m *Migration) configureDeadLetterQueue(queueUrl *string, dltQueue string) error {
	// check if dead-letter attribute is already set
	queueAttributes, err := m.connection.GetQueueAttributes(
		&sqsSdk.GetQueueAttributesInput{
			QueueUrl: queueUrl,
			AttributeNames: []*string{
				aws.String(sqsSdk.QueueAttributeNameRedrivePolicy),
			},
		},
	)

	if err != nil {
		return err
	}

	if _, ok := queueAttributes.Attributes[sqsSdk.QueueAttributeNameRedrivePolicy]; ok {
		// queue dead-letter attribute already exists
		return nil
	}

	// add dead-letter attribute to queue
	policy := map[string]string{
		"deadLetterTargetArn": m.getQueueArn(dltQueue),
		"maxReceiveCount":     "1",
	}

	bytes, err := json.Marshal(policy)
	if err != nil {
		return err
	}

	_, err = m.connection.SetQueueAttributes(
		&sqsSdk.SetQueueAttributesInput{
			QueueUrl: queueUrl,
			Attributes: map[string]*string{
				sqsSdk.QueueAttributeNameRedrivePolicy: aws.String(string(bytes)),
			},
		},
	)

	return err
}

func (m *Migration) getQueueArn(queue string) string {
	return fmt.Sprintf("arn:aws:sqs:%s:%s:%s", m.config.Credentials.Region, m.config.Credentials.IdAccount, queue)
}
