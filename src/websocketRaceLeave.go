package main

import (
	"strconv"

	"github.com/Zamiell/isaac-racing-server/src/log"

	melody "gopkg.in/olahol/melody.v1"
)

func websocketRaceLeave(s *melody.Session, d *IncomingWebsocketData) {
	// Local variables
	username := d.v.Username

	/*
		Validation
	*/

	// Validate that the race exists
	var race *Race
	if v, ok := races[d.ID]; !ok {
		return
	} else {
		race = v
	}

	// Validate that the race has started
	if race.Status != "open" {
		return
	}

	// Validate that they are in the race
	if _, ok := race.Racers[username]; !ok {
		return
	}

	/*
		Leave
	*/

	// Disconnect the user from the channel for that race
	d.Room = "_race_" + strconv.Itoa(race.ID)
	websocketRoomLeaveSub(s, d)

	// Remove this racer from the map
	delete(race.Racers, username)

	// Send everyone a notification that the user left the race
	for _, s := range websocketSessions {
		type RaceLeftMessage struct {
			ID   int    `json:"id"`
			Name string `json:"name"`
		}
		websocketEmit(s, "raceLeft", &RaceLeftMessage{
			race.ID,
			username,
		})
	}

	if len(race.Racers) == 0 {
		// Remove this race if this is the last person to leave
		delete(races, d.ID)

		// Also delete it from the database
		if err := db.Races.Delete(d.ID); err != nil {
			log.Error("Database error:", err)
			websocketError(s, d.Command, "")
			return
		}
	} else if len(race.Racers) == 1 {
		// If the race went from 2 people to 1, check to see if the last person is ready
		for _, lastRacer := range race.Racers {
			if lastRacer.Status == "ready" {
				// Automatically unready the last person so that they don't start the race by themsevles
				race.SetRacerStatus(lastRacer.Name, "not ready")
			}
		}
	} else {
		race.CheckStart()
	}
}