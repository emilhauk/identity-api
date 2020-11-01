package endpoint

import "net/http"

func WebHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	fs := http.FileServer(http.Dir("./web"))
	fs.ServeHTTP(w, r)
}
