package security

import (
	"context"
	"github.com/kataras/iris/v12"
	"github.com/lestrrat-go/jwx/jwk"
	"github.com/lestrrat-go/jwx/jwt"
	"time"
)

func SetUpOAuth2(app *iris.Application, jwk Jwk, role string) {
	sets, _ := jwk.JwkSets()
	var middleware = NewOAuth2Middleware(sets, role)
	app.Use(middleware)
}

func NewOAuth2Middleware(keySet jwk.Set, allowedAuthority string) func(ctx iris.Context) {
	return func(ctx iris.Context) {
		authorization := authorizationHeaderFor(ctx)

		jwt, _ := jwt.ParseString(authorization)
		if time.Now().After(jwt.Expiration()) {
			ctx.StatusCode(401)
			return
		}
		userName, _ := jwt.PrivateClaims()["user_name"].(string)
		authorities, _ := jwt.PrivateClaims()["authorities"].([]interface{})

		if ok := contains(*toStringSlice(authorities), allowedAuthority); !ok {
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
