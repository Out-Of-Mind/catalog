package main

import (
	"github.com/out-of-mind/catalog/config"

	"database/sql"
	_ "github.com/lib/pq"

	"flag"
	"fmt"
	"log"
)

var (
	db         *sql.DB
	c          config.Config
	configFile string
)

func main() {
	flag.StringVar(&configFile, "-c", "/catalog/backend/app/config.json", "usage: -c ./config.json to set ./config.json as config file")
	flag.Parse()
	c = config.ParseConfig(configFile)

	initDB()

	log.Println("Started migrations")

	execSql("DROP TABLE IF EXISTS items")
	execSql("DROP TABLE IF EXISTS categories")
	execSql("DROP TABLE IF EXISTS groups_users")
	execSql("DROP TABLE IF EXISTS owned_groups")
	execSql("DROP TABLE IF EXISTS users")
	execSql("DROP TABLE IF EXISTS groups")

	execSql("CREATE TABLE groups(group_id INT GENERATED ALWAYS AS IDENTITY, group_name VARCHAR(255) NOT NULL, invite_link VARCHAR(32) NOT NULL, UNIQUE(invite_link), select_link VARCHAR(32) NOT NULL, UNIQUE(select_link), UNIQUE(group_name), PRIMARY KEY(group_id))")
	execSql("CREATE TABLE users(user_id INT GENERATED ALWAYS AS IDENTITY, user_name VARCHAR(255) NOT NULL, email TEXT NOT NULL, password TEXT NOT NULL, group_id INT,  PRIMARY KEY(user_id), UNIQUE(user_name), CONSTRAINT fk_group FOREIGN KEY(group_id) REFERENCES groups(group_id) ON DELETE SET NULL)")
	execSql("CREATE TABLE groups_users(user_id INT NOT NULL REFERENCES users(user_id) ON DELETE CASCADE, group_id INT NOT NULL REFERENCES groups(group_id) ON DELETE CASCADE, CONSTRAINT pk_group_user PRIMARY KEY (user_id, group_id))")
	execSql("CREATE TABLE owned_groups(user_id INT NOT NULL REFERENCES users(user_id) ON DELETE CASCADE, group_id INT NOT NULL REFERENCES groups(group_id) ON DELETE CASCADE, CONSTRAINT pk_owned_group PRIMARY KEY (user_id, group_id))")
	execSql("CREATE TABLE categories(category_id INT GENERATED ALWAYS AS IDENTITY, category_name TEXT NOT NULL, group_id INT NOT NULL, PRIMARY KEY(category_id), UNIQUE(category_name), CONSTRAINT fk_group FOREIGN KEY(group_id) REFERENCES groups(group_id) ON DELETE CASCADE)")
	execSql("CREATE TABLE items(item_id INT GENERATED ALWAYS AS IDENTITY, item_name TEXT NOT NULL, category_id INT NOT NULL, UNIQUE(item_name), PRIMARY KEY(item_id), CONSTRAINT fk_category FOREIGN KEY(category_id) REFERENCES categories(category_id) ON DELETE CASCADE)")

	execSql("CREATE INDEX users_id_idx ON users(user_id)")
	execSql("CREATE INDEX users_group_id_idx ON users(group_id)")
	execSql("CREATE INDEX items_category_id_idx ON items(category_id)")

	log.Println("Migrations finished!")
}

func initDB() {
	var err error
	connStr := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=%s", c.DB_USER, c.DB_PASSOWRD, c.DB_NAME, c.DB_SSLMODE)
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Cannot access postgresql db, exit with error: ", err)
	}
}

func execSql(sql string) {
	_, err := db.Exec(sql)
	if err != nil {
		log.Fatal(fmt.Sprintf("Operation %s exited with error: %s", sql, err))
	}
}
