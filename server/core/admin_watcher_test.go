package core

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/gogo/protobuf/proto"
	"github.com/perenecabuto/CatchCatch/server/mocks"
	"github.com/perenecabuto/CatchCatch/server/model"
	"github.com/perenecabuto/CatchCatch/server/protobuf"
	"github.com/perenecabuto/CatchCatch/server/service"
	"github.com/perenecabuto/CatchCatch/server/websocket"
)

func TestNewAdminWatcher(t *testing.T) {
	if w := createAdminWatcher(); w == nil {
		t.Fatal("Can't create AdminWatcher")
	}
}

func TestWatchCheckPointsMustNotifyPlayersNearToCheckPoinstsTheDistToIt(t *testing.T) {
	w := createAdminWatcher()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err := w.WatchCheckpoints(ctx)
	assert.NoError(t, err)

	c := &mocks.WSConnection{}
	c.On("Send", mock.Anything).Return(nil)

	wss := w.wss
	cListener := wss.Add(c)
	playerID := cListener.ID

	distToCheckPoint := 12.0
	checkPoint := model.Feature{Group: "checkpoint", ID: "checkpoint1"}

	geoS := w.service.(*MockGeoServiceWithCallback)
	geoS.PlayerNearToFeatureCallback(playerID, distToCheckPoint, checkPoint)

	expected, _ := proto.Marshal(&protobuf.Detection{
		EventName:    proto.String("checkpoint:detected"),
		Id:           &checkPoint.ID,
		FeatId:       &checkPoint.ID,
		NearByFeatId: &playerID,
		NearByMeters: &distToCheckPoint,
	})

	c.AssertCalled(t, "Send", expected)
}

func createAdminWatcher() *AdminWatcher {
	wss := websocket.NewWSServer(&mocks.WSDriver{})
	geoService := &MockGeoServiceWithCallback{}
	return NewAdminWatcher(geoService, wss)
}

type MockGeoServiceWithCallback struct {
	mocks.GeoFeatureService
	PlayersAroundCallback       service.PlayersAroundCallback
	PlayerNearToFeatureCallback service.PlayerNearToFeatureCallback
}

func (gs *MockGeoServiceWithCallback) ObservePlayersAround(_ context.Context, cb service.PlayersAroundCallback) error {
	gs.PlayersAroundCallback = cb
	return nil
}
func (gs *MockGeoServiceWithCallback) ObservePlayerNearToFeature(_ context.Context, _ string, cb service.PlayerNearToFeatureCallback) error {
	gs.PlayerNearToFeatureCallback = cb
	return nil
}
