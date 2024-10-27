package tracer

import (
	"context"
	"fmt"

	"github.com/guilhermealegre/go-clean-arch-infrastructure-lib/sqs/middlewares"

	"github.com/guilhermealegre/go-clean-arch-infrastructure-lib/domain"

	"github.com/aws/aws-sdk-go/service/sqs"
)

type tracerMiddleware struct {
	app domain.IApp
}

// NewTracerMiddleware returns a new tracerMiddleware
func NewTracerMiddleware(app domain.IApp) middlewares.Middleware {
	return &tracerMiddleware{
		app: app,
	}
}

// Run implements the middleware
func (t *tracerMiddleware) Run(ctx context.Context, data any, err error) context.Context {
	t.app.Tracer().Trace(ctx, t.app.SQS().Name(), t.mapData(data), err)
	return ctx
}

// mapData maps data according to its type
func (t *tracerMiddleware) mapData(data any) map[string]any {
	switch d := data.(type) {
	case *sqs.ReceiveMessageOutput:
		return t.getReceiveMessageOutput(d)
	case []*sqs.SendMessageBatchRequestEntry:
		return t.getSendMessageBatchInput(d)
	}

	return nil
}

// getReceiveMessageOutput gets data of type ReceiveMessageOutput
func (t *tracerMiddleware) getReceiveMessageOutput(messageOutput *sqs.ReceiveMessageOutput) map[string]any {
	attrs := make(map[string]any)
	if messageOutput != nil {
		attrs["message.output"] = messageOutput.String()
	}
	return attrs
}

// getSendMessageBatchInput gets data of type SendMessageBatchInput
func (t *tracerMiddleware) getSendMessageBatchInput(batch []*sqs.SendMessageBatchRequestEntry) map[string]any {
	attrs := make(map[string]any)
	for i, entry := range batch {
		if entry != nil {
			attrs[fmt.Sprintf("message[%d].input", i+1)] = entry.String()
		}
	}
	return attrs
}
