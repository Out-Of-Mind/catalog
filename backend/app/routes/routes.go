package routes

import "net/http"

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	homeHandler(w, r)
}

func DashboardHandler(w http.ResponseWriter, r *http.Request) {
	dashboardHandler(w, r)
}

func SelectHandler(w http.ResponseWriter, r *http.Request) {
	selectHandler(w, r)
}

func InviteHandler(w http.ResponseWriter, r *http.Request) {
	inviteHandler(w, r)
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	loginHandler(w, r)
}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	logoutHandler(w, r)
}

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	registerHandler(w, r)
}

func APIHandler(w http.ResponseWriter, r *http.Request) {
	apiHandler(w, r)
}
