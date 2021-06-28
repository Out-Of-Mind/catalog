package middlewares

import (
    vars "github.com/out-of-mind/catalog/variables"
    "github.com/go-redis/redis/v8"

    "net/http"
    "log"
)

func LoggingMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

        log.Println(r.RequestURI)

        if r.RequestURI == "/login" {
            log.Println("login Handler")
            next.ServeHTTP(w, r)
        } else {
            c, err := r.Cookie("session_token")
            if err != nil {
                if err == http.ErrNoCookie {
                    log.Println("no cookie, redirect")
                    http.Redirect(w, r, "login", http.StatusTemporaryRedirect)
                    return
                }
                log.Println("bad request")
                w.WriteHeader(http.StatusBadRequest)
                return
            }
            sessionToken := c.Value

            _, err = vars.Cache.Get(vars.CTX, sessionToken).Result()
            if err == redis.Nil {
                log.Println("no cookie in redis, redirect")
                http.Redirect(w, r, "login", http.StatusTemporaryRedirect)
                return
            } else if err != nil {
                log.Println(err)
            } else {
                log.Println("cookie found, next")
                next.ServeHTTP(w, r)
            }
        }
    })
}