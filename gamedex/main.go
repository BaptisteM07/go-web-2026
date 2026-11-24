package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"strconv"
)

func accueil(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprint(w, `<!DOCTYPE html>
<html lang="fr">
<head>
  <meta charset="utf-8">
  <title>GameDex</title>
  <link rel="stylesheet" href="/static/css/style.css">
</head>
<body>
  <h1>Bienvenue sur GameDex</h1>
  <p>Un catalogue de jeux et de leurs tournois.</p>
  <p><a href="/jeux">Voir les jeux</a></p>
</body>
</html>`)
}

func listeJeux(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Liste des jeux")
}

func detailJeu(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || id < 1 {
		http.Error(w, "identifiant invalide", http.StatusBadRequest) // 400
		return
	}
	if id > 12 {
		http.NotFound(w, r) 
		return
	}
	fmt.Fprintf(w, "Fiche du jeu %d\n", id)
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

	fmt.Println("Serveur sur http://localhost:" + *port)
	log.Fatal(http.ListenAndServe(":"+*port, mux))
}
