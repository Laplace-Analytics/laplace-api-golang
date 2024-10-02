package laplace

import (
	"context"
	"testing"

	"github.com/samber/lo"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CustomThemeTestSuite struct {
	*ClientTestSuite
}

func TestCustomTheme(t *testing.T) {
	suite.Run(t, &CustomThemeTestSuite{
		NewClientTestSuite(),
	})
}

func (s *CustomThemeTestSuite) TestGetAllCustomThemes() {
	client := NewClient(s.Config, logrus.New())

	ctx := context.Background()

	resp, err := client.GetAllCustomThemes(ctx, LocaleTr)
	require.NoError(s.T(), err)

	require.NotEmpty(s.T(), resp)
}

func (s *CustomThemeTestSuite) TestCreateUpdateDeleteCustomTheme() {
	client := NewClient(s.Config, logrus.New())

	ctx := context.Background()

	stocks, err := client.GetAllStocks(ctx, RegionTr)
	require.NoError(s.T(), err)
	require.NotEmpty(s.T(), stocks)

	createParams := &CreateCustomThemeParams{
		Title: LocaleString{
			LocaleTr: "Test Custom Theme",
		},
		Description: LocaleString{
			LocaleTr: "Test Custom Theme Description",
		},
		Region:         []Region{RegionTr},
		ImageURL:       "Test Custom Theme Image URL",
		Image:          "Test Custom Theme Image",
		AvatarImageURL: "Test Custom Theme Avatar Image",
		Stocks:         []primitive.ObjectID{stocks[0].ID, stocks[1].ID},
		Status:         CollectionStatusActive,
	}
	id := testCreateCustomTheme(s, client, ctx, createParams)
	testGetDetails(s, *id, LocaleTr, client, ctx, createParams)

	updateParams := &UpdateCustomThemeParams{
		Stocks: []primitive.ObjectID{stocks[0].ID, stocks[2].ID},
	}
	testUpdateCustomTheme(s, *id, client, ctx, updateParams)
	applyUpdateParams(updateParams, createParams)
	testGetDetails(s, *id, LocaleTr, client, ctx, createParams)

	updateParams = &UpdateCustomThemeParams{
		Title: LocaleString{
			LocaleTr: "Test Custom Theme Title Updated",
			LocaleEn: "Test Custom Theme Title Updated",
		},
		Description: LocaleString{
			LocaleTr: "Test Custom Theme Description Updated",
			LocaleEn: "Test Custom Theme Description Updated",
		},
	}
	testUpdateCustomTheme(s, *id, client, ctx, updateParams)
	applyUpdateParams(updateParams, createParams)
	testGetDetails(s, *id, LocaleTr, client, ctx, createParams)
	testGetDetails(s, *id, LocaleEn, client, ctx, createParams)

	updateParams = &UpdateCustomThemeParams{
		Status: CollectionStatusInactive,
	}
	testUpdateCustomTheme(s, *id, client, ctx, updateParams)
	applyUpdateParams(updateParams, createParams)
	testGetDetails(s, *id, LocaleTr, client, ctx, createParams)

	testDeleteCustomTheme(s, *id, client, ctx)
	resp, err := client.GetCustomThemeDetail(ctx, id.Hex(), "", LocaleTr, "")
	require.Error(s.T(), err)
	require.Empty(s.T(), resp)
}

func testCreateCustomTheme(s *CustomThemeTestSuite, client *Client, ctx context.Context, createParams *CreateCustomThemeParams) *primitive.ObjectID {
	resp, err := client.CreateCustomTheme(ctx, *createParams)
	require.NoError(s.T(), err)
	require.NotEmpty(s.T(), resp)
	return resp
}

func testGetDetails(s *CustomThemeTestSuite, id primitive.ObjectID, locale Locale, client *Client, ctx context.Context, createParams *CreateCustomThemeParams) {
	resp, err := client.GetCustomThemeDetail(ctx, id.Hex(), "", locale, "")
	require.NoError(s.T(), err)

	require.NotEmpty(s.T(), resp)

	require.Equal(s.T(), createParams.Title[locale], resp.Title)
	require.Equal(s.T(), createParams.Description[locale], resp.Description)
	require.Equal(s.T(), createParams.Region, resp.Region)
	require.Equal(s.T(), createParams.ImageURL, resp.ImageUrl)
	require.Equal(s.T(), createParams.Image, resp.Image)
	require.Equal(s.T(), createParams.AvatarImageURL, resp.AvatarUrl)
	require.Equal(s.T(), createParams.Stocks, lo.Map(resp.Stocks, func(stock Stock, _ int) primitive.ObjectID {
		return stock.ID
	}))
	require.Equal(s.T(), createParams.Status, resp.Status)
}

func testUpdateCustomTheme(s *CustomThemeTestSuite, id primitive.ObjectID, client *Client, ctx context.Context, updateParams *UpdateCustomThemeParams) {
	err := client.UpdateCustomTheme(ctx, id, *updateParams)
	require.NoError(s.T(), err)
}

func applyUpdateParams(updateParams *UpdateCustomThemeParams, createParams *CreateCustomThemeParams) {
	if updateParams.Stocks != nil {
		createParams.Stocks = updateParams.Stocks
	}
	if updateParams.Title != nil {
		createParams.Title = updateParams.Title
	}
	if updateParams.Description != nil {
		createParams.Description = updateParams.Description
	}
	if updateParams.ImageURL != "" {
		createParams.ImageURL = updateParams.ImageURL
	}
	if updateParams.Image != "" {
		createParams.Image = updateParams.Image
	}
	if updateParams.AvatarImageURL != "" {
		createParams.AvatarImageURL = updateParams.AvatarImageURL
	}
	if updateParams.Status != "" {
		createParams.Status = updateParams.Status
	}
	if updateParams.MetaData != nil {
		createParams.MetaData = updateParams.MetaData
	}
}

func testDeleteCustomTheme(s *CustomThemeTestSuite, id primitive.ObjectID, client *Client, ctx context.Context) {
	err := client.DeleteCustomTheme(ctx, id)
	require.NoError(s.T(), err)
}
