package routes

import (
	vars "github.com/out-of-mind/catalog/variables"
	"github.com/out-of-mind/catalog/structures"

	"html/template"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"log"
)

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	c, err := r.Cookie("session_token")
    sessionToken := c.Value
	user_id, _ := vars.Cache.Get(vars.CTX, sessionToken).Result()

	log.Println(sessionToken, user_id)
	
	tmpl, err := template.ParseFiles(vars.TemplateDir+"index.html")
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("500 Internal Server Error"))
		return
	}

	rows, err := vars.DB.Query("SELECT category_name, item_name FROM categories, items, users WHERE users.user_id=$1 AND categories.group_id=users.group_id AND items.category_id=categories.category_id", user_id)
    if err != nil {
        log.Println(err)
    }
    defer rows.Close()

    var indexItems structures.IndexItems

    data := make(map[string][]string)

    for rows.Next() {
    	var (
        	categoryName, itemName string
    	)
        err = rows.Scan(&categoryName, &itemName)
        if err != nil {
            log.Println(err)
        }
        data[categoryName] = append(data[categoryName], itemName)
    }

    for categoryName := range data {
    	var indexData structures.IndexData
    	indexData.CategoryName = categoryName

    	for _, itemName := range data[categoryName] {
    		indexData.ItemNames = append(indexData.ItemNames, itemName)
    	}

    	indexItems.Items = append(indexItems.Items, indexData)
    }
	
	w.Header().
	Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
    tmpl.Execute(w, indexItems)
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().
	Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("/login"))
}

func APIHandler(w http.ResponseWriter, r *http.Request) {
	body, _ := ioutil.ReadAll(r.Body)

	if len(body) > 0 {
		var requestJSON structures.RequestJSON
		var responseJSON structures.ResponseJSON
		var data structures.ResponseDataJSON

		log.Println(string(body))

		err := json.Unmarshal(body, &requestJSON)

		if err != nil {
			log.Println(err)

			w.Header().
			Set("Access-Control-Allow-Origin", "http://catalog.cc")
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("smth"))
		} else {
			switch requestJSON.Action {
			case "add_item":
				log.Println("add_item")

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
				log.Println("delete_item")

				responseJSON.Succes = true
				responseJSON.Data = data
	
				response, _ := json.Marshal(responseJSON)

				w.Header().
				Set("Access-Control-Allow-Origin", "http://catalog.cc")
				w.Header().
				Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				w.Write(response)
			case "add_category":
				log.Println("add_category")

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
				log.Println("delete_category")

				responseJSON.Succes = true
				responseJSON.Data = data
		
				response, _ := json.Marshal(responseJSON)

				w.Header().
				Set("Access-Control-Allow-Origin", "http://catalog.cc")
				w.Header().
				Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				w.Write(response)
			case "new_jwt":
				log.Println("new_jwt")

				responseJSON.Succes = true
				responseJSON.Data = data
	
				response, _ := json.Marshal(responseJSON)

				w.Header().
				Set("Access-Control-Allow-Origin", "http://catalog.cc")
				w.Header().
				Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				w.Write(response)
			case "new_rjwt":
				log.Println("new_rjwt")

				responseJSON.Succes = true
				responseJSON.Data = data
		
				response, _ := json.Marshal(responseJSON)

				w.Header().
				Set("Access-Control-Allow-Origin", "http://catalog.cc")
				w.Header().
				Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				w.Write(response)
			default:
				log.Println("default")
				w.Header().
				Set("Access-Control-Allow-Origin", "http://catalog.cc")
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte("smth"))
			}
		}
	} else {
		log.Println("body is empty")

		w.Header().
		Set("Access-Control-Allow-Origin", "http://catalog.cc")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("smth"))
	}
}