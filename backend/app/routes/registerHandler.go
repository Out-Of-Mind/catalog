package routes

import (
	"github.com/gorilla/csrf"
	"github.com/out-of-mind/catalog/structures"
	vars "github.com/out-of-mind/catalog/variables"
	"github.com/satori/uuid"
	"golang.org/x/crypto/bcrypt"

	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func registerHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		login := r.FormValue("login")
		email := r.FormValue("email")
		password := r.FormValue("password")
		repassword := r.FormValue("repassword")

		if login == "" || email == "" || password == "" || repassword == "" {
			vars.Log.Error("login or email or password or repassword set to null")

			var (
				data  structures.RegisterData
				error structures.ErrorTemplate
			)

			error.Show = true
			error.Text = "Заполните все поля!"

			data.CSRFToken = csrf.Token(r)
			data.Error = error

			showHTML(w, "register.html", data)
			return
		}

		matched, _ := regexp.Match(`^(([^<>()\[\]\\.,;:\s@"]+(\.[^<>()\[\]\\.,;:\s@"]+)*)|(".+"))@((\[[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}])|(([a-zA-Z\-0-9]+\.)+[a-zA-Z]{2,}))$`, []byte(email))
		if !matched {
			vars.Log.Error("email regex failed")

			var (
				data  structures.RegisterData
				error structures.ErrorTemplate
			)

			error.Show = true
			error.Text = "Введён не верный эмейл!"

			data.CSRFToken = csrf.Token(r)
			data.Error = error

			showHTML(w, "register.html", data)
			return
		}

		if !verifyPassword(password) {
			vars.Log.Error("password regex failed")

			var (
				data  structures.RegisterData
				error structures.ErrorTemplate
			)

			error.Show = true
			error.Text = "Пароль должен содержать: одну маленькую букву, одну цифру, как минимум 8 символов!"

			data.CSRFToken = csrf.Token(r)
			data.Error = error

			showHTML(w, "register.html", data)
			return
		}

		if password != repassword {
			vars.Log.Error("passwords isn't match")

			var (
				data  structures.RegisterData
				error structures.ErrorTemplate
			)

			error.Show = true
			error.Text = "Пароли не совпадают!"

			data.CSRFToken = csrf.Token(r)
			data.Error = error

			showHTML(w, "register.html", data)
			return
		}

		vars.Log.Println("inserting new user")

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 8)
		if err != nil {
			vars.Log.Error(err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("500 Internal Server Error"))
			return
		}

		if _, err := vars.DB.Exec("INSERT INTO users(user_name, email, password) VALUES($1, $2, $3)", login, email, string(hashedPassword)); err != nil {
			vars.Log.Error(err)

			if strings.Contains(err.Error(), "users_user_name_key") {
				vars.Log.Error("password regex failed")

				var (
					data  structures.RegisterData
					error structures.ErrorTemplate
				)

				error.Show = true
				error.Text = "Пользователь с таким логином уже существует!"

				data.CSRFToken = csrf.Token(r)
				data.Error = error

				showHTML(w, "register.html", data)
				return
			}

			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("500 Internal Server Error"))
			return
		}

		var userId int

		row := vars.DB.QueryRow("SELECT user_id FROM users WHERE user_name=$1", login)
		err = row.Scan(&userId)
		if err != nil {
			vars.Log.Error(err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("500 Internal Server Error"))
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
			Secure:   true,
			SameSite: 2,
		})

		http.Redirect(w, r, "dashboard", http.StatusTemporaryRedirect)
	} else if r.Method == "GET" {
		var (
			data  structures.RegisterData
			error structures.ErrorTemplate
		)

		error.Show = false

		data.CSRFToken = csrf.Token(r)
		data.Error = error

		showHTML(w, "register.html", data)
	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("405 Method Not Allowed"))
	}
}
