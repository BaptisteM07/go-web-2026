package store

var inscrits = map[int][]string{} // id du tournoi -> pseudos

func AddPlayer(tournamentID int, pseudo string) {
	mu.Lock()
	defer mu.Unlock()
	inscrits[tournamentID] = append(inscrits[tournamentID], pseudo)
}

func Players(tournamentID int) []string {
	mu.Lock()
	defer mu.Unlock()
	return append([]string(nil), inscrits[tournamentID]...)
}