package tracing

import (
	"github.com/kataras/iris/v12"
)

var TracingHeaders = []string{
	"x-request-id",
	"x-b3-traceid",
	"x-b3-spanid",
	"x-b3-sampled",
	"x-b3-parentspanid",
	"x-b3-flags",
	"x-ot-span-context",
}

func GetTracingHeadersFrom(ctx iris.Context) map[string]string {
	headers := make(map[string]string)
	for _, key := range TracingHeaders {
		if value := ctx.GetHeader(key); value != "" {
			headers[key] = value
		}

	}
	return headers
}
