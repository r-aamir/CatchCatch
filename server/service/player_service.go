package service

import (
	"context"
	"fmt"

	"github.com/perenecabuto/CatchCatch/server/model"
	"github.com/perenecabuto/CatchCatch/server/service/repository"
	gjson "github.com/tidwall/gjson"
)

// PlayerLocationService manage players and features
type PlayerLocationService interface {
	Set(p *model.Player) error
	Remove(p *model.Player) error
	All() (model.PlayerList, error)

	ObservePlayersAround(context.Context, PlayersAroundCallback) error
	ObservePlayerNearToFeature(context.Context, string, PlayerNearToFeatureCallback) error
}

type PlayerNearToFeatureCallback func(playerID string, distTo float64, f model.Feature) error
type PlayersAroundCallback func(playerID string, movingPlayer model.Player, exit bool) error

// Tile38PlayerLocationService manages player locations
type Tile38PlayerLocationService struct {
	repo   repository.Repository
	stream repository.EventStream
}

// NewPlayerLocationService build a PlayerLocationService
func NewPlayerLocationService(repo repository.Repository, stream repository.EventStream) PlayerLocationService {
	return &Tile38PlayerLocationService{repo, stream}
}

// Exists add new player
func (s *Tile38PlayerLocationService) Exists(p *model.Player) (bool, error) {
	return s.repo.Exists("player", p.ID)
}

// Set player data
func (s *Tile38PlayerLocationService) Set(p *model.Player) error {
	_, err := s.repo.SetFeature("player", p.ID,
		fmt.Sprintf(`{"type": "Point", "coordinates": [%f, %f]}`, p.Lon, p.Lat))
	return err
}

// Remove player
func (s *Tile38PlayerLocationService) Remove(p *model.Player) error {
	return s.repo.RemoveFeature("player", p.ID)
}

// All return all registered players
func (s *Tile38PlayerLocationService) All() (model.PlayerList, error) {
	features, err := s.repo.Features("player")
	if err != nil {
		return nil, err
	}
	list := make(model.PlayerList, len(features))
	for i, f := range features {
		coords := gjson.Get(f.Coordinates, "coordinates").Array()
		if len(coords) != 2 {
			coords = make([]gjson.Result, 2)
		}
		list[i] = &model.Player{ID: f.ID, Lat: coords[1].Float(), Lon: coords[0].Float()}
	}
	return list, nil
}

func (s *Tile38PlayerLocationService) ObservePlayersAround(ctx context.Context, callback PlayersAroundCallback) error {
	return s.stream.StreamNearByEvents(ctx, "player", "player", "*", DefaultGeoEventRange, func(d *repository.Detection) error {
		playerID := d.NearByFeatID
		movingPlayer := model.Player{ID: d.FeatID, Lon: d.Lon, Lat: d.Lat}
		return callback(playerID, movingPlayer, d.Intersects == repository.Exit)
	})
}

func (s *Tile38PlayerLocationService) ObservePlayerNearToFeature(ctx context.Context, group string, callback PlayerNearToFeatureCallback) error {
	return s.stream.StreamNearByEvents(ctx, group, "player", "*", DefaultGeoEventRange, func(d *repository.Detection) error {
		if d.Intersects == repository.Inside {
			playerID := d.NearByFeatID
			f := model.Feature{ID: d.FeatID, Group: group, Coordinates: d.Coordinates}
			return callback(playerID, d.NearByMeters, f)
		}
		return nil
	})
}