package api

import (
	"github.com/kataras/iris/v12"
	"github/mrflick72/i18n-message/src/internal/message/repository"
)

type MessageEndpoints struct {
	Repository repository.MessageRepository
}

func (endpoint *MessageEndpoints) ResgisterEndpoint(application *iris.Application) {
	application.Get("/messages/{application}", endpoint.getMessagesEndpoint)
}

func (endpoint *MessageEndpoints) getMessagesEndpoint(ctx iris.Context) {
	application := ctx.Params().Get("application")
	lang := urlParam(ctx, "lang", "")
	find, _ := endpoint.Repository.Find(application, lang)

	ctx.JSON(find)
	ctx.StatusCode(iris.StatusOK)
	return
}

func urlParam(ctx iris.Context, paramName string, defaultValue string) string {
	lang := ctx.URLParam(paramName)
	if &lang == nil {
		return defaultValue
	}
	return lang
}
