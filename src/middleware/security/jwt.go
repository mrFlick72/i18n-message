package security

import (
	"context"
	"fmt"
	"github.com/kataras/iris/v12"
	"github.com/lestrrat-go/jwx/jwk"
	"github.com/lestrrat-go/jwx/jwt"
	"github/mrflick72/i18n-message/src/internal/logging"
	"time"
)

var logger = logging.GetLoggerInstance()

func SetUpOAuth2(app *iris.Application, jwk Jwk, role []string) {
	sets, _ := jwk.JwkSets()
	var middleware = NewOAuth2Middleware(sets, role)
	app.Use(middleware)
}

func NewOAuth2Middleware(keySet jwk.Set, allowedAuthority []string) func(ctx iris.Context) {
	return func(ctx iris.Context) {
		authorization := authorizationHeaderFor(ctx)

		jwt, err := jwt.ParseString(authorization)
		if err != nil {
			logger.LogInfoFor(fmt.Sprintf("failed to create parse jwt: %v", err))

			ctx.StatusCode(401)
			return
		}

		if time.Now().After(jwt.Expiration()) {
			ctx.StatusCode(401)
			return
		}
		userName, _ := jwt.PrivateClaims()["user_name"].(string)
		authorities, _ := jwt.PrivateClaims()["authorities"].([]interface{})
		scopes, _ := jwt.PrivateClaims()["scope"].([]interface{})

		if ok := allowedRole(authorities, allowedAuthority) || allowedRole(scopes, allowedAuthority); !ok {
			ctx.StatusCode(403)
			return
		}

		newContext := context.WithValue(ctx.Request().Context(), "user", OAuth2User{
			UserName: userName,
		})
		ctx.ResetRequest(ctx.Request().WithContext(newContext))
		ctx.Next()
	}
}

func allowedRole(authorities []interface{}, allowedAuthority []string) bool {
	for _, authority := range allowedAuthority {
		if contains(*toStringSlice(authorities), authority) {
			return true
		}
	}
	return false
}

func contains(slice []string, item string) bool {
	set := make(map[string]struct{}, len(slice))
	for _, s := range slice {
		set[s] = struct{}{}
	}

	_, ok := set[item]
	return ok
}

func toStringSlice(slice []interface{}) *[]string {
	result := make([]string, 0)

	for _, item := range slice {
		result = append(result, item.(string))
	}

	return &result
}

func authorizationHeaderFor(ctx iris.Context) string {
	authorization := ctx.GetHeader("Authorization")
	authorization = authorization[7:]
	return authorization
}

type OAuth2User struct {
	UserName    string
	Authorities []string
}
