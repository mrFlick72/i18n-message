package api

import (
	"github.com/kataras/iris/v12"
	"github/mrflick72/i18n-message/internal/message/repository"
)

type MessageEndpoints struct {
	Repository repository.MessageRepository
}

func (endpoint *MessageEndpoints) GetMessagesEndpoint(ctx iris.Context) {
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
