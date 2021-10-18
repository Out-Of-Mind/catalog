package routes

import (
	"github.com/out-of-mind/catalog/structures"
	vars "github.com/out-of-mind/catalog/variables"
	"github.com/satori/uuid"

	"net/http"
	"strconv"
	"time"
)

func dashboardHandler(w http.ResponseWriter, r *http.Request) {
	c, _ := r.Cookie("session_token")
	sessionToken := c.Value

	userId, _ := vars.Cache.Get(vars.CTX, sessionToken).Result()
	vars.Log.Debug("user_id: ", userId)

	vars.Cache.Del(vars.CTX, sessionToken)
	vars.Log.Debug("setting new cookie")

	sessionToken = uuid.NewV4().String()

	_, err := vars.Cache.Set(vars.CTX, sessionToken, userId, 720*time.Hour).Result()
	if err != nil {
		vars.Log.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("500 Internal Server Error"))
		return
	}
	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    sessionToken,
		Domain:   ".catalog.cc",
		Expires:  time.Now().Add(720 * time.Hour),
		HttpOnly: true,
		Secure: true,
		SameSite: 2,
	})

	if err != nil {
		vars.Log.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("500 Internal Server Error"))
		return
	}

	userIdInt, err := strconv.Atoi(userId)
	if err != nil {
		vars.Log.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("500 Internal Server Error"))
		return
	}

	var dashboardData structures.DashboardData
	rows, err := vars.DB.Query("SELECT groups.group_name, groups.select_link FROM groups_users, groups WHERE groups_users.group_id=groups.group_id AND groups_users.user_id=$1", userIdInt)
	if err != nil {
		vars.Log.Error(err)
	}
	defer rows.Close()

	for rows.Next() {
		var dashboardGroup structures.DashboardGroup
		err = rows.Scan(&dashboardGroup.GroupName, &dashboardGroup.GroupWelcomeLink)
		if err != nil {
			vars.Log.Println(err)
		}
		dashboardData.Groups = append(dashboardData.Groups, dashboardGroup)
	}

	rows, err = vars.DB.Query("SELECT groups.group_name, groups.invite_link FROM owned_groups, groups WHERE owned_groups.group_id=groups.group_id AND owned_groups.user_id=$1", userIdInt)
	if err != nil {
		vars.Log.Error(err)
	}
	defer rows.Close()

	for rows.Next() {
		var dashboardOwnedGroup structures.DashboardOwnedGroup
		err = rows.Scan(&dashboardOwnedGroup.GroupName, &dashboardOwnedGroup.GroupWelcomeLink)
		if err != nil {
			vars.Log.Println(err)
		}
		dashboardData.OwnedGroups = append(dashboardData.OwnedGroups, dashboardOwnedGroup)
	}

	err = vars.DB.QueryRow("SELECT users.user_name FROM users WHERE users.user_id=$1", userIdInt).Scan(&dashboardData.UserName)
	if err != nil {
		vars.Log.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("500 Internal Server Error"))
		return
	}

	var (
		jwt, rjwt structures.JWT

		jwtHeader, rjwtHeader   structures.JWTHeader
		jwtPayload, rjwtPayload structures.JWTPayload
	)

	jwtHeader.Alg = "HS256"
	jwtHeader.Type = "JWT"

	rjwtHeader = jwtHeader

	jwtPayload.Exp = time.Now().Add(15 * time.Minute)
	jwtPayload.Value = userIdInt

	rjwtPayload.Exp = time.Now().Add(24 * 7 * time.Hour)
	rjwtPayload.Value = userIdInt

	jwt.Header = jwtHeader
	jwt.Payload = jwtPayload

	rjwt.Header = rjwtHeader
	rjwt.Payload = rjwtPayload

	jwtStr, err := newJWT(jwt)
	rjwtStr, err := newJWT(rjwt)

	dashboardData.JWT = jwtStr
	dashboardData.RJWT = rjwtStr

	showHTML(w, "dashboard.html", dashboardData)
}
