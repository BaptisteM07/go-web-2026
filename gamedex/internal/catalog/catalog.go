package catalog

import (
	"slices"
	"strings"
	"gamedex/internal/api"
)

func FindGame(games []api.Game, id int) (api.Game, bool) {
	for _, g := range games {
		if g.ID == id {
			return g, true
		}
	}
	return api.Game{}, false
}

func TournamentsOf(tournaments []api.Tournament, gameID int) []api.Tournament {
	result := []api.Tournament{}
	for _, t := range tournaments {
		if t.GameID == gameID {
			result = append(result, t)
		}
	}
	return result
}

func Filter(games []api.Game, genre, query string) []api.Game {
	result := []api.Game{}
	for _, g := range games {
		if genre != "" && !slices.Contains(g.Genres, genre) {
			continue
		}
		if query != "" && !strings.Contains(strings.ToLower(g.Name), strings.ToLower(query)) {
			continue
		}
		result = append(result, g)
	}
	return result
}