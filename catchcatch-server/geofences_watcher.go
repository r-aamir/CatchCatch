package main

import (
	"context"
	"log"

	"github.com/golang/protobuf/proto"
	"github.com/perenecabuto/CatchCatch/catchcatch-server/protobuf"
)

// WatchGeofences watch for geofences events and notify players around
func (gw *GameWatcher) WatchGeofences(ctx context.Context) error {
	// TODO: only notify admins about new geofences
	return gw.stream.StreamNearByEvents(ctx, "geofences", "player", "*", DefaultWatcherRange, func(d *Detection) error {
		switch d.Intersects {
		case Inside:
			coords := `{"type":"Polygon","coordinates":` + d.Coordinates + "}"
			err := gw.wss.Emit(d.NearByFeatID, &protobuf.Feature{EventName: proto.String("admin:feature:added"), Id: &d.FeatID,
				Group: proto.String("geofences"), Coords: &coords})
			if err != nil {
				log.Println("admin:feature:added error", err.Error())
			}
		}
		return nil
	})
}
