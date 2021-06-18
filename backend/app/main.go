package main

import (
	"github.com/out-of-mind/catalog/routes"
	"github.com/gorilla/mux"
	
	"os/signal"
	"net/http"
	"context"
	//"flag"
	"time"
	"log"
	"fmt"
	"os"
)

func main() {
	wait := time.Second * 10;

	r := mux.NewRouter()

	r.HandleFunc("/", routes.HomeHandler)
	r.HandleFunc("/login", routes.LoginHandler)
	r.HandleFunc("/api", routes.APIHandler)/*.
	Methods("POST")*/

	srv := &http.Server{
		Handler:      r,
		Addr:         "127.0.0.1:8000",
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

	log.Println(fmt.Sprintf("Server started at %s", "127.0.0.1:8000"))

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	<-c

	ctx, cancel := context.WithTimeout(context.Background(), wait)
	defer cancel()
	srv.Shutdown(ctx)
	log.Println("shutting down")
	os.Exit(0)
}