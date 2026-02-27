package laplace

import (
	"context"
	"encoding/json"
	"strings"
	"testing"

	"github.com/stretchr/testify/suite"
)

type WebSocketTestSuite struct {
	*ClientTestSuite
}

func TestWebSocket(t *testing.T) {
	suite.Run(t, &WebSocketTestSuite{
		NewClientTestSuite(),
	})
}

func (s *WebSocketTestSuite) TestGetWebSocketUrl() {
	client := newTestClient(s.Config)
	ctx := context.Background()

	resp, err := client.GetWebSocketUrl(ctx, "test-user", []FeedType{FeedTypeLivePriceTR})
	s.Require().NoError(err)
	s.Require().IsType("", resp)
	s.Require().NotEmpty(resp)
}

func (s *WebSocketTestSuite) TestGetWebsocketUsageForMonth() {
	client := newTestClient(s.Config)
	ctx := context.Background()

	resp, err := client.GetWebsocketUsageForMonth(ctx, 1, 2025, FeedTypeLivePriceTR)
	s.Require().NoError(err)
	s.Require().NotNil(resp)

	if len(resp) > 0 {
		usage := resp[0]
		s.Require().NotEmpty(usage.ExternalUserID)
		s.Require().NotZero(usage.FirstConnectionTime)
		s.Require().Greater(usage.UniqueDeviceCount, int64(0))
	}
}

func (s *WebSocketTestSuite) TestRevokeWebSocketConnection() {
	client := newTestClient(s.Config)
	ctx := context.Background()

	// Generate a WebSocket URL to get a valid connection ID
	url, err := client.GetWebSocketUrl(ctx, "test-revoke-user", []FeedType{FeedTypeLivePriceTR})
	s.Require().NoError(err)
	s.Require().NotEmpty(url)

	// Extract the UUID from the URL (last path segment)
	parts := strings.Split(url, "/")
	id := parts[len(parts)-1]
	s.Require().NotEmpty(id)

	// Revoke the connection — may return 403 if the API key lacks permission
	err = client.RevokeWebSocketConnection(ctx, id)
	if err != nil {
		s.Require().ErrorIs(err, ErrYouDoNotHaveAccessToEndpoint)
	}
}

func (s *WebSocketTestSuite) TestSendWebsocketEvent() {
	client := newTestClient(s.Config)
	ctx := context.Background()

	event, _ := json.Marshal(map[string]string{"type": "test", "message": "hello"})
	broadcastToAll := true

	err := client.SendWebsocketEvent(ctx, SendWebsocketEventRequest{
		Event:          event,
		Transient:      &broadcastToAll,
		BroadCastToAll: true,
	})
	s.Require().NoError(err)
}
