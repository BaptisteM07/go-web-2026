package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
	"net/url"
)

const baseURL = "http://localhost:8081"

var client = &http.Client{Timeout: 3 * time.Second}

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

func FetchGames(baseURL, q, genre string) ([]Game, error) {
	params := url.Values{}
	params.Set("q", q)
	params.Set("genre", genre)

resp, err := http.Get(baseURL + "/api/games?" + params.Encode())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API: statut %d", resp.StatusCode)
	}

	var games []Game
	if err := json.NewDecoder(resp.Body).Decode(&games); err != nil {
		return nil, err
	}
	return games, nil
}

func FetchTournaments() ([]Tournament, error) {
	resp, err := client.Get(baseURL + "/api/tournaments")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("l'API a répondu %d", resp.StatusCode)
	}

	var tournaments []Tournament
	err = json.NewDecoder(resp.Body).Decode(&tournaments)
	return tournaments, err
}