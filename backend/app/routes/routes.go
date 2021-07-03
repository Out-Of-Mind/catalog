package routes

import (
	"github.com/gorilla/csrf"
	"github.com/out-of-mind/catalog/structures"
	vars "github.com/out-of-mind/catalog/variables"
	"github.com/satori/uuid"
	"golang.org/x/crypto/bcrypt"

	"crypto/hmac"
	"crypto/sha256"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"errors"
	"html"
	"html/template"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
	"unicode"
)

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	c, _ := r.Cookie("session_token")
	sessionToken := c.Value

	userId, _ := vars.Cache.Get(vars.CTX, sessionToken).Result()
	vars.Log.Debug("user_id: ", userId)

	vars.Cache.Del(vars.CTX, sessionToken)
	vars.Log.Debug("setting new cookie")

	sessionToken = uuid.NewV4().String()

	_, err := vars.Cache.Set(vars.CTX, sessionToken, userId, 720*time.Hour).Result()
	if err != nil {
		vars.Log.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
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

	tmpl, err := template.ParseFiles(vars.TemplateDir + "index.html")
	if err != nil {
		vars.Log.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("500 Internal Server Error"))
		return
	}

	userIdInt, err := strconv.Atoi(userId)
	if err != nil {
		vars.Log.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("500 Internal Server Error"))
		return
	}

	rows, err := vars.DB.Query("SELECT categories.category_name, categories.category_id FROM categories, users WHERE users.user_id=$1 AND categories.group_id=users.group_id", userIdInt)
	if err != nil {
		vars.Log.Error(err)
	}
	defer rows.Close()

	categoriesMap := make(map[int64]string)

	for rows.Next() {
		var (
			categoryName string
			categoryId   int64
		)
		err = rows.Scan(&categoryName, &categoryId)
		if err != nil {
			vars.Log.Error(err)
		}
		categoriesMap[categoryId] = categoryName
	}

	rows, err = vars.DB.Query("SELECT items.item_name, items.category_id FROM items, categories, users WHERE users.user_id=$1 AND categories.group_id=users.group_id AND items.category_id=categories.category_id", userIdInt)
	if err != nil {
		vars.Log.Error(err)
	}
	defer rows.Close()

	itemsMap := make(map[int64][]string)

	for rows.Next() {
		var (
			itemName   string
			categoryId int64
		)
		err = rows.Scan(&itemName, &categoryId)
		if err != nil {
			vars.Log.Println(err)
		}
		itemsMap[categoryId] = append(itemsMap[categoryId], itemName)
	}

	var indexItems structures.IndexItems

	for id := range categoriesMap {
		var indexData structures.IndexData

		indexData.ID = id
		indexData.CategoryName = categoriesMap[id]
		indexData.CategoryID = strings.ReplaceAll(categoriesMap[id], " ", "_")
		for _, itemName := range itemsMap[id] {
			indexData.ItemNames = append(indexData.ItemNames, itemName)
		}
		indexItems.Items = append(indexItems.Items, indexData)
	}

	sort.Sort(structures.ByID(indexItems.Items))

	var (
		jwt, rjwt structures.JWT

		jwtHeader, rjwtHeader   structures.JWTHeader
		jwtPayload, rjwtPayload structures.JWTPayload
	)

	jwtHeader.Alg = "HS256"
	jwtHeader.Type = "JWT"

	rjwtHeader = jwtHeader

	jwtPayload.Exp = time.Now().Add(15 * time.Minute)
	jwtPayload.Value = userIdInt

	rjwtPayload.Exp = time.Now().Add(24 * 7 * time.Hour)
	rjwtPayload.Value = userIdInt

	jwt.Header = jwtHeader
	jwt.Payload = jwtPayload

	rjwt.Header = rjwtHeader
	rjwt.Payload = rjwtPayload

	jwtStr, err := newJWT(jwt)
	rjwtStr, err := newJWT(rjwt)

	indexItems.JWT = jwtStr
	indexItems.RJWT = rjwtStr

	w.Header().
		Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	tmpl.Execute(w, indexItems)
}

func DashboardHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().
		Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("/dashboard"))
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
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

			showLoginHTML(w, data, error)
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

					showLoginHTML(w, data, error)
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

				showLoginHTML(w, data, error)
				return
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

		showLoginHTML(w, data, error)
	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("405 Method Not Allowed"))
	}
}
func showLoginHTML(w http.ResponseWriter, data structures.LoginData, error structures.ErrorTemplate) {
	tmpl, err := template.ParseFiles(vars.TemplateDir + "login.html")
	if err != nil {
		vars.Log.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("500 Internal Server Error"))
		return
	}

	tmpl.Execute(w, data)
}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	c, _ := r.Cookie("session_token")
	sessionToken := c.Value

	vars.Log.Println("deleting cookie")

	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    "",
		Domain:   ".catalog.cc",
		Expires:  time.Now(),
		HttpOnly: true,
		SameSite: 2,
	})

	vars.Cache.Del(vars.CTX, sessionToken)

	http.Redirect(w, r, "login", http.StatusTemporaryRedirect)
}

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
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

			showRegisterHTML(w, data, error)
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

			showRegisterHTML(w, data, error)
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

			showRegisterHTML(w, data, error)
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

			showRegisterHTML(w, data, error)
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

				showRegisterHTML(w, data, error)
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
	} else if r.Method == "GET" {
		var (
			data  structures.RegisterData
			error structures.ErrorTemplate
		)

		error.Show = false

		data.CSRFToken = csrf.Token(r)
		data.Error = error

		showRegisterHTML(w, data, error)
	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("405 Method Not Allowed"))
	}
}
func showRegisterHTML(w http.ResponseWriter, data structures.RegisterData, error structures.ErrorTemplate) {
	tmpl, err := template.ParseFiles(vars.TemplateDir + "register.html")
	if err != nil {
		vars.Log.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("500 Internal Server Error"))
		return
	}

	tmpl.Execute(w, data)
}
func verifyPassword(s string) bool {
	var (
		hasMinLen = false
		hasLower  = false
		hasNumber = false
	)
	if len(s) >= 8 {
		hasMinLen = true
	}
	for _, char := range s {
		switch {
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsNumber(char):
			hasNumber = true
		}
	}
	return hasMinLen && hasLower && hasNumber
}

