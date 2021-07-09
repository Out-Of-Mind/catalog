package routes

import (
	"github.com/out-of-mind/catalog/structures"
	vars "github.com/out-of-mind/catalog/variables"
	"github.com/satori/uuid"

	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"
)

func homeHandler(w http.ResponseWriter, r *http.Request) {
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

	rows, err := vars.DB.Query("SELECT categories.category_name, categories.category_id FROM categories, users WHERE users.user_id=$1 AND categories.group_id=users.group_id", userIdInt)
	if err != nil {
		vars.Log.Error(err)
	}
	defer rows.Close()

	categoriesMap := make(map[int64]string)

	for rows.Next() {
		var (
			categoryName string
			categoryId   int64
		)
		err = rows.Scan(&categoryName, &categoryId)
		if err != nil {
			vars.Log.Error(err)
		}
		categoriesMap[categoryId] = categoryName
	}

	rows, err = vars.DB.Query("SELECT items.item_name, items.category_id FROM items, categories, users WHERE users.user_id=$1 AND categories.group_id=users.group_id AND items.category_id=categories.category_id", userIdInt)
	if err != nil {
		vars.Log.Error(err)
	}
	defer rows.Close()

	itemsMap := make(map[int64][]string)

	for rows.Next() {
		var (
			itemName   string
			categoryId int64
		)
		err = rows.Scan(&itemName, &categoryId)
		if err != nil {
			vars.Log.Println(err)
		}
		itemsMap[categoryId] = append(itemsMap[categoryId], itemName)
	}

	var indexItems structures.IndexItems

	for id := range categoriesMap {
		var indexData structures.IndexData

		indexData.ID = id
		indexData.CategoryName = categoriesMap[id]
		indexData.CategoryID = strings.ReplaceAll(categoriesMap[id], " ", "_")
		for _, itemName := range itemsMap[id] {
			indexData.ItemNames = append(indexData.ItemNames, itemName)
		}
		indexItems.Items = append(indexItems.Items, indexData)
	}

	sort.Sort(structures.ByID(indexItems.Items))

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

	indexItems.JWT = jwtStr
	indexItems.RJWT = rjwtStr

	showHTML(w, "index.html", indexItems)
}
