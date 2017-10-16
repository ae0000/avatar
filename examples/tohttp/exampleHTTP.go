package main

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/ae0000/avatar"

	"github.com/go-chi/chi"
)

func main() {
	r := chi.NewRouter()

	// Get the png based on the initials, You would use it like this:
	//    <img src="http://localhost:3000/avatar/ae/png" width="150">
	r.Get("/avatar/{initials}.png", func(w http.ResponseWriter, r *http.Request) {
		initials := chi.URLParam(r, "initials")

		avatar.ToHTTP(initials, w)
	})

	// Show a heap of avatars
	r.Get("/showoff", func(w http.ResponseWriter, r *http.Request) {
		abc := strings.Repeat("ABCDEFGHIJKLMNOPQRSTUVQXYZ", 10)
		image := `<img src="http://localhost:3000/avatar/%s%s.png" width="100" height="100">`
		html := "<body>"

		for _, c := range abc {
			html += fmt.Sprintf(image, string(c), string(c))
		}

		html += "</body>"
		w.Write([]byte(html))
	})

	http.ListenAndServe(":3000", r)
}
