package routes

import (
	vars "github.com/out-of-mind/catalog/variables"

	"github.com/gorilla/mux"

	"context"
	"net/http"
	"strconv"
	"strings"
	"database/sql"
	_ "github.com/lib/pq"
)

func inviteHandler(w http.ResponseWriter, r *http.Request) {
	c, _ := r.Cookie("session_token")
	sessionToken := c.Value

	userId, _ := vars.Cache.Get(vars.CTX, sessionToken).Result()
	userIdInt, err := strconv.Atoi(userId)
	if err != nil {
		vars.Log.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("500 Internal Server Error"))
		return
	}

	var groupId int
	id := mux.Vars(r)["id"]
	err = vars.DB.QueryRow("SELECT groups.group_id FROM groups WHERE groups.invite_link=$1", id).Scan(&groupId)
	if err != nil {
		vars.Log.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("500 Internal Server Error"))
		return
	}

	txCtx := context.Background()

	tx, err := vars.DB.BeginTx(txCtx, &sql.TxOptions{Isolation: sql.LevelReadCommitted})

	_, err = vars.DB.ExecContext(txCtx, "UPDATE users SET group_id=$1 WHERE user_id=$2", groupId, userIdInt)
	if err != nil {
		tx.Rollback()
		vars.Log.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("500 Internal Server Error"))
		return
	}

	_, err = vars.DB.ExecContext(txCtx, "INSERT INTO groups_users(user_id, group_id) VALUES($1, $2)",  userIdInt, groupId)
	if err != nil {
		if strings.Contains(err.Error(), "pk_group_user") {

		} else {
			tx.Rollback()
			vars.Log.Error(err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("500 Internal Server Error"))
			return
		}
	}

	if err := tx.Commit(); err != nil {
		vars.Log.Error(err)
		w.Header().
			Set("Access-Control-Allow-Origin", "http://catalog.cc")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("500 Internal Server Error"))
		return
	}

	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
}
