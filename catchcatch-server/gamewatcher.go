package main

import (
	"context"
	"log"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/perenecabuto/CatchCatch/catchcatch-server/protobuf"
)

const (
	// DefaultWatcherRange set the watcher radar radius size
	DefaultWatcherRange = 5000
	// MinPlayersPerGame ...
	MinPlayersPerGame = 3
	// DefaultGameDuration ...
	DefaultGameDuration = time.Minute
)

// GameContext stores game and its canel (and stop eventualy) function
type GameContext struct {
	game   *Game
	cancel context.CancelFunc
}

// GameWatcher is made to start/stop games by player presence
// and notify players events to each game by geo position
type GameWatcher struct {
	games  map[string]*GameContext
	wss    *WSServer
	stream EventStream
	Clear  context.CancelFunc
}

// NewGameWatcher builds GameWatecher
func NewGameWatcher(stream EventStream, wss *WSServer) *GameWatcher {
	return &GameWatcher{make(map[string]*GameContext), wss, stream, func() {}}
}

// observeGamePlayers events
// TODO: monitor game player watches
func (gw *GameWatcher) observeGamePlayers(ctx context.Context, g *Game) error {
	// TODO: change game status to running and set the server name
	// TODO: check if the game is already started
	// TODO: stop this when game stops
	// TODO: monitor game start
	return gw.stream.StreamIntersects(ctx, "player", "geofences", g.ID, func(d *Detection) error {
		switch d.Intersects {
		case Enter:
			if err := g.SetPlayer(d.FeatID, d.Lat, d.Lon); err != nil {
				return err
			}
		case Inside:
			if err := g.SetPlayer(d.FeatID, d.Lat, d.Lon); err != nil {
				return err
			}
		case Exit:
			g.RemovePlayer(d.FeatID)
		}
		return nil
	})
}

// WatchGames starts this gamewatcher to listen to player events over games
// TODO: monitor game watches
// TODO: create a timer to anounce the server as a game idle watcher
// TODO: before starting idle game watcher verify if the server watcher is the same of this server
// TODO: check if the server is the watcher and isn't running
func (gw *GameWatcher) WatchGames(ctx context.Context) error {
	var watcherCtx context.Context
	// TODO: remove this cancel
	watcherCtx, gw.Clear = context.WithCancel(ctx)
	defer gw.Clear()

	// TODO: find games idle
	return gw.stream.StreamNearByEvents(watcherCtx, "player", "geofences", "status = idle", 0, func(d *Detection) error {
		gameID := d.NearByFeatID
		if gameID == "" {
			return nil
		}

		go func() {
			if err := gw.watchGame(watcherCtx, gameID); err != nil {
				log.Println(err)
				// TODO: remove this cancel
				gw.games[gameID].cancel()
			}
		}()
		return nil
	})
}

// WatchGamesForever restart game wachter util context done
func (gw *GameWatcher) WatchGamesForever(ctx context.Context) error {
	done := ctx.Done()
	for {
		select {
		case <-done:
			return nil
		default:
			if err := gw.WatchGames(ctx); err != nil {
				return err
			}
		}
	}
}

func (gw *GameWatcher) watchGame(ctx context.Context, gameID string) error {
	_, exists := gw.games[gameID]
	if exists {
		return nil
	}
	g := NewGame(gameID, DefaultGameDuration, gw)
	gCtx, cancel := context.WithCancel(ctx)
	gw.games[gameID] = &GameContext{game: g, cancel: func() {
		delete(gw.games, gameID)
		cancel()
	}}

	errChan := make(chan error)
	go func() {
		errChan <- gw.observeGamePlayers(gCtx, g)
	}()
	go func() {
		if err := gw.startGameWhenReady(gCtx, g); err != nil {
			errChan <- err
		}
	}()
	if err := <-errChan; err != nil {
		gw.games[gameID].cancel()
		return err
	}
	return nil
}

func (gw *GameWatcher) startGameWhenReady(ctx context.Context, g *Game) error {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return nil
		case <-ticker.C:
			ready := len(g.players) >= MinPlayersPerGame
			if ready {
				return g.Start(ctx)
			}
		}
	}
}

// game callbacks

// OnGameStarted implements GameEvent.OnGameStarted
func (gw *GameWatcher) OnGameStarted(g Game, p GamePlayer) {
	// TODO: set game started
	gw.wss.Emit(p.ID, &protobuf.GameInfo{
		EventName: proto.String("game:started"),
		Id:        &g.ID,
		Game:      &g.ID, Role: proto.String(string(p.Role))})
}

// OnTargetWin implements GameEvent.OnTargetWin
func (gw *GameWatcher) OnTargetWin(p GamePlayer) {
	gw.wss.Emit(p.ID, &protobuf.Simple{EventName: proto.String("game:target:win")})
}

// OnGameFinish implements GameEvent.OnGameFinish
func (gw *GameWatcher) OnGameFinish(rank GameRank) {
	// TODO remove this game cancel and set game stopped
	log.Printf("gamewatcher:stop:game:%s", rank.Game)
	if game := gw.games[rank.Game]; game != nil && game.cancel != nil {
		game.cancel()
	}

	playersRank := make([]*protobuf.PlayerRank, len(rank.PlayerRank))
	for i, pr := range rank.PlayerRank {
		playersRank[i] = &protobuf.PlayerRank{Player: &pr.Player, Points: proto.Int32(int32(pr.Points))}
	}
	gw.wss.BroadcastTo(rank.PlayerIDs, &protobuf.GameRank{
		EventName: proto.String("game:finish"),
		Id:        &rank.Game,
		Game:      &rank.Game, PlayersRank: playersRank,
	})
}

// OnPlayerLoose implements GameEvent.OnPlayerLoose
func (gw *GameWatcher) OnPlayerLoose(g Game, p GamePlayer) {
	gw.wss.Emit(p.ID, &protobuf.Simple{EventName: proto.String("game:loose"), Id: &g.ID})
}

// OnTargetReached implements GameEvent.OnTargetReached
func (gw *GameWatcher) OnTargetReached(p GamePlayer) {
	gw.wss.Emit(p.ID, &protobuf.Distance{EventName: proto.String("game:target:reached"),
		Dist: &p.DistToTarget})
}

// OnPlayerNearToTarget implements GameEvent.OnPlayerNearToTarget
func (gw *GameWatcher) OnPlayerNearToTarget(p GamePlayer) {
	gw.wss.Emit(p.ID, &protobuf.Distance{EventName: proto.String("game:target:near"),
		Dist: &p.DistToTarget})
}
