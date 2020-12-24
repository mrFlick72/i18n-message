package heath

import (
	"github.com/kataras/iris/v12"
)

type HealthEndpoint struct {
}

func (endpoint *HealthEndpoint) ResgisterEndpoints(application *iris.Application) {
	application.Get("/health/liveness", endpoint.livenessEndpoint)
	application.Get("/health/readiness", endpoint.readinessEndpoint)
}

func (endpoint *HealthEndpoint) readinessEndpoint(ctx iris.Context) {
	ctx.StatusCode(iris.StatusOK)
}

func (endpoint *HealthEndpoint) livenessEndpoint(ctx iris.Context) {
	ctx.StatusCode(iris.StatusOK)
}
