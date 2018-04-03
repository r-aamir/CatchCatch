package core

import (
	"context"
	"log"

	"github.com/golang/protobuf/proto"
	"github.com/tidwall/gjson"

	"github.com/perenecabuto/CatchCatch/server/protobuf"
	"github.com/perenecabuto/CatchCatch/server/service"
	"github.com/perenecabuto/CatchCatch/server/service/messages"
	"github.com/perenecabuto/CatchCatch/server/websocket"
)

// AdminHandler handle websocket events
type AdminHandler struct {
	server   *websocket.WSServer
	players  service.PlayerLocationService
	messages messages.Dispatcher
}

// NewAdminHandler AdminHandler builder
func NewAdminHandler(s *websocket.WSServer, p service.PlayerLocationService, m messages.Dispatcher) *AdminHandler {
	handler := &AdminHandler{s, p, m}
	return handler
}

// OnConnection handles game and admin connection events
func (h *AdminHandler) OnConnection(c *websocket.WSConnListener) {
	log.Println("[AdminHandler] [admin] connected", c.ID)

	lat, lon := 0.0, 0.0
	h.players.SetAdmin(c.ID, lat, lon)

	c.OnDisconnected(func() {
		h.players.RemoveAdmin(c.ID)
		log.Println("[AdminHandler] [admin] disconnected", c.ID)
	})

	c.On("admin:players:disconnect", h.onDisconnectPlayer(c))
	c.On("admin:players:request", h.onRequestPlayers(c))
	c.On("admin:feature:request", h.onRequestFeatures(c))
	c.On("admin:feature:add", h.onAddFeature())
	c.On("admin:clear", h.onClear())
}

func (h *AdminHandler) onRequestPlayers(so *websocket.WSConnListener) func([]byte) {
	return func([]byte) {
		players, err := h.players.All()
		if err != nil {
			log.Println("[AdminHandler] player:request-remotes event error: " + err.Error())
		}
		event := "remote-player:new"
		for _, p := range players {
			if p == nil {
				continue
			}
			err := so.Emit(&protobuf.Player{EventName: &event, Id: &p.ID, Lon: &p.Lon, Lat: &p.Lat})
			if err != nil {
				log.Println("[AdminHandler] player:request-remotes event error: " + err.Error())
			}
		}
	}
}

func (h *AdminHandler) onDisconnectPlayer(c *websocket.WSConnListener) func([]byte) {
	return func(buf []byte) {
		msg := &protobuf.Simple{}
		proto.Unmarshal(buf, msg)
		playerID := msg.GetId()
		log.Println("[AdminHandler] admin:players:disconnect", playerID)
		err := h.server.Close(playerID)
		if err != nil {
			log.Println("[AdminHandler] admin:players:disconnect error", playerID)
		}
		err = h.players.Remove(playerID)
		if err == service.ErrFeatureNotFound {
			// Notify remote-player removal to ghost players on admin
			log.Println("[AdminHandler] admin:players:disconnect:force", playerID)
			c.Emit(&protobuf.Player{EventName: proto.String("admin:players:disconnected"), Id: &playerID})
		}
	}
}

func (h *AdminHandler) onClear() func([]byte) {
	return func([]byte) {
		// TODO: send this message by broaker
		// h.games.Clear()
		h.players.Clear()
		h.server.CloseAll()
	}
}

// Map events

func (h *AdminHandler) onAddFeature() func([]byte) {
	return func(buf []byte) {
		msg := &protobuf.Feature{}
		proto.Unmarshal(buf, msg)

		var err error
		switch msg.GetGroup() {
		case "geofences":
			err = h.players.SetGeofence(msg.GetId(), msg.GetCoords())
		case "checkpoint":
			err = h.players.SetCheckpoint(msg.GetId(), msg.GetCoords())
		}
		if err != nil {
			log.Println("[AdminHandler] Error to create feature:", err)
			return
		}
	}
}

func (h *AdminHandler) onRequestFeatures(c *websocket.WSConnListener) func([]byte) {
	return func(buf []byte) {
		msg := &protobuf.Feature{}
		proto.Unmarshal(buf, msg)

		features, err := h.players.Features()
		if err != nil {
			log.Println("[AdminHandler] Error on sendFeatures:", err)
		}
		event := "admin:feature:added"
		for _, f := range features {
			c.Emit(&protobuf.Feature{EventName: &event, Id: &f.ID, Group: &f.Group, Coords: &f.Coordinates})
		}
	}
}

func (h *AdminHandler) WatchFeatureEvents(ctx context.Context) error {
	stream := make(chan []byte)
	err := h.messages.Subscribe(ctx, FeaturesMessageTopic, func(data []byte) error {
		stream <- data
		return nil
	})
	if err != nil {
		return err
	}

	for {
		select {
		case <-ctx.Done():
			return nil
		case data := <-stream:
			adminID, featID, group, coords, action :=
				gjson.GetBytes(data, "id").String(),
				gjson.GetBytes(data, "featID").String(),
				gjson.GetBytes(data, "group").String(),
				gjson.GetBytes(data, "coordinates").String(),
				gjson.GetBytes(data, "action").String()
			err := h.server.Emit(adminID, &protobuf.Feature{
				EventName: proto.String("admin:feature:" + action), Id: &featID,
				Group: &group, Coords: &coords})
			if err != nil {
				log.Fatal("[AdminHandler] Error to create feature:", err)
			}
		}
	}

}
