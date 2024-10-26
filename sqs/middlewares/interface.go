package middlewares

import (
	"context"
)

type Middleware interface {
	Run(ctx context.Context, data any, err error) context.Context
}
