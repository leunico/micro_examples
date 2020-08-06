package main

import (
	"net/http"

	"myauth/lib/token"

	log "github.com/micro/go-micro/v2/logger"
	"github.com/micro/go-plugins/micro/metrics/v2"
	"github.com/micro/micro/v2/cmd"
	"github.com/micro/micro/v2/plugin"
)

func main() {
	tk := &token.Token{}
	tk.Init([]byte("key123456"))
	plugin.Register(plugin.NewPlugin(
		plugin.WithName("auth"),
		plugin.WithHandler(
			JWTAuthWrapper(tk),
		),
	))
	plugin.Register(metrics.NewPlugin())
	cmd.Init()
}

func JWTAuthWrapper(t *token.Token) plugin.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log.Info("===========", r.URL.Path)
			//不需要登录的url地址 strings.HasPrefix(r.URL.Path, "/hello") ||
			if r.URL.Path == "/myauth/Myauth/GetJwt" ||
				r.URL.Path == "/myauth/Myauth/InspectJwt" ||
				r.URL.Path == "/metrics" {
				h.ServeHTTP(w, r)
				return
			}

			// tokenstr := r.Header.Get("Authorization")//现在不可以用Authorization，需要用Bearer
			tokenstr := r.Header.Get("Bearer")
			log.Info("tokenstr", tokenstr)
			userFromToken, e := t.Decode(tokenstr)
			log.Info("userFromToken", userFromToken)
			if e != nil {
				_, _ = w.Write([]byte("unauthorized"))
				return
			}

			// r.Header.Set("X-Example-Username", userFromToken.UserName)
			h.ServeHTTP(w, r)
			return
		})
	}
}
