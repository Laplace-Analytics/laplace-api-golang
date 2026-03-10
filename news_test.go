package laplace

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
)

type NewsTestSuite struct {
	*ClientTestSuite
}

func TestNews(t *testing.T) {
	suite.Run(t, &NewsTestSuite{
		NewClientTestSuite(),
	})
}

func (s *NewsTestSuite) TestGetNewsHighlights() {
	client := newTestClient(s.Config)

	ctx := context.Background()

	resp, err := client.GetNewsHighlights(ctx, RegionUs, LocaleEn)
	s.Require().NoError(err)
	s.Require().NotNil(resp)
	s.Require().NotNil(resp.Consumer)
	s.Require().NotNil(resp.EnergyAndUtilities)
	s.Require().NotNil(resp.Finance)
	s.Require().NotNil(resp.Healthcare)
	s.Require().NotNil(resp.IndustrialsAndMaterials)
	s.Require().NotNil(resp.Tech)
	s.Require().NotNil(resp.Other)
}

func (s *NewsTestSuite) TestGetNews() {
	client := newTestClient(s.Config)

	ctx := context.Background()

	size := 10
	resp, err := client.GetNews(ctx, GetNewsParams{
		Region: RegionUs,
		Locale: LocaleEn,
		Size:   &size,
	})
	s.Require().NoError(err)
	s.Require().NotNil(resp)
	s.Require().Greater(len(resp.Items), 0)

	news := resp.Items[0]
	s.Require().NotEmpty(news.URL)
	s.Require().NotZero(news.Timestamp)
	s.Require().NotEmpty(news.PublisherUrl)
	s.Require().NotZero(news.CreatedAt)
	s.Require().NotEmpty(news.Publisher.Name)
	s.Require().NotNil(news.RelatedTickers)

	if len(news.RelatedTickers) > 0 {
		ticker := news.RelatedTickers[0]
		s.Require().NotEmpty(ticker.ID)
		s.Require().NotEmpty(ticker.Name)
	}

	if news.Categories != nil {
		s.Require().NotEmpty(news.Categories.Name)
		s.Require().GreaterOrEqual(news.Categories.NewsCount, int64(0))
	}

	if news.Sectors != nil {
		s.Require().NotEmpty(news.Sectors.Name)
		s.Require().GreaterOrEqual(news.Sectors.NewsCount, int64(0))
	}

	if news.Industries != nil {
		s.Require().NotEmpty(news.Industries.Name)
	}

	if news.Content != nil {
		s.Require().NotEmpty(news.Content.Title)
		s.Require().NotEmpty(news.Content.Description)
		s.Require().NotNil(news.Content.Content)
		s.Require().NotNil(news.Content.Summary)
	}
}

func (s *NewsTestSuite) TestGetNewsWithType() {
	client := newTestClient(s.Config)

	ctx := context.Background()

	size := 10
	resp, err := client.GetNews(ctx, GetNewsParams{
		Region:   RegionUs,
		Locale:   LocaleEn,
		NewsType: NewsTypeBloomberg,
		Size:     &size,
	})
	s.Require().NoError(err)
	s.Require().NotNil(resp)
}

func (s *NewsTestSuite) TestGetNewsStream() {
	client := newTestClient(s.Config)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	stream, err := client.CreateNewsStream(ctx, StreamNewsParams{Locale: LocaleEn})
	s.Require().NoError(err)
	defer stream.Close()

	receiveChan := stream.Receive()

	select {
	case data := <-receiveChan:
		if data.Error != nil {
			s.T().Logf("Received error: %v", data.Error)
		} else {
			s.T().Logf("Received stream data slice length: %d", len(data.Data))
			if len(data.Data) > 0 {
				news := data.Data[0]
				s.Require().NotEmpty(news.URL)
				s.Require().NotZero(news.Timestamp)
			}
		}
	case <-ctx.Done():
		s.T().Log("Timeout waiting for stream data")
	}
}
