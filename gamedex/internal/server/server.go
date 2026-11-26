package server

import (
	"cmp"
	"html/template"
	"log"
	"net/http"
	"slices"
	"strconv"
	"strings"

	"gamedex/internal/api"
	"gamedex/internal/catalog"
	"gamedex/internal/store"
)

var tmpl = template.Must(template.ParseGlob("web/templates/*.html"))

var genres = []string{"plateforme", "rogue-like", "sport", "combat", "simulation", "moba", "fps", "indé", "action", "stratégie", "tactique", "course", "party", "metroidvania", "sandbox", "survie"}

type Server struct {
	Games       []api.Game
	Tournaments []api.Tournament
}

type ListPage struct {
	Games  []api.Game
	Query  string
	Genre  string
	Genres []string
	Likes  map[int]int
}

type GamePage struct {
	Game        api.Game
	Tournaments []api.Tournament
	Likes       int
	Players     map[int][]string
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
	mux.HandleFunc("POST /jeux/{id}/like", s.like)
	mux.HandleFunc("POST /tournois/{id}/inscription", s.inscription)

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
	genre := r.URL.Query().Get("genre")

	games, err := api.FetchGames("http://localhost:8081", q, genre)
	if err != nil {
		showError(w, http.StatusBadGateway, "La mini API ne répond pas.")
		return
	}

	// du plus liké au moins liké
	slices.SortStableFunc(games, func(a, b api.Game) int {
		return cmp.Compare(store.GetLikes(b.ID), store.GetLikes(a.ID))
	})

	likes := map[int]int{}
	for _, g := range games {
		likes[g.ID] = store.GetLikes(g.ID)
	}

	show(w, "liste.html", ListPage{
		Games:  games,
		Query:  q,
		Genre:  genre,
		Genres: genres,
		Likes:  likes,
	})
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

	tournaments := catalog.TournamentsOf(s.Tournaments, id)
	players := map[int][]string{}
	for _, t := range tournaments {
		players[t.ID] = store.Players(t.ID)
	}

	show(w, "jeu.html", GamePage{
		Game:        game,
		Tournaments: tournaments,
		Likes:       store.GetLikes(id),
		Players:     players,
	})
}

func (s *Server) like(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || id < 1 {
		showError(w, http.StatusBadRequest, "Identifiant invalide.")
		return
	}

	if _, ok := catalog.FindGame(s.Games, id); !ok {
		showError(w, http.StatusNotFound, "Ce jeu n'existe pas.")
		return
	}

	store.AddLike(id)
	http.Redirect(w, r, "/jeux/"+strconv.Itoa(id), http.StatusSeeOther)
}

func (s *Server) inscription(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || id < 1 {
		showError(w, http.StatusBadRequest, "Identifiant invalide.")
		return
	}

	gameID := 0
	found := false
	for _, t := range s.Tournaments {
		if t.ID == id {
			gameID = t.GameID
			found = true
		}
	}
	if !found {
		showError(w, http.StatusNotFound, "Ce tournoi n'existe pas.")
		return
	}

	pseudo := strings.TrimSpace(r.FormValue("pseudo"))
	if pseudo == "" || len(pseudo) > 20 {
		showError(w, http.StatusBadRequest, "Le pseudo doit faire entre 1 et 20 caractères.")
		return
	}

	store.AddPlayer(id, pseudo)
	http.Redirect(w, r, "/jeux/"+strconv.Itoa(gameID), http.StatusSeeOther)
}