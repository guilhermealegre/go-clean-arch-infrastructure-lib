package tracer

import (
	"context"

	"github.com/guilhermealegre/go-clean-arch-infrastucture-lib/tracer"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/smithy-go/middleware"
	"github.com/guilhermealegre/go-clean-arch-infrastucture-lib/domain"
)

// tracerMiddleware
type tracerMiddleware struct {
	app domain.IApp
}

// NewTracerMiddleware creates a new TracerMiddleware
func NewTracerMiddleware(app domain.IApp) middleware.InitializeMiddleware {
	return &tracerMiddleware{
		app: app,
	}
}

// ID returns the id of the middlewares
func (t *tracerMiddleware) ID() string {
	return "tracerMiddleware"
}

// HandleInitialize implements the middleware
func (t *tracerMiddleware) HandleInitialize(ctx context.Context, in middleware.InitializeInput, next middleware.InitializeHandler) (
	middleware.InitializeOutput, middleware.Metadata, error) {

	out, metadata, err := next.HandleInitialize(ctx, in)

	t.app.Tracer().Trace(ctx, t.app.S3().Name(), t.mapParameters(in.Parameters), err)

	return out, metadata, err
}

// mapParameters maps the parameters according to its type
func (t *tracerMiddleware) mapParameters(parameters any) map[string]any {
	switch p := parameters.(type) {
	case *s3.PutObjectInput:
		return t.mapPutObjectInput(p)
	}
	return nil
}

// mapPutObjectInput maps data of type PutObjectInput
func (t *tracerMiddleware) mapPutObjectInput(input *s3.PutObjectInput) map[string]any {
	attrs := make(map[string]any)
	if input.Bucket != nil {
		attrs[tracer.TracerTagBucket] = *input.Bucket
	}
	if input.ContentType != nil {
		attrs[tracer.TracerTagContentType] = *input.ContentType
	}
	if input.Key != nil {
		attrs[tracer.TracerTagKey] = *input.Key
	}

	return attrs
}
