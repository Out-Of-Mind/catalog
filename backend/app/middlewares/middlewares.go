package middlewares

import (
	"github.com/go-redis/redis/v8"
	"github.com/gorilla/csrf"
	vars "github.com/out-of-mind/catalog/variables"

	"net/http"
)

func LoggingMiddleware(next http.Handler) http.Handler {
	vars.Log.Debug("Logging middleware")

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		vars.Log.Println(r.RequestURI)

		if r.RequestURI == "/login" || r.RequestURI == "/register" || r.RequestURI == "/api" { // test only (!=) in prod set to (==)
			vars.Log.Debug("Next handler without cookies set")
			next.ServeHTTP(w, r)
		} else {
			c, err := r.Cookie("session_token")
			if err != nil {
				if err == http.ErrNoCookie {
					vars.Log.Debug("no cookie, redirect")
					http.Redirect(w, r, "login", http.StatusTemporaryRedirect)
					return
				}
				vars.Log.Debug("bad request")
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			sessionToken := c.Value

			_, err = vars.Cache.Get(vars.CTX, sessionToken).Result()
			if err == redis.Nil {
				vars.Log.Debug("no cookie in redis, redirect")
				http.Redirect(w, r, "login", http.StatusTemporaryRedirect)
				return
			} else if err != nil {
				vars.Log.Debug(err)
			} else {
				vars.Log.Debug("cookie found, next")
				next.ServeHTTP(w, r)
			}
		}
	})
}

func CSRFMiddleware(next http.Handler) http.Handler {
	vars.Log.Debug("CSRF middleware")

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.RequestURI == "/login" || r.RequestURI == "/register" {
			CSRF := csrf.Protect(vars.Secret, csrf.Secure(false))
			// csrf.Secure(false) only for debug
			CSRF(next).ServeHTTP(w, r)
		} else {
			next.ServeHTTP(w, r)
		}
	})
}
