package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"
	"encoding/json"
	"html/template"
)

const baseURL = "http://localhost:8081"

var client = &http.Client{Timeout: 3 * time.Second}
var tmpl = template.Must(template.ParseGlob("templates/*.html"))

type Game struct {
	ID     int      `json:"id"`
	Name   string   `json:"name"`
	Year   int      `json:"year"`
	Studio string   `json:"studio"`
	Genres []string `json:"genres"`
}

type Tournament struct {
	ID     int    `json:"id"`
	GameID int    `json:"gameId"`
	Name   string `json:"name"`
	City   string `json:"city"`
	Date   string `json:"date"`
}

type ListePage struct {
	Games []Game
	Query string
}

type JeuPage struct {
	Game        Game
	Tournaments []Tournament
}

type ErreurPage struct {
	Code    int
	Message string
}


func fetchJSON(path string, dest any) error {
	resp, err := client.Get(baseURL + path)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("l'API a répondu %d", resp.StatusCode)
	}

	return json.NewDecoder(resp.Body).Decode(dest)
}

func FetchGames() ([]Game, error) {
	var games []Game
	err := fetchJSON("/api/games", &games)
	return games, err
}

func FetchTournaments() ([]Tournament, error) {
	var tournaments []Tournament
	err := fetchJSON("/api/tournaments", &tournaments)
	return tournaments, err
}

func accueil(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprint(w, `<!DOCTYPE html>
<html lang="fr">
<head>
  <meta charset="utf-8">
  <meta name="viewport" content="width=device-width, initial-scale=1">
  <title>GameDex</title>
  <link rel="stylesheet" href="/static/css/style.css">
</head>
<body>
  <main>
    <h1>Bienvenue sur GameDex</h1>
    <p>Un catalogue de jeux et de leurs tournois.</p>
    <a class="bouton" href="/jeux">Voir les jeux</a>
  </main>
  <footer>GameDex · B1 Ynov Lyon</footer>
</body>
</html>`)
}

func listeJeux(w http.ResponseWriter, r *http.Request) {
	games, err := FetchGames()
	if err != nil {
		log.Println(err)
		http.Error(w, "Le service est indisponible", http.StatusBadGateway)
		return
	}

page := ListePage{Games: games, Query: r.URL.Query().Get("q")}
err = tmpl.ExecuteTemplate(w, "liste.html", page)
	if err != nil {
		log.Println(err)
	}
}

func detailJeu(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || id < 1 {
		http.Error(w, "Identifiant invalide", http.StatusBadRequest)
		return
	}
	games, err := FetchGames()
	if err != nil {
		log.Println(err)
		http.Error(w, "Le service est indisponible", http.StatusBadGateway)
		return
	}
	var game Game
	trouve := false
	for _, g := range games {
		if g.ID == id {
			game = g
			trouve = true
		}
	}
	if !trouve {
		http.Error(w, "Ce jeu n'existe pas", http.StatusNotFound)
		return
	}
	tournaments, err := FetchTournaments()
	if err != nil {
		log.Println(err)
		http.Error(w, "Le service est indisponible", http.StatusBadGateway)
		return
	}
	var deCeJeu []Tournament
	for _, t := range tournaments {
		if t.GameID == id {
			deCeJeu = append(deCeJeu, t)
		}
	}
	err = tmpl.ExecuteTemplate(w, "jeu.html", JeuPage{Game: game, Tournaments: deCeJeu})
	if err != nil {
		log.Println(err)
	}
}

func erreur(w http.ResponseWriter, code int, message string) {
	w.WriteHeader(code)
	err := tmpl.ExecuteTemplate(w, "erreur.html", ErreurPage{Code: code, Message: message})
	if err != nil {
		log.Println(err)
	}
}

func main() {
	port := flag.String("port", "8080", "port d'écoute du serveur")
	flag.Parse()

	mux := http.NewServeMux()
	mux.HandleFunc("GET /{$}", accueil)
	mux.HandleFunc("GET /jeux", listeJeux)
	mux.HandleFunc("GET /jeux/{id}", detailJeu)

	fs := http.FileServer(http.Dir("static"))
	mux.Handle("GET /static/", http.StripPrefix("/static/", fs))


	games, err := FetchGames()
	if err != nil {
		log.Fatal("Impossible de charger les jeux : ", err)
	}
	tournaments, err := FetchTournaments()
	if err != nil {
		log.Fatal("Impossible de charger les tournois : ", err)
	}
	byGame := make(map[int][]Tournament)
	for _, t := range tournaments {
		byGame[t.GameID] = append(byGame[t.GameID], t)
	}
	for _, g := range games {
		fmt.Printf("%-20s (%d)  %d tournoi(s)\n", g.Name, g.Year, len(byGame[g.ID]))
}

	fmt.Println("Serveur sur http://localhost:" + *port)
	log.Fatal(http.ListenAndServe(":"+*port, mux))
}