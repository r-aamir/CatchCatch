package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"strings"
	"time"

	redis "gopkg.in/redis.v5"
)

// Player payload
type Player struct {
	ID string  `json:"id"`
	X  float32 `json:"x"`
	Y  float32 `json:"y"`
}

func (p *Player) String() string {
	return fmt.Sprintln("id:", p.ID, "x:", p.X, "y:", p.Y)
}

// PlayerList payload for list of players
type PlayerList struct {
	Players []*Player `json:"players"`
}

// PlayerLocationService manages player locations
type PlayerLocationService struct {
	client *redis.Client
}

// Register add new player
func (s *PlayerLocationService) Register(p *Player) error {
	return s.Update(p)
}

// Update player data
func (s *PlayerLocationService) Update(p *Player) error {
	cmd := redis.NewStringCmd("SET", "player", p.ID, "POINT", p.X, p.Y)
	s.client.Process(cmd)
	return cmd.Err()
}

// Remove player
func (s *PlayerLocationService) Remove(p *Player) error {
	cmd := redis.NewStringCmd("DEL", "player", p.ID)
	s.client.Process(cmd)
	return cmd.Err()
}

type position struct {
	Coords [2]float32 `json:"coordinates"`
}

// All return all registred players
func (s *PlayerLocationService) All() (*PlayerList, error) {
	cmd := redis.NewSliceCmd("SCAN", "player")
	s.client.Process(cmd)
	res, err := cmd.Result()
	if err != nil {
		return nil, err
	}

	payload, _ := redis.NewSliceResult(res[1].([]interface{}), err).Result()
	list := make([]*Player, len(payload))
	for i, item := range payload {
		itemRes, _ := redis.NewSliceResult(item.([]interface{}), nil).Result()
		id, data, geo := itemRes[0].(string), []byte(itemRes[1].(string)), &position{}
		json.Unmarshal(data, geo)

		list[i] = &Player{ID: id, X: geo.Coords[1], Y: geo.Coords[0]}
	}

	return &PlayerList{list}, nil
}

// AddGeofence persist geofence
func (s *PlayerLocationService) AddGeofence(name string, geojson string) error {
	cmd := redis.NewStringCmd("SET", "mapfences", name, "OBJECT", geojson)
	s.client.Process(cmd)
	return cmd.Err()
}

// StreamGeofenceEvents ...
func (s *PlayerLocationService) StreamGeofenceEvents(callback func(msg string)) error {
	conn, err := net.Dial("tcp", ":9851")
	if err != nil {
		return err
	}
	defer conn.Close()

	cmd := "NEARBY player FENCE ROAM mapfences * 0\r\n"
	log.Println("REDIS DEBUG:", cmd)
	if _, err = fmt.Fprintf(conn, cmd); err != nil {
		return err
	}
	buf := make([]byte, 4096)
	n, err := conn.Read(buf)
	res := string(buf[:n])
	if res != "+OK\r\n" {
		return fmt.Errorf("expected OK, got '%v'", res)
	}

	t := time.NewTicker(100 * time.Microsecond)
	for range t.C {
		if n, err = conn.Read(buf); err != nil {
			return err
		}
		for _, line := range strings.Split(string(buf[:n]), "\n") {
			if len(line) == 0 || line[0] != '{' {
				continue
			}
			callback(line)
		}
	}

	return nil
}
