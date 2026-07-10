package laplace

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"
)

type ScreenerTestSuite struct {
	*ClientTestSuite
}

func TestScreener(t *testing.T) {
	suite.Run(t, &ScreenerTestSuite{
		NewClientTestSuite(),
	})
}

func (s *ScreenerTestSuite) TestScreenerBasic() {
	client := newTestClient(s.Config)

	ctx := context.Background()

	resp, err := client.Screener(ctx, RegionTr, ScreenerRequest{
		PageSize: 5,
	})
	s.Require().NoError(err)
	s.Require().Greater(resp.RecordCount, 0)
	s.Require().NotEmpty(resp.Items)
	s.Require().LessOrEqual(len(resp.Items), 5)

	for _, item := range resp.Items {
		s.Require().NotEmpty(item.Symbol)
	}
}

func (s *ScreenerTestSuite) TestScreenerWithFiltersAndSort() {
	client := newTestClient(s.Config)

	ctx := context.Background()

	minPrice := 10.0
	maxPrice := 1000.0

	resp, err := client.Screener(ctx, RegionTr, ScreenerRequest{
		Filters: &ScreenerFilters{
			Price: &ScreenerRange{Min: &minPrice, Max: &maxPrice},
		},
		SortBy:    ScreenerSortByMarketCap,
		SortOrder: SortDirectionDesc,
		Page:      1,
		PageSize:  10,
	})
	s.Require().NoError(err)
	s.Require().NotEmpty(resp.Items)

	var prev *float64
	for _, item := range resp.Items {
		s.Require().NotEmpty(item.Symbol)
		if item.Price != nil {
			s.Require().GreaterOrEqual(*item.Price, minPrice)
			s.Require().LessOrEqual(*item.Price, maxPrice)
		}
		if prev != nil && item.MarketCap != nil {
			s.Require().LessOrEqual(*item.MarketCap, *prev)
		}
		if item.MarketCap != nil {
			mc := *item.MarketCap
			prev = &mc
		}
	}
}

func (s *ScreenerTestSuite) TestScreenerRatingAndTechnicalFilters() {
	client := newTestClient(s.Config)

	ctx := context.Background()

	minComposite := 90.0
	minPriceVsSma200 := 0.0
	minOffHigh := -15.0
	maxOffHigh := 0.0
	epsAcceleration := true

	resp, err := client.Screener(ctx, RegionTr, ScreenerRequest{
		Filters: &ScreenerFilters{
			CompositeRating: &ScreenerRange{Min: &minComposite},
			SMRRating:       []ScreenerRating{ScreenerRatingA, ScreenerRatingB},
			ADRating:        []ScreenerRating{ScreenerRatingA, ScreenerRatingB},
			EPSAcceleration: &epsAcceleration,
			PriceVsSMA200:   &ScreenerRange{Min: &minPriceVsSma200},
			OffHighPct:      &ScreenerRange{Min: &minOffHigh, Max: &maxOffHigh},
		},
		SortBy:    ScreenerSortByCompositeRating,
		SortOrder: SortDirectionDesc,
		Page:      1,
		PageSize:  50,
	})
	s.Require().NoError(err)
	s.Require().LessOrEqual(len(resp.Items), 50)

	var prev *int
	for _, item := range resp.Items {
		s.Require().NotEmpty(item.Symbol)
		if item.CompositeRating != nil {
			s.Require().GreaterOrEqual(float64(*item.CompositeRating), minComposite)
		}
		if item.SMRRating != nil {
			s.Require().Contains([]ScreenerRating{ScreenerRatingA, ScreenerRatingB}, *item.SMRRating)
		}
		if item.ADRating != nil {
			s.Require().Contains([]ScreenerRating{ScreenerRatingA, ScreenerRatingB}, *item.ADRating)
		}
		if item.EPSAcceleration != nil {
			s.Require().True(*item.EPSAcceleration)
		}
		if prev != nil && item.CompositeRating != nil {
			s.Require().LessOrEqual(*item.CompositeRating, *prev)
		}
		if item.CompositeRating != nil {
			cr := *item.CompositeRating
			prev = &cr
		}
	}
}
