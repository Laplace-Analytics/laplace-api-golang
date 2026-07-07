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

	resp, err := client.GetNewsHighlights(ctx, GetNewsHighlightsParams{
		Region: RegionUs,
		Locale: LocaleEn,
	})
	s.Require().NoError(err)
	s.Require().NotNil(resp)
	s.Require().Greater(len(resp.Items), 0)

	highlight := resp.Items[0]
	s.Require().NotEmpty(highlight.ID)
	s.Require().NotZero(highlight.CreatedAt)
	s.Require().NotNil(highlight.Consumer)
	s.Require().NotNil(highlight.EnergyAndUtilities)
	s.Require().NotNil(highlight.Finance)
	s.Require().NotNil(highlight.Healthcare)
	s.Require().NotNil(highlight.IndustrialsAndMaterials)
	s.Require().NotNil(highlight.Tech)
	s.Require().NotNil(highlight.Other)
}

func (s *NewsTestSuite) TestGetNewsCategories() {
	client := newTestClient(s.Config)

	ctx := context.Background()

	resp, err := client.GetNewsCategories(ctx, LocaleEn)
	s.Require().NoError(err)
	s.Require().NotNil(resp)
	s.Require().Greater(len(resp), 0)

	for _, category := range resp {
		s.Require().NotEmpty(category.ID)
		s.Require().NotEmpty(category.Name)
	}
}

func (s *NewsTestSuite) TestGetNewsLanes() {
	client := newTestClient(s.Config)

	ctx := context.Background()

	resp, err := client.GetNewsLanes(ctx, GetNewsLanesParams{})
	s.Require().NoError(err)
	s.Require().NotNil(resp)
	s.Require().Greater(len(resp), 0)

	for _, lane := range resp {
		s.Require().NotEmpty(lane.ID)
		s.Require().NotEmpty(lane.Label)
	}
}

func (s *NewsTestSuite) TestGetNewsApiSourceNames() {
	client := newTestClient(s.Config)

	ctx := context.Background()

	resp, err := client.GetNewsApiSourceNames(ctx, GetNewsApiSourceNamesParams{})
	s.Require().NoError(err)
	s.Require().NotNil(resp)
	s.Require().Greater(len(resp), 0)

	for _, source := range resp {
		s.Require().NotEmpty(source.ID)
		s.Require().NotEmpty(source.Name)
	}
}

func (s *NewsTestSuite) TestGetNewsWithLane() {
	client := newTestClient(s.Config)

	ctx := context.Background()

	size := 10
	resp, err := client.GetNews(ctx, GetNewsParams{
		Region: RegionUs,
		Locale: LocaleEn,
		Lane:   NewsLaneGlobalMacro,
		Size:   &size,
	})
	s.Require().NoError(err)
	s.Require().NotNil(resp)
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

	stream, err := client.CreateNewsStream(ctx, StreamNewsParams{
		Region:      RegionUs,
		Locale:      LocaleEn,
		Symbols:     []string{"AAPL", "GOOGL"},
		CategoryIds: []string{"1", "2"},
		SectorIds:   []string{"65533e047844ee7afe9941bf"},
		IndustryIds: []string{"65533e441fa5c7b58afa0944"},
	})
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
