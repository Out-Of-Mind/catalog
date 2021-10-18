package routes

import (
	"github.com/gorilla/csrf"
	"github.com/out-of-mind/catalog/structures"
	vars "github.com/out-of-mind/catalog/variables"
	"github.com/satori/uuid"
	"golang.org/x/crypto/bcrypt"

	"database/sql"
	"net/http"
	"strconv"
	"time"
	_ "github.com/lib/pq"
)

func loginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		login := r.FormValue("login")
		password := r.FormValue("password")

		if login == "" || password == "" {
			vars.Log.Println("login or password set to null")

			var (
				data  structures.LoginData
				error structures.ErrorTemplate
			)

			error.Show = true
			error.Text = "Логин или пароль не указан!"

			data.CSRFToken = csrf.Token(r)
			data.Error = error

			showHTML(w, "login.html", data)
		} else {
			var (
				result string
				userId int
			)

			row := vars.DB.QueryRow("SELECT password, user_id FROM users WHERE user_name=$1", login)
			err := row.Scan(&result, &userId)
			vars.Log.Debug(result, userId)

			if err != nil {
				if err == sql.ErrNoRows {
					vars.Log.Println("Not found user by user_name")

					var (
						data  structures.LoginData
						error structures.ErrorTemplate
					)

					error.Show = true
					error.Text = "Логин не верный!"

					data.CSRFToken = csrf.Token(r)
					data.Error = error

					showHTML(w, "login.html", data)
					return
				}
				vars.Log.Error(err)
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte("500 Internal Server Error"))
				return
			}

			if err = bcrypt.CompareHashAndPassword([]byte(result), []byte(password)); err != nil {
				vars.Log.Println("Password isn't match with db password")

				var (
					data  structures.LoginData
					error structures.ErrorTemplate
				)

				error.Show = true
				error.Text = "Пароль не верный!"

				data.CSRFToken = csrf.Token(r)
				data.Error = error

				showHTML(w, "login.html", data)
				return
			}

			c, err := r.Cookie("session_token")
			if err == nil {
				vars.Cache.Del(vars.CTX, c.Value)
			}

			sessionToken := uuid.NewV4().String()

			_, err = vars.Cache.Set(vars.CTX, sessionToken, strconv.Itoa(userId), 720*time.Hour).Result()
			if err != nil {
				vars.Log.Error(err)
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte("500 Internal Server Error"))
				return
			}

			http.SetCookie(w, &http.Cookie{
				Name:     "session_token",
				Value:    sessionToken,
				Domain:   ".catalog.cc",
				Expires:  time.Now().Add(720 * time.Hour),
				HttpOnly: true,
				SameSite: 2,
			})

			http.Redirect(w, r, "dashboard", http.StatusTemporaryRedirect)
		}
	} else if r.Method == "GET" {
		var (
			data  structures.LoginData
			error structures.ErrorTemplate
		)

		error.Show = false
		data.CSRFToken = csrf.Token(r)
		data.Error = error

		showHTML(w, "login.html", data)
	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("405 Method Not Allowed"))
	}
}
