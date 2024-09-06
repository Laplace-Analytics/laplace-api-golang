package client

import (
	"context"
	"testing"
	"time"

	"finfree.co/laplace/utilities"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type LivePriceTestSuite struct {
	*utilities.ClientTestSuite
}

func TestLivePrice(t *testing.T) {
	suite.Run(t, &LivePriceTestSuite{
		utilities.NewClientTestSuite(),
	})
}

func (s *LivePriceTestSuite) TestBISTLivePrice() {
	client := NewClient(s.Config, logrus.New())
	symbols := []string{"TUPRS", "SASA", "THYAO", "GARAN", "YKBN"}

	testLivePrice(s.T(), client, symbols, RegionTr)
}

func (s *LivePriceTestSuite) TestUSLivePrice() {
	client := NewClient(s.Config, logrus.New())
	symbols := []string{"AAPL", "GOOGL", "MSFT", "AMZN", "META"}

	testLivePrice(s.T(), client, symbols, RegionUs)
}

func testLivePrice(t *testing.T, client *Client, symbols []string, region Region) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	livePrices, errs, err := client.GetLivePriceForBIST(ctx, symbols, region)
	require.NoError(t, err)

	livePriceCount := 0

	for {
		select {
		case livePrice, ok := <-livePrices:
			if !ok {
				// Channel closed
				break
			}
			livePriceCount++
			require.NotEmpty(t, livePrice)
			if livePriceCount > 3 {
				return
			}
		case err, ok := <-errs:
			if !ok {
				// Error channel closed
				continue
			}
			require.Fail(t, "Error occurred during live price retrieval", err)
			return
		}
	}
}
