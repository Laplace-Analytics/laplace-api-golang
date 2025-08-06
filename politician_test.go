package laplace

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
)

type PoliticianTestSuite struct {
	*ClientTestSuite
}

func TestPolitician(t *testing.T) {
	suite.Run(t, &PoliticianTestSuite{
		NewClientTestSuite(),
	})
}

func (s *PoliticianTestSuite) TestGetAllPoliticians() {
	client := newTestClient(s.Config)
	ctx := context.Background()

	resp, err := client.GetAllPoliticians(ctx)
	s.Require().NoError(err)
	s.Require().NotNil(resp)
	s.Require().Greater(len(resp), 0)

	politician := resp[0]
	s.Require().Greater(politician.Id, int32(0))
	s.Require().NotEmpty(politician.PoliticianName)
	s.Require().Greater(politician.TotalHoldings, int32(0))
	s.Require().NotEqual(time.Time{}, politician.LastUpdated)
}

func (s *PoliticianTestSuite) TestGetPoliticianHoldingsBySymbol() {
	client := newTestClient(s.Config)
	ctx := context.Background()

	resp, err := client.GetPoliticianHoldingsBySymbol(ctx, "TUPRS")
	s.Require().NoError(err)
	s.Require().NotNil(resp)
	s.Require().Greater(len(resp), 0)

	holding := resp[0]
	s.Require().NotEmpty(holding.PoliticianName)
	s.Require().Equal("TUPRS", holding.Symbol)
	s.Require().NotEmpty(holding.Company)
	s.Require().NotEmpty(holding.Holding)
	s.Require().NotEmpty(holding.Allocation)
	s.Require().NotEqual(time.Time{}, holding.LastUpdated)
}

func (s *PoliticianTestSuite) TestGetAllTopHoldings() {
	client := newTestClient(s.Config)
	ctx := context.Background()

	resp, err := client.GetAllTopHoldings(ctx)
	s.Require().NoError(err)
	s.Require().NotNil(resp)
	s.Require().Greater(len(resp), 0)

	topHolding := resp[0]
	s.Require().NotEmpty(topHolding.Symbol)
	s.Require().NotEmpty(topHolding.Company)
	s.Require().Greater(len(topHolding.Politicians), 0)
	s.Require().Greater(topHolding.Count, int32(0))

	politician := topHolding.Politicians[0]
	s.Require().NotEmpty(politician.Name)
	s.Require().NotEmpty(politician.Holding)
	s.Require().NotEmpty(politician.Allocation)
}

func (s *PoliticianTestSuite) TestGetPoliticianDetail() {
	client := newTestClient(s.Config)
	ctx := context.Background()

	resp, err := client.GetPoliticianDetail(ctx, 1)
	s.Require().NoError(err)
	s.Require().NotNil(resp)

	s.Require().Equal(int32(1), resp.Id)
	s.Require().NotEmpty(resp.Name)
	s.Require().Greater(len(resp.Holdings), 0)
	s.Require().Greater(resp.TotalHoldings, int32(0))
	s.Require().NotEqual(time.Time{}, resp.LastUpdated)

	holding := resp.Holdings[0]
	s.Require().NotEmpty(holding.Symbol)
	s.Require().NotEmpty(holding.Company)
	s.Require().NotEmpty(holding.Holding)
	s.Require().NotEmpty(holding.Allocation)
}
