package laplace

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

type PaginatedResponse[T any] struct {
	RecordCount int `json:"recordCount"`
	Items       []T `json:"items"`
}

type CapitalIncrease struct {
	ID                            int        `json:"id"`
	BoardDecisionDate             *time.Time `json:"boardDecisionDate"`
	RegisteredCapitalCeiling      string     `json:"registeredCapitalCeiling"`
	CurrentCapital                string     `json:"currentCapital"`
	TargetCapital                 string     `json:"targetCapital"`
	Types                         []string   `json:"types"`
	SpkApplicationResult          *string    `json:"spkApplicationResult"`
	SpkApplicationDate            *time.Time `json:"spkApplicationDate"`
	SpkApprovalDate               *time.Time `json:"spkApprovalDate"`
	PaymentDate                   *time.Time `json:"paymentDate"`
	RegistrationDate              *time.Time `json:"registrationDate"`
	SpecifiedCurrency             string     `json:"specifiedCurrency"`
	Symbol                        string     `json:"symbol"`
	RelatedDisclosureIDs          []int      `json:"relatedDisclosureIds"`
	RightsRate                    string     `json:"rightsRate"`
	RightsPrice                   string     `json:"rightsPrice"`
	RightsTotalAmount             string     `json:"rightsTotalAmount"`
	RightsStartDate               *time.Time `json:"rightsStartDate"`
	RightsEndDate                 *time.Time `json:"rightsEndDate"`
	RightsLastSellDate            *string    `json:"rightsLastSellDate"`
	BonusRate                     string     `json:"bonusRate"`
	BonusTotalAmount              string     `json:"bonusTotalAmount"`
	BonusStartDate                *time.Time `json:"bonusStartDate"`
	BonusDividendRate             string     `json:"bonusDividendRate"`
	BonusDividendTotalAmount      string     `json:"bonusDividendTotalAmount"`
	ExternalCapitalIncreaseAmount string     `json:"externalCapitalIncreaseAmount"`
	ExternalCapitalIncreaseRate   string     `json:"externalCapitalIncreaseRate"`
}

// GetAllCapitalIncreases retrieves all capital increase announcements and events with pagination.
func (c *Client) GetAllCapitalIncreases(ctx context.Context, page int, pageSize int, region Region) (*PaginatedResponse[CapitalIncrease], error) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/api/v1/capital-increase/all", c.baseUrl), nil)
	if err != nil {
		return nil, err
	}

	q := req.URL.Query()
	q.Add("page", strconv.Itoa(page))
	q.Add("size", strconv.Itoa(pageSize))
	q.Add("region", string(region))
	req.URL.RawQuery = q.Encode()

	resp, err := sendRequest[PaginatedResponse[CapitalIncrease]](ctx, c, req)
	if err != nil {
		return nil, err
	}

	return &resp, nil
}

// GetCapitalIncreasesForInstrument fetches capital increase events for a specific stock symbol.
func (c *Client) GetCapitalIncreasesForInstrument(ctx context.Context, symbol string, page int, pageSize int, region Region) (*PaginatedResponse[CapitalIncrease], error) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/api/v1/capital-increase/%s", c.baseUrl, symbol), nil)
	if err != nil {
		return nil, err
	}

	q := req.URL.Query()
	q.Add("page", strconv.Itoa(page))
	q.Add("size", strconv.Itoa(pageSize))
	q.Add("region", string(region))
	req.URL.RawQuery = q.Encode()

	resp, err := sendRequest[PaginatedResponse[CapitalIncrease]](ctx, c, req)
	if err != nil {
		return nil, err
	}

	return &resp, nil
}

// GetActiveRightsForInstrument retrieves active rights offerings for a specific stock on a given date.
func (c *Client) GetActiveRightsForInstrument(ctx context.Context, symbol string, date string, region Region) ([]CapitalIncrease, error) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/api/v1/rights/active/%s", c.baseUrl, symbol), nil)
	if err != nil {
		return nil, err
	}

	q := req.URL.Query()
	q.Add("date", date)
	q.Add("region", string(region))
	req.URL.RawQuery = q.Encode()

	resp, err := sendRequest[[]CapitalIncrease](ctx, c, req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
