package laplace

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"
)

type SectorTestSuite struct {
	*ClientTestSuite
}

func TestSector(t *testing.T) {
	suite.Run(t, &SectorTestSuite{
		NewClientTestSuite(),
	})
}

func (s *SectorTestSuite) TestGetAllSectors() {
	client := newTestClient(s.Config)

	ctx := context.Background()

	resp, err := client.GetAllSectors(ctx, RegionTr, LocaleTr)
	s.Require().NoError(err)
	s.Require().NotEmpty(resp)

	hasStocks := false
	for _, sector := range resp {
		s.Require().NotEmpty(sector.ID)
		s.Require().NotEmpty(sector.Title)
		if sector.NumStocks > 0 {
			hasStocks = true
		}
	}
	s.Require().True(hasStocks)
}
