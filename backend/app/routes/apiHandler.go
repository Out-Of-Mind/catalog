package routes

import (
	"github.com/out-of-mind/catalog/structures"
	vars "github.com/out-of-mind/catalog/variables"

	"encoding/json"
	"html"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

func apiHandler(w http.ResponseWriter, r *http.Request) {
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
