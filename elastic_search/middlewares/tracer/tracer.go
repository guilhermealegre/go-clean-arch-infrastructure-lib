package tracer

import (
	"bytes"
	"fmt"
	"io"
	"net/http"

	"github.com/guilhermealegre/go-clean-arch-infrastucture-lib/domain"

	errorsInfra "github.com/guilhermealegre/go-clean-arch-infrastucture-lib/errors"
	httpInfra "github.com/guilhermealegre/go-clean-arch-infrastucture-lib/http"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"

	"github.com/guilhermealegre/go-clean-arch-infrastucture-lib/tracer"

	"github.com/elastic/go-elasticsearch/v7/estransport"
)

type tracerMiddleware struct {
	app       domain.IApp
	transport estransport.Interface
}

// NewTracerMiddleware returns a new tracerMiddleware
func NewTracerMiddleware(app domain.IApp, client *elasticsearch.Client) estransport.Interface {
	return &tracerMiddleware{
		app:       app,
		transport: client.Transport,
	}
}

// Perform implements the middleware
func (t *tracerMiddleware) Perform(request *http.Request) (*http.Response, error) {
	if request.URL.Path == "/" {
		// productCheck, it should be ignored!
		return t.transport.Perform(request)
	}

	attrs := make(map[string]any)
	attrs[tracer.TracerTagHttpMethod] = request.Method
	attrs[tracer.TracerTagPath] = request.URL.Path

	if request.Body != nil {
		reqBody, _ := io.ReadAll(request.Body)
		attrs[tracer.TracerTagRequestBody] = string(reqBody)
		// resetting the body buffer to the request
		request.Body = httpInfra.NewReader(bytes.NewBuffer(reqBody), true)
	}

	otel.GetTextMapPropagator().Inject(request.Context(), propagation.HeaderCarrier(request.Header))
	response, err := t.transport.Perform(request)
	var errConn error
	if response != nil {
		if response.StatusCode >= http.StatusBadRequest {
			if err == nil {
				err = errorsInfra.ErrorElastic()
			}

			errConn = fmt.Errorf("error [%s]", response.Status)
		}

		if response.Body != nil {
			respBody, _ := io.ReadAll(response.Body)
			attrs[tracer.TracerTagResponseBody] = string(respBody)
			response.Body = httpInfra.NewReader(bytes.NewBuffer(respBody), true)
		}
	}

	t.app.Tracer().Trace(request.Context(), t.app.ElasticSearch().Name(), attrs, errConn)

	if response == nil {
		response = &http.Response{}
	}

	return response, err
}
