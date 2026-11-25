package main

import (
	"flag"
	"fmt"
	"gamedex/internal/api"
	"gamedex/internal/server"
	"log"
	"net/http"
)

func main() {
	port := flag.String("port", "8080", "port d'écoute du serveur")
	flag.Parse()

	games, err := api.FetchGames()
	if err != nil {
		log.Fatal("Impossible de charger les jeux : ", err)
	}

	tournaments, err := api.FetchTournaments()
	if err != nil {
		log.Fatal("Impossible de charger les tournois : ", err)
	}

	srv := &server.Server{Games: games, Tournaments: tournaments}

	fmt.Println("GameDex sur http://localhost:" + *port)
	log.Fatal(http.ListenAndServe(":"+*port, srv.Routes()))
}