func APIHandler(w http.ResponseWriter, r *http.Request) {
	body, _ := ioutil.ReadAll(r.Body)

	if len(body) > 0 {
		var (
			requestJSON  structures.RequestJSON
			responseJSON structures.ResponseJSON
			data         structures.ResponseDataJSON
			jwt          structures.JWT
			rjwt         structures.JWT
		)

		vars.Log.Debug(string(body))

		err := json.Unmarshal(body, &requestJSON)

		if err != nil {
			vars.Log.Error(err)

			w.Header().
				Set("Access-Control-Allow-Origin", "http://catalog.cc")
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("400 Bad Request"))
			return
		}

		if requestJSON.Data.RJWT != "" {
			rjwt, err = validateAndParseJWT(requestJSON.Data.RJWT)
		} else if requestJSON.Data.JWT != "" {
			jwt, err = validateAndParseJWT(requestJSON.Data.JWT)
		} else {
			vars.Log.Debug("jwt or rjwt set to null")

			w.Header().
				Set("Access-Control-Allow-Origin", "http://catalog.cc")
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("401 Unauthorized"))
			return
		}

		if err != nil {
			vars.Log.Error(err)

			w.Header().
				Set("Access-Control-Allow-Origin", "http://catalog.cc")
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("400 Bad Request"))
			return
		} else {
			switch requestJSON.Action {
			case "add_item":
				vars.Log.Println("add_item")

				itemName, _ := unescapeUrl(requestJSON.Data.ItemName)
				categoryName, _ := unescapeUrl(requestJSON.Data.CategoryName)

				itemName = html.EscapeString(itemName)
				categoryName = html.EscapeString(categoryName)

				if _, err := vars.DB.Exec("INSERT INTO items(item_name, category_id) VALUES($1, (SELECT category_id FROM categories WHERE category_name=$2))", itemName, categoryName); err != nil {
					vars.Log.Error(err)
					if strings.Contains(err.Error(), "NULL") {
						data.Error = "Категории с таким именем не найдено!"
						responseJSON.Succes = false
						responseJSON.Data = data

						response, _ := json.Marshal(responseJSON)

						w.Header().
							Set("Access-Control-Allow-Origin", "http://catalog.cc")
						w.WriteHeader(http.StatusOK)
						w.Write(response)
						return
					} else if strings.Contains(err.Error(), "items_item_name_key") {
						data.Error = "У вещей не могут быть одинаковые имена!"
						responseJSON.Succes = false
						responseJSON.Data = data

						response, _ := json.Marshal(responseJSON)

						w.Header().
							Set("Access-Control-Allow-Origin", "http://catalog.cc")
						w.WriteHeader(http.StatusOK)
						w.Write(response)
						return
					}
					w.Header().
						Set("Access-Control-Allow-Origin", "http://catalog.cc")
					w.WriteHeader(http.StatusInternalServerError)
					w.Write([]byte("500 Internal Server Error"))
					return
				}

				responseJSON.Succes = true
				responseJSON.Data = data

				response, _ := json.Marshal(responseJSON)

				w.Header().
					Set("Access-Control-Allow-Origin", "http://catalog.cc")
				w.Header().
					Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				w.Write(response)
			case "delete_item":
				vars.Log.Println("delete_item")

				responseJSON.Succes = true
				responseJSON.Data = data

				itemName, _ := unescapeUrl(requestJSON.Data.ItemName)
				categoryName, _ := unescapeUrl(requestJSON.Data.CategoryName)

				itemName = html.EscapeString(itemName)
				categoryName = html.EscapeString(categoryName)

				result, err := vars.DB.Exec("DELETE FROM items WHERE item_name=$1 AND category_id=(SELECT category_id FROM categories WHERE category_name=$2)", itemName, categoryName)
				if err != nil {
					vars.Log.Error(err)
					w.Header().
						Set("Access-Control-Allow-Origin", "http://catalog.cc")
					w.WriteHeader(http.StatusInternalServerError)
					w.Write([]byte("500 Internal Server Error"))
					return
				}

				rows, err := result.RowsAffected()
				if err != nil || rows == 0 {
					vars.Log.Error(err)
					data.Error = "Нет вещи с таким именем!"
					responseJSON.Succes = false
					responseJSON.Data = data
				}

				response, _ := json.Marshal(responseJSON)

				w.Header().
					Set("Access-Control-Allow-Origin", "http://catalog.cc")
				w.Header().
					Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				w.Write(response)
			case "add_category":
				vars.Log.Println("add_category")

				categoryName, _ := unescapeUrl(requestJSON.Data.CategoryName)
				categoryName = html.EscapeString(categoryName)

				_, err = vars.DB.Exec("INSERT INTO categories(category_name, group_id) VALUES($1, (SELECT group_id FROM users WHERE user_id=$2))", categoryName, jwt.Payload.Value)
				if err != nil {
					vars.Log.Error(err)
					w.Header().
						Set("Access-Control-Allow-Origin", "http://catalog.cc")
					w.WriteHeader(http.StatusInternalServerError)
					w.Write([]byte("500 Internal Server Error"))
					return
				}

				responseJSON.Succes = true
				responseJSON.Data = data

				response, _ := json.Marshal(responseJSON)

				w.Header().
					Set("Access-Control-Allow-Origin", "http://catalog.cc")
				w.Header().
					Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				w.Write(response)
			case "delete_category":
				vars.Log.Println("delete_category")

				categoryName, _ := unescapeUrl(requestJSON.Data.CategoryName)
				categoryName = html.EscapeString(categoryName)

				responseJSON.Succes = true
				responseJSON.Data = data

				result, err := vars.DB.Exec("DELETE FROM categories WHERE category_name=$1 AND group_id=(SELECT group_id FROM users WHERE user_id=$2)", categoryName, jwt.Payload.Value)
				if err != nil {
					vars.Log.Error(err)
					w.Header().
						Set("Access-Control-Allow-Origin", "http://catalog.cc")
					w.WriteHeader(http.StatusInternalServerError)
					w.Write([]byte("500 Internal Server Error"))
					return
				}

				rows, err := result.RowsAffected()
				if err != nil || rows == 0 {
					vars.Log.Error(err)
					data.Error = "Нет такой категории!"
					responseJSON.Succes = false
					responseJSON.Data = data
				}

				response, _ := json.Marshal(responseJSON)

				w.Header().
					Set("Access-Control-Allow-Origin", "http://catalog.cc")
				w.Header().
					Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				w.Write(response)
			case "new_jwt":
				vars.Log.Println("new_jwt")

				if jwt.Payload.Exp.Sub(time.Now()).Seconds() >= 30 {
					vars.Log.Debug("big expiration: ", jwt.Payload.Exp.Sub(time.Now()).Seconds())
					w.Header().
						Set("Access-Control-Allow-Origin", "http://catalog.cc")
					w.WriteHeader(http.StatusInternalServerError)
					w.Write([]byte("500 Internal Server Error"))
					return
				}

				jwt.Payload.Exp = time.Now().Add(15 * time.Minute)

				jwtStr, err := newJWT(jwt)
				if err != nil {
					vars.Log.Error(err)
					w.Header().
						Set("Access-Control-Allow-Origin", "http://catalog.cc")
					w.WriteHeader(http.StatusInternalServerError)
					w.Write([]byte("500 Internal Server Error"))
					return
				}

				data.JWT = jwtStr

				responseJSON.Succes = true
				responseJSON.Data = data

				response, _ := json.Marshal(responseJSON)

				vars.Log.Debug("new jwt is set")

				w.Header().
					Set("Access-Control-Allow-Origin", "http://catalog.cc")
				w.Header().
					Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				w.Write(response)
			case "new_rjwt":
				vars.Log.Println("new_rjwt")

				if rjwt.Payload.Exp.Sub(time.Now()).Seconds() >= 30 {
					vars.Log.Debug("big expiration: ", rjwt.Payload.Exp.Sub(time.Now()).Seconds())
					w.Header().
						Set("Access-Control-Allow-Origin", "http://catalog.cc")
					w.WriteHeader(http.StatusInternalServerError)
					w.Write([]byte("500 Internal Server Error"))
					return
				}

				var (
					jwt structures.JWT

					jwtHeader  structures.JWTHeader
					jwtPayload structures.JWTPayload
				)

				jwtHeader.Alg = "HS256"
				jwtHeader.Type = "JWT"
				jwtPayload.Value = rjwt.Payload.Value

				jwt.Header = jwtHeader
				jwt.Payload = jwtPayload

				jwt.Payload.Exp = time.Now().Add(15 * time.Minute)
				rjwt.Payload.Exp = time.Now().Add(24 * 7 * time.Hour)

				jwtStr, err := newJWT(jwt)
				if err != nil {
					vars.Log.Error(err)
					w.Header().
						Set("Access-Control-Allow-Origin", "http://catalog.cc")
					w.WriteHeader(http.StatusInternalServerError)
					w.Write([]byte("500 Internal Server Error"))
					return
				}
				rjwtStr, err := newJWT(rjwt)
				if err != nil {
					vars.Log.Error(err)
					w.Header().
						Set("Access-Control-Allow-Origin", "http://catalog.cc")
					w.WriteHeader(http.StatusInternalServerError)
					w.Write([]byte("500 Internal Server Error"))
					return
				}

				data.JWT = jwtStr
				data.RJWT = rjwtStr

				responseJSON.Succes = true
				responseJSON.Data = data

				response, _ := json.Marshal(responseJSON)

				vars.Log.Debug("new rjwt and jwt is set")

				w.Header().
					Set("Access-Control-Allow-Origin", "http://catalog.cc")
				w.Header().
					Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				w.Write(response)
			default:
				vars.Log.Debug("default")
				w.Header().
					Set("Access-Control-Allow-Origin", "http://catalog.cc")
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte("400 Bad Request"))
			}
		}
	} else {
		vars.Log.Debug("body is empty")

		w.Header().
			Set("Access-Control-Allow-Origin", "http://catalog.cc")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("400 Bad Request"))
	}
}
func unescapeUrl(path string) (string, error) {
	unescapedPath, err := url.PathUnescape(path)

	return unescapedPath, err
}
func validateAndParseJWT(jwt string) (structures.JWT, error) {
	jwtPieces := strings.Split(jwt, ".")
	header := jwtPieces[0]
	payload := jwtPieces[1]
	signature := jwtPieces[2]

	signatureBytes, err := base64.RawURLEncoding.DecodeString(signature)
	if err != nil {
		return structures.JWT{}, err
	}
	message := header + "." + payload

	if validMAC([]byte(message), signatureBytes, []byte(vars.Secret)) {
		headerBytes, err := base64.RawURLEncoding.DecodeString(header)
		if err != nil {
			return structures.JWT{}, err
		}

		payloadBytes, err := base64.RawURLEncoding.DecodeString(payload)
		if err != nil {
			return structures.JWT{}, err
		}

		var (
			jwt          structures.JWT
			jwtHeader    structures.JWTHeader
			jwtPayload   structures.JWTPayload
			jwtSignature structures.JWTSignature
		)

		err = json.Unmarshal(headerBytes, &jwtHeader)
		if err != nil {
			return structures.JWT{}, err
		}

		err = json.Unmarshal(payloadBytes, &jwtPayload)
		if err != nil {
			return structures.JWT{}, err
		}

		if jwtPayload.Exp.Sub(time.Now()).Minutes() <= 0 {
			return structures.JWT{}, errors.New("jwt: jwt token expired")
		}

		jwtSignature.Hash = string(signatureBytes)
		jwt.Header = jwtHeader
		jwt.Payload = jwtPayload

		jwt.Signature = jwtSignature

		return jwt, nil
	} else {
		return structures.JWT{}, errors.New("jwt: signatures isn't matched")
	}
}
func validMAC(message, messageMAC, key []byte) bool {
	mac := hmac.New(sha256.New, key)
	mac.Write(message)
	expectedMAC := mac.Sum(nil)

	return hmac.Equal(messageMAC, expectedMAC)
}
func newJWT(jwt structures.JWT) (string, error) {
	header, err := json.Marshal(jwt.Header)
	if err != nil {
		return "", err
	}
	payload, err := json.Marshal(jwt.Payload)
	if err != nil {
		return "", err
	}

	headerStr := base64.RawURLEncoding.EncodeToString(header)
	payloadStr := base64.RawURLEncoding.EncodeToString(payload)

	message := headerStr + "." + payloadStr

	mac := hmac.New(sha256.New, []byte(vars.Secret))
	mac.Write([]byte(message))
	signature := mac.Sum(nil)

	signatureStr := base64.RawURLEncoding.EncodeToString(signature)

	JWT := message + "." + signatureStr

	return JWT, nil
}
