package routes

import (
	vars "github.com/out-of-mind/catalog/variables"

	"github.com/gorilla/mux"

	"net/http"
	"strconv"
)

func selectHandler(w http.ResponseWriter, r *http.Request) {
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

	id := mux.Vars(r)["id"]
	_, err = vars.DB.Exec("UPDATE users SET group_id=(SELECT groups.group_id FROM groups WHERE groups.select_link=$1) WHERE user_id=$2", id, userIdInt)
	if err != nil {
		vars.Log.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("500 Internal Server Error"))
		return
	}

	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
}
