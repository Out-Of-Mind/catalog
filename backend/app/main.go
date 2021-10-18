package main

import (
	"github.com/go-redis/redis/v8"
	"github.com/out-of-mind/catalog/config"
	"github.com/out-of-mind/catalog/middlewares"
	"github.com/out-of-mind/catalog/routes"
	vars "github.com/out-of-mind/catalog/variables"
	"github.com/sirupsen/logrus"
	"github.com/t-tomalak/logrus-easy-formatter"

	"github.com/gorilla/mux"

	"context"
	"database/sql"
	"flag"
	"fmt"
	_ "github.com/lib/pq"
	"net/http"
	"os"
	"os/signal"
	"time"
)

var log *logrus.Logger
var c config.Config
var configFile string

func main() {
	flag.StringVar(&configFile, "-c", "/catalog/backend/app/config.json", "usage: -c ./config.json to set ./config.json as config file")
	flag.Parse()

	vars.Log = initLogger("catalog.log", "INFO")

	c = config.ParseConfig(configFile)
	initDB()
	defer vars.DB.Close()
	initCache()

	wait := time.Second * 10

	r := mux.NewRouter()

	r.HandleFunc("/", routes.HomeHandler)
	r.HandleFunc("/dashboard", routes.DashboardHandler)
	r.HandleFunc("/select/{id}", routes.SelectHandler)
	r.HandleFunc("/invite/{id}", routes.InviteHandler)
	r.HandleFunc("/login", routes.LoginHandler)
	r.HandleFunc("/register", routes.RegisterHandler)
	r.HandleFunc("/logout", routes.LogoutHandler)
	r.HandleFunc("/api", routes.APIHandler).
		Methods("POST").
		Host("catalog.cc").
		Host("dashboard.catalog.cc")

	r.Use(middlewares.CSRFMiddleware)
	r.Use(middlewares.LoggingMiddleware)

	srv := &http.Server{
		Handler:           r,
		Addr:              "localhost:8080",
		ReadTimeout:       15 * time.Second,
		WriteTimeout:      15 * time.Second,
		IdleTimeout:       20 * time.Second,
		ReadHeaderTimeout: 20 * time.Second,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			fmt.Println(err)
		}
	}()

	fmt.Printf("Server started at %s", "localhost:8080\n")
	vars.Log.Println(fmt.Sprintf("Server started at %s", "localhost:8080"))

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	<-c

	ctx, cancel := context.WithTimeout(context.Background(), wait)
	defer cancel()
	srv.Shutdown(ctx)
	fmt.Println("shutting down")
	os.Exit(0)
}

func initCache() {
	vars.Cache = redis.NewClient(&redis.Options{
		Addr:     c.REDIS_IP + ":" + c.REDIS_PORT,
		Password: c.REDIS_PASSWORD,
		DB:       c.REDIS_DB,
	})

	_, err := vars.Cache.Do(vars.CTX, "keys", "*").Result()
	if err != nil {
		fmt.Println("Cannot access redis, exit with error: ", err)
		vars.Log.Fatal("Cannot access redis, exit with error: ", err)
	}
}

func initDB() {
	var err error
	connStr := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=%s", c.DB_USER, c.DB_PASSOWRD, c.DB_NAME, c.DB_SSLMODE)
	vars.DB, err = sql.Open("postgres", connStr)
	if err != nil {
		fmt.Println("Cannot access postgresql db, exit with error: ", err)
		vars.Log.Fatal("Cannot access postgresql db, exit with error: ", err)
	}

	err = vars.DB.Ping()
	if err != nil {
		fmt.Println("Cannot access postgresql db, exit with error: ", err)
		vars.Log.Fatal("Cannot access postgresql db, exit with error: ", err)
	}

	vars.DB.SetMaxOpenConns(30)
	vars.DB.SetMaxIdleConns(30)
	vars.DB.SetConnMaxLifetime(15*time.Minute)
}

func initLogger(path, logLevel string) *logrus.Logger {
	logFile, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println(err)
		vars.Log.Fatal(err)
	}
	lvl, err := logrus.ParseLevel(logLevel)
	if err != nil {
		fmt.Println(err)
		vars.Log.Fatal(err)
	}
	logger := &logrus.Logger{
		Out:   logFile,
		Level: lvl,
		Formatter: &easy.Formatter{
			TimestampFormat: "2006-01-02 15:04:05",
			LogFormat:       "[%lvl%]: %time% - %msg%\n",
		},
	}

	return logger
}
