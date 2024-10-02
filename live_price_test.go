package laplace

// import (
// 	"context"
// 	"testing"
// 	"time"

// 	"github.com/sirupsen/logrus"
// 	"github.com/stretchr/testify/require"
// 	"github.com/stretchr/testify/suite"
// )

// type LivePriceTestSuite struct {
// 	*ClientTestSuite
// }

// func TestLivePrice(t *testing.T) {
// 	suite.Run(t, &LivePriceTestSuite{
// 		NewClientTestSuite(),
// 	})
// }

// func (s *LivePriceTestSuite) TestBISTLivePrice() {
// 	client := NewClient(s.Config, logrus.New())
// 	symbols := []string{"TUPRS", "SASA", "THYAO", "GARAN", "YKBN"}

// 	testLivePrice(s.T(), client, symbols, RegionTr)
// }

// func (s *LivePriceTestSuite) TestUSLivePrice() {
// 	client := NewClient(s.Config, logrus.New())
// 	symbols := []string{"AAPL", "GOOGL", "MSFT", "AMZN", "META"}

// 	testLivePrice(s.T(), client, symbols, RegionUs)
// }

// func testLivePrice(t *testing.T, client *Client, symbols []string, region Region) {
// 	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
// 	defer cancel()

// 	var livePrices <-chan StockLiveData
// 	var errs <-chan error
// 	var close func()
// 	var err error
// 	if region == RegionTr {
// 		var bistLiveData <-chan BISTStockLiveData
// 		bistLiveData, errs, close, err = client.GetLivePriceForBIST(ctx, symbols)
// 		livePrices = convertChannel(bistLiveData)
// 	} else {
// 		var usLiveData <-chan USStockLiveData
// 		usLiveData, errs, close, err = getLivePrice[USStockLiveData](client, ctx, symbols, region)
// 		livePrices = convertChannel(usLiveData)
// 	}
// 	require.NoError(t, err)
// 	defer close()
// 	livePriceCount := 0

// 	for {
// 		select {
// 		case livePrice, ok := <-livePrices:
// 			if !ok {
// 				// Channel closed
// 				break
// 			}
// 			livePriceCount++
// 			require.NotEmpty(t, livePrice)
// 			if livePriceCount > 3 {
// 				return
// 			}
// 		case err, ok := <-errs:
// 			if !ok {
// 				// Error channel closed
// 				continue
// 			}
// 			require.Fail(t, "Error occurred during live price retrieval", err)
// 			return
// 		}
// 	}
// }

// // Add this helper function at the end of the file
// func convertChannel[T StockLiveData](ch <-chan T) <-chan StockLiveData {
// 	out := make(chan StockLiveData)
// 	go func() {
// 		defer close(out)
// 		for v := range ch {
// 			out <- v
// 		}
// 	}()
// 	return out
// }
