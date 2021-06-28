package routes

import (
	vars "github.com/out-of-mind/catalog/variables"
	"github.com/out-of-mind/catalog/structures"
	"github.com/satori/uuid"

	"encoding/json"
	"html/template"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
	"sort"
	"log"
)

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	c, err := r.Cookie("session_token")
    sessionToken := c.Value
    expiresAt := c.Expires

    if time.Now().Sub(expiresAt).Hours() <= 168 {
    	vars.Cache.Del(vars.CTX, sessionToken)
    	log.Println("setting new cookie")

    	sessionToken := uuid.NewV4().String()

		_, err = vars.Cache.Set(vars.CTX, sessionToken, "k", 720 * time.Hour).Result()
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		http.SetCookie(w, &http.Cookie{
			Name:    "session_token",
			Value:   sessionToken,
			Expires: time.Now().Add(720 * time.Hour),
		})
    }

	userId, _ := vars.Cache.Get(vars.CTX, sessionToken).Result()
	log.Println(userId)
	
	tmpl, err := template.ParseFiles(vars.TemplateDir+"index.html")
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("500 Internal Server Error"))
		return
	}

	rows, err := vars.DB.Query("SELECT categories.category_name, categories.category_id FROM categories, users WHERE users.user_id=$1 AND categories.group_id=users.group_id", userId)
    if err != nil {
        log.Println(err)
    }
    defer rows.Close()

    categoriesMap := make(map[int64]string)

    for rows.Next() {
    	var (
        	categoryName string
        	categoryId int64
    	)
        err = rows.Scan(&categoryName, &categoryId)
        if err != nil {
            log.Println(err)
        }
        categoriesMap[categoryId] = categoryName
    }

    rows, err = vars.DB.Query("SELECT items.item_name, items.category_id FROM items, categories, users WHERE users.user_id=$1 AND categories.group_id=users.group_id AND items.category_id=categories.category_id", userId)
	if err != nil {
        log.Println(err)
    }
    defer rows.Close()

    itemsMap := make(map[int64][]string)

    for rows.Next() {
    	var (
        	itemName string
        	categoryId int64
    	)
        err = rows.Scan(&itemName, &categoryId)
        if err != nil {
            log.Println(err)
        }
        itemsMap[categoryId] = append(itemsMap[categoryId], itemName)
    }

    var indexItems structures.IndexItems

    for id := range categoriesMap {
    	var indexData structures.IndexData

    	indexData.ID = id
    	indexData.CategoryName = categoriesMap[id]
    	indexData.CategoryID = strings.ReplaceAll(strings.ToLower(categoriesMap[id]), " ", "_")
    	for _, itemName := range itemsMap[id] {
    		indexData.ItemNames = append(indexData.ItemNames, itemName)
    	}
    	indexItems.Items = append(indexItems.Items, indexData)
    }

    sort.Sort(structures.ByID(indexItems.Items))
	
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