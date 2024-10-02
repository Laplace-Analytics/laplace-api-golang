package laplace

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type LaplaceClientTestSuite struct {
	*ClientTestSuite
}

func TestLaplaceClient(t *testing.T) {
	suite.Run(t, &LaplaceClientTestSuite{
		NewClientTestSuite(),
	})
}

func (s *LaplaceClientTestSuite) TestClient() {
	client := NewClient(s.Config, logrus.New())

	req, err := http.NewRequest("GET", s.Config.BaseURL+"/api/v1/industry", nil)
	s.Require().NoError(err)
	q := req.URL.Query()
	q.Add("region", string(RegionTr))
	q.Add("locale", string(LocaleTr))
	req.URL.RawQuery = q.Encode()

	res, err := sendRequest[any](context.Background(), client, req)
	s.Require().NoError(err)
	s.Require().NotEmpty(res)
}

func (s *LaplaceClientTestSuite) TestYouDontHaveAccessError() {
	client := NewClient(s.Config, logrus.New())

	_, err := client.GetAllCollections(context.Background(), "aaa", LocaleTr)
	require.Error(s.T(), err)
	require.True(s.T(), errors.Is(err, ErrYouDoNotHaveAccessToEndpoint))
}

func (s *LaplaceClientTestSuite) TestInvalidToken() {
	invalidConfig := &LaplaceConfiguration{
		BaseURL: s.Config.BaseURL,
		APIKey:  "invalid",
	}

	client := NewClient(*invalidConfig, logrus.New())

	_, err := client.GetAllCollections(context.Background(), RegionTr, LocaleTr)
	require.Error(s.T(), err)
	require.True(s.T(), errors.Is(err, ErrInvalidToken))
}

func (s *LaplaceClientTestSuite) TestInvalidID() {
	client := NewClient(s.Config, logrus.New())

	_, err := client.GetCollectionDetail(context.Background(), "invalid", RegionTr, LocaleTr)
	require.Error(s.T(), err)
	require.True(s.T(), errors.Is(err, ErrInvalidID))

	_, err = client.GetStockDetailByID(context.Background(), "invalid", LocaleTr)
	require.Error(s.T(), err)
	require.True(s.T(), errors.Is(err, ErrInvalidID))
}
