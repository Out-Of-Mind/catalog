package routes

import (
	vars "github.com/out-of-mind/catalog/variables"

	"net/http"
	"time"
)

func logoutHandler(w http.ResponseWriter, r *http.Request) {
	c, _ := r.Cookie("session_token")
	sessionToken := c.Value

	vars.Log.Println("deleting cookie")

	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    "",
		Domain:   ".catalog.cc",
		Expires:  time.Now(),
		HttpOnly: true,
		Secure:   true,
		SameSite: 2,
	})

	vars.Cache.Del(vars.CTX, sessionToken)

	http.Redirect(w, r, "login", http.StatusTemporaryRedirect)
}
