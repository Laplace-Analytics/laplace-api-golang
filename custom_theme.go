package laplace

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (c *Client) GetAllCustomThemes(ctx context.Context, locale Locale) ([]Collection, error) {
	return c.getAllCollections(ctx, CollectionTypeCustomTheme, "", locale)
}

func (c *Client) GetCustomThemeDetail(ctx context.Context, id string, region Region, locale Locale, sortBy SortBy) (CollectionDetail, error) {
	return c.getCollectionDetail(ctx, id, CollectionTypeCustomTheme, region, locale, sortBy)
}

type CollectionStatus string

const (
	CollectionStatusActive   CollectionStatus = "active"
	CollectionStatusInactive CollectionStatus = "inactive"
)

type CreateCustomThemeParams struct {
	Title          LocaleString         `json:"title,omitempty" bson:"title,omitempty" validate:"required"`
	Description    LocaleString         `json:"description,omitempty" bson:"description,omitempty"`
	Region         []Region             `json:"region,omitempty" bson:"region,omitempty"`
	ImageURL       string               `json:"image_url" bson:"image_url"`
	Image          string               `json:"image" bson:"image"`
	AvatarImageURL string               `json:"avatar_url" bson:"avatar_image_url"`
	Stocks         []primitive.ObjectID `json:"stocks" bson:"stocks" validate:"required"`
	Order          int                  `json:"order" bson:"order"`
	Status         CollectionStatus     `json:"status" bson:"status" validate:"required,oneof=active inactive"`
	MetaData       map[string]any       `json:"meta_data,omitempty" bson:"meta_data,omitempty"`
}

func (c *Client) CreateCustomTheme(ctx context.Context, params CreateCustomThemeParams) (*primitive.ObjectID, error) {
	bodyJSON, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}
	body := bytes.NewBuffer(bodyJSON)

	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/api/v1/custom-theme", c.baseUrl), body)
	if err != nil {
		return nil, err
	}

	resp, err := sendRequest[string](ctx, c, req)
	if err != nil {
		return nil, err
	}

	id, err := primitive.ObjectIDFromHex(resp)
	if err != nil {
		return nil, err
	}

	return &id, nil
}

type UpdateCustomThemeParams struct {
	Title          LocaleString         `json:"title,omitempty" bson:"title,omitempty"`
	Description    LocaleString         `json:"description,omitempty" bson:"description,omitempty"`
	ImageURL       string               `json:"image_url" bson:"image_url"`
	Image          string               `json:"image" bson:"image"`
	AvatarImageURL string               `json:"avatar_url" bson:"avatar_image_url"`
	Stocks         []primitive.ObjectID `json:"stockIds" bson:"stockIds"`
	Status         CollectionStatus     `json:"status" bson:"status"`
	MetaData       map[string]any       `json:"meta_data,omitempty" bson:"meta_data,omitempty"`
}

func (c *Client) UpdateCustomTheme(ctx context.Context, id primitive.ObjectID, params UpdateCustomThemeParams) error {
	bodyJSON, err := json.Marshal(params)
	if err != nil {
		return err
	}
	body := bytes.NewBuffer(bodyJSON)

	req, err := http.NewRequest(http.MethodPatch, fmt.Sprintf("%s/api/v1/custom-theme/%s", c.baseUrl, id.Hex()), body)
	if err != nil {
		return err
	}

	_, err = sendRequest[any](ctx, c, req)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) DeleteCustomTheme(ctx context.Context, id primitive.ObjectID) error {
	req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("%s/api/v1/custom-theme/%s", c.baseUrl, id.Hex()), nil)
	if err != nil {
		return err
	}

	_, err = sendRequest[string](ctx, c, req)
	if err != nil {
		return err
	}

	return nil
}
