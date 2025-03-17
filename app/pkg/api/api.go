package api

import (
	"net/http"
	"url_shortener/pkg/db"

	"github.com/gorilla/mux"
)

type API struct {
	hostname string
	router   *mux.Router
	DB       db.Interface
}

func New(hostname string, db db.Interface) *API {
	api := API{hostname: hostname, router: mux.NewRouter(), DB: db}
	api.endpoints()
	return &api
}

func (api *API) Router() *mux.Router {
	return api.router
}

func (api *API) endpoints() {
	api.router.HandleFunc("/", api.saveURL).Methods(http.MethodPost)
	api.router.HandleFunc("/{shortURL}", api.getLongURL).Methods(http.MethodGet)
}

func (api *API) saveURL(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Form parsing error: "+err.Error(), http.StatusBadRequest)
		return
	}
	url := r.FormValue("url")
	shortURL, err := api.DB.MakeShort(url)
	if err != nil {
		http.Error(w, "Database error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write([]byte(api.hostname + "/" + shortURL))
}

func (api *API) getLongURL(w http.ResponseWriter, r *http.Request) {
	shortURL := mux.Vars(r)["shortURL"]
	if len(shortURL) != int(api.DB.URLSize()) {
		http.Error(w, "Wrong short url", http.StatusBadRequest)
		return
	}

	originalURL, err := api.DB.GetOriginal(shortURL)
	if err != nil {
		http.Error(w, "Database error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	if originalURL == "" {
		http.Error(w, "No such url", http.StatusBadRequest)
	}

	w.Write([]byte(originalURL))
}
