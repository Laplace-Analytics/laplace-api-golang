package laplace

import (
	"context"
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
