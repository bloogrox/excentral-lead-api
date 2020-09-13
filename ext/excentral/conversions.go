package excentral

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/pkg/errors"
)

// RequestGetConversions ...
type RequestGetConversions struct {
	DateFrom string
	DateTo   string
	Limit    int
	Offset   int
}

// ResponseGetConversions ...
type ResponseGetConversions struct {
	ReturnCode         int      `json:"returnCode"`
	Description        string   `json:"description"`
	InvalidFields      []string `json:"invalidFields"`
	TimestampGenerated string   `json:"timestampGenerated"`
	Data               struct {
		AllLeadsCount            string `json:"allLeadsCount"`
		SelectedLeadsCount       string `json:"selectedLeadsCount"`
		Limit                    string `json:"limit"`
		Offset                   string `json:"offset"`
		AllLeadsEarnedTotal      string `json:"allLeadsEarnedTotal"`
		SelectedLeadsEarnedTotal string `json:"selectedLeadsEarnedTotal"`
		Leads                    []struct {
			LeadID         string `json:"leadID"`
			DateRegistered string `json:"dateRegistered"`
			DateConverted  string `json:"dateConverted"`
			Status         string `json:"status"`
		} `json:"leads"`
	} `json:"data"`
}

// GetConversions ...
func (a *API) GetConversions(r RequestGetConversions) (*ResponseGetConversions, error) {
	data := url.Values{}
	data.Set("affiliateID", fmt.Sprintf("%d", a.affiliateID))
	data.Set("dateFrom", r.DateFrom)
	data.Set("dateTo", r.DateTo)
	data.Set("limit", fmt.Sprintf("%d", r.Limit))
	data.Set("offset", fmt.Sprintf("%d", r.Offset))

	checksum := a.generateChecksum(data)
	data.Set("checksum", checksum)

	req, err := http.NewRequest(
		"POST",
		baseURL+"lead/getConversions",
		bytes.NewBuffer([]byte(data.Encode())),
	)
	if err != nil {
		return nil, errors.Wrap(err, "API.GetConversions")
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "API.GetConversions")
	}

	var response ResponseGetConversions
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return nil, errors.Wrap(err, "API.GetConversions")
	}

	switch response.ReturnCode {
	case 1:
		return &response, nil
	case 2, 3:
		fieldsStr := strings.Join(response.InvalidFields, ", ")
		return nil, errors.Wrap(
			errMap[response.ReturnCode],
			"API.GetConversions: fields: "+fieldsStr,
		)
	case 4, 5, 7, 8, 9, 11, 24:
		return nil, errors.Wrap(errMap[response.ReturnCode], "API.GetConversions")
	default:
		return nil, errors.Wrap(errors.New("Unknown Return Code"), "API.GetConversions")
	}
}
