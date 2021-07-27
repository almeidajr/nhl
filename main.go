package main

import (
	"io"
	"log"
	"nhl/api"
	"os"
	"sync"
	"time"
)

func main() {
	now := time.Now()

	file, err := os.OpenFile("data.txt", os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		log.Fatal("error opening the file data.txt", err)
	}
	defer file.Close()

	wrt := io.MultiWriter(os.Stdout, file)
	log.SetOutput(wrt)

	teams, err := api.GetAllTeams()
	if err != nil {
		log.Fatal("error while getting all teams", err)
	}

	var waitGroup sync.WaitGroup
	waitGroup.Add(len(teams))

	// unbuffered channel
	results := make(chan []api.Roster)

	for _, team := range teams {
		go func(team api.Team) {
			roster, err := api.GetAllRosters(team.ID)
			if err != nil {
				log.Fatal("error while getting roster", err)
			}

			results <- roster

			waitGroup.Done()
		}(team)
	}

	go func() {
		waitGroup.Wait()
		close(results)
	}()

	display(results)

	log.Printf("took %v", time.Now().Sub(now).String())
}

func display(results chan []api.Roster) {
	for rosters := range results {
		for _, roster := range rosters {
			log.Printf("ID: %d\n", roster.Person.ID)
			log.Printf("Name: %s\n", roster.Person.FullName)
			log.Printf("Position: %s\n", roster.Position.Abbreviation)
			log.Printf("Jersey: %s", roster.JerseyNumber)
			log.Println("### ### ### ### ### ### ###")
		}
	}
}
