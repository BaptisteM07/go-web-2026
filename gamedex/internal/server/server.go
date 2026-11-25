package server

import (
	"html/template"
	"log"
	"net/http"
	"strconv"

	"gamedex/internal/api"
	"gamedex/internal/catalog"
)

var tmpl = template.Must(template.ParseGlob("web/templates/*.html"))

type Server struct {
	Games       []api.Game
	Tournaments []api.Tournament
}

type ListPage struct {
	Games []api.Game
	Query string
}

type GamePage struct {
	Game        api.Game
	Tournaments []api.Tournament
}

type ErrorPage struct {
	Code    int
	Message string
}

func (s *Server) Routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /{$}", s.accueil)
	mux.HandleFunc("GET /jeux", s.list)
	mux.HandleFunc("GET /jeux/{id}", s.detail)

	fs := http.FileServer(http.Dir("web/static"))
	mux.Handle("GET /static/", http.StripPrefix("/static/", fs))
	return mux
}

func show(w http.ResponseWriter, name string, data any) {
	err := tmpl.ExecuteTemplate(w, name, data)
	if err != nil {
		log.Println(err)
	}
}

func showError(w http.ResponseWriter, code int, message string) {
	w.WriteHeader(code)
	show(w, "erreur.html", ErrorPage{Code: code, Message: message})
}

func (s *Server) accueil(w http.ResponseWriter, r *http.Request) {
	show(w, "accueil.html", nil)
}

func (s *Server) list(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query().Get("q")
	show(w, "liste.html", ListPage{Games: s.Games, Query: q})
}

func (s *Server) detail(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || id < 1 {
		showError(w, http.StatusBadRequest, "Identifiant invalide.")
		return
	}

	game, ok := catalog.FindGame(s.Games, id)
	if !ok {
		showError(w, http.StatusNotFound, "Ce jeu n'existe pas.")
		return
	}

	show(w, "jeu.html", GamePage{
		Game:        game,
		Tournaments: catalog.TournamentsOf(s.Tournaments, id),
	})
}