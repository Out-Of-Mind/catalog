package routes

import "net/http"

func dashboardHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().
		Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("/dashboard"))
}
