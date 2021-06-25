package main

import (
	vars "github.com/out-of-mind/catalog/variables"
	"github.com/out-of-mind/catalog/middlewares"
	"github.com/out-of-mind/catalog/routes"
	"github.com/go-redis/redis/v8"
	"github.com/gorilla/mux"
	
	"database/sql"
    _ "github.com/lib/pq"
	"os/signal"
	"net/http"
	"context"
	"time"
	"log"
	"fmt"
	"os"
)

func main() {	
	initDB()
	defer vars.DB.Close()
	initCache()

	wait := time.Second * 10;

	r := mux.NewRouter()

	r.HandleFunc("/", routes.HomeHandler)
	r.HandleFunc("/login", routes.LoginHandler)
	r.HandleFunc("/api", routes.APIHandler)/*.
	Methods("POST")*/
	r.Use(middlewares.LoggingMiddleware)

	srv := &http.Server{
		Handler:      r,
		Addr:         "127.0.0.1:8080",
		ReadTimeout: 15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout: 20 * time.Second,
		ReadHeaderTimeout: 20 * time.Second,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Println(err)
		}
	}()

	log.Println(fmt.Sprintf("Server started at %s", "127.0.0.1:8080"))

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	<-c

	ctx, cancel := context.WithTimeout(context.Background(), wait)
	defer cancel()
	srv.Shutdown(ctx)
	log.Println("shutting down")
	os.Exit(0)
}

func initCache() {
	vars.Cache = redis.NewClient(&redis.Options{
        Addr:     "localhost:6379",
        Password: "",
        DB:       0,
    })
}

func initDB() {
	var err error
	connStr := "user=catalog_user password=password dbname=catalog_db sslmode=disable"
    vars.DB, err = sql.Open("postgres", connStr)
    if err != nil {
        log.Println(err)
        os.Exit(1)
    }
}