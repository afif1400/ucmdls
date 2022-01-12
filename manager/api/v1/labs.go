package v1

import "net/http"

func HandleLabs() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("labs"))
	}
}
