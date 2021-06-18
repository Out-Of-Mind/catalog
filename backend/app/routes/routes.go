package routes

import (
	"github.com/out-of-mind/catalog/structures"

	"encoding/json"
	"io/ioutil"
	"net/http"
	"log"
)

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().
	Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("hi"))
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().
	Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("hi"))
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