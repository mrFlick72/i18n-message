package api

import (
	"github.com/kataras/iris/v12"
	"github/mrflick72/i18n-message/src/internal/message/repository"
	"github/mrflick72/i18n-message/src/internal/tracing"
)

type MessageEndpoints struct {
	Repository repository.MessageRepository
}

func (endpoint *MessageEndpoints) RegisterEndpoint(application *iris.Application) {
	application.Get("/messages/{application}", endpoint.getMessagesEndpoint)
}

func (endpoint *MessageEndpoints) getMessagesEndpoint(ctx iris.Context) {
	application := ctx.Params().Get("application")
	lang := urlParam(ctx, "lang", "")

	context := tracingContextFrom(ctx)

	find, _ := endpoint.Repository.Find(application, lang, context)

	ctx.JSON(find)
	ctx.StatusCode(iris.StatusOK)
	return
}

func tracingContextFrom(ctx iris.Context) map[string]string {
	return tracing.GetTracingHeadersFrom(ctx)
}

func urlParam(ctx iris.Context, paramName string, defaultValue string) string {
	lang := ctx.URLParam(paramName)
	if &lang == nil {
		return defaultValue
	}
	return lang
}
