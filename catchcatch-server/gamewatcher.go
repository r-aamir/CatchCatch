package main

import (
	"log"
	"time"

	io "github.com/googollee/go-socket.io"
)

type GameWatcher struct {
	games    map[string]*Game
	sessions *SessionManager
	stream   EventStream
}

func NewGameWatcher(stream EventStream, sessions *SessionManager) *GameWatcher {
	return &GameWatcher{make(map[string]*Game), sessions, stream}
}

// WatchPlayers events
func (gw *GameWatcher) WatchPlayers(g *Game) error {
	return gw.stream.StreamIntersects("player", "geofences", g.ID, func(d *Detection) {
		p := &Player{ID: d.FeatID, Lat: d.Lat, Lon: d.Lon}
		switch d.Intersects {
		case Enter:
			g.SetPlayerUntilReady(p, gw.sessions)
		case Exit:
			g.RemovePlayer(p, gw.sessions)
		case Inside:
			if !g.started {
				g.SetPlayerUntilReady(p, gw.sessions)
			} else if g.HasPlayer(p.ID) {
				g.UpdateAndNofityPlayer(p, gw.sessions)
			}
		}
	})
}

func (gw *GameWatcher) WatchGames() {
	err := gw.stream.StreamNearByEvents("player", "geofences", 0, func(d *Detection) {
		gameID := d.NearByFeatID
		game, exists := gw.games[gameID]
		if !exists {
			log.Println("Creating game", gameID)
			gameDuration := time.Minute
			game = NewGame(gameID, gameDuration)
			gw.games[gameID] = game

			go func() {
				if err := gw.WatchPlayers(game); err != nil {
					delete(gw.games, gameID)
					log.Printf("Error to start gamewatcher:%s - err: %v", game.ID, err)
				}
			}()
		}
	})
	if err != nil {
		log.Println("Error to stream geofence:event", err)
	}
}

func (gw *GameWatcher) WatchCheckpoints(server *io.Server) {
	err := gw.stream.StreamNearByEvents("player", "checkpoint", 1000, func(d *Detection) {
		if err := gw.sessions.Emit(d.FeatID, "checkpoint:detected", d); err != nil {
			log.Println("Error to notify player", d.FeatID, err)
		}
		server.BroadcastTo("main", "admin:feature:checkpoint", d)
	})
	if err != nil {
		log.Println("Error to stream geofence:event", err)
	}
}
