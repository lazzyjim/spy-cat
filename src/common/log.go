package common

import (
	"github.com/qiangxue/fasthttp-routing"
	"time"
)

func LoggingMiddleware(ctx *routing.Context) error {
	start := time.Now()
	logger := ctx.RequestCtx.Logger()
	logger.Printf("Started %s %s", ctx.Method(), ctx.RequestURI())
	err := ctx.Next()
	logger.Printf("Completed %s in %v with status %d", ctx.RequestURI(), time.Since(start), ctx.Response.StatusCode())
	return err
}
