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

// RequestCreateLead ...
type RequestCreateLead struct {
	Email     string
	FirstName string
	LastName  string
	Phone     string
	Language  string
	Country   string
	Marker    string
}

// ResponseCreateLead ...
type ResponseCreateLead struct {
	ReturnCode         int      `json:"returnCode"`
	Description        string   `json:"description"`
	InvalidFields      []string `json:"invalidFields"`
	TimestampGenerated string   `json:"timestampGenerated"`
	Data               struct {
		LeadID         string `json:"leadID"`
		DateRegistered string `json:"dateRegistered"`
		Status         string `json:"status"`
	} `json:"data"`
}

// CreateLead ...
func (a *API) CreateLead(r RequestCreateLead) (*ResponseCreateLead, error) {
	data := url.Values{}
	data.Set("affiliateID", fmt.Sprintf("%d", a.affiliateID))
	data.Set("email", r.Email)
	data.Set("firstName", r.FirstName)
	data.Set("lastName", r.LastName)
	data.Set("phone", r.Phone)
	data.Set("language", r.Language)
	data.Set("country", r.Country)
	data.Set("marker", r.Marker)

	checksum := a.generateChecksum(data)
	data.Set("checksum", checksum)

	req, err := http.NewRequest("POST", baseURL+"lead/create", bytes.NewBuffer([]byte(data.Encode())))
	if err != nil {
		return nil, errors.Wrap(err, "API.CreateLead")
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "API.CreateLead")
	}

	var response ResponseCreateLead
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return nil, errors.Wrap(err, "API.CreateLead")
	}

	switch response.ReturnCode {
	case 1:
		return &response, nil
	case 2, 3:
		fieldsStr := strings.Join(response.InvalidFields, ", ")
		return nil, errors.Wrap(
			errMap[response.ReturnCode],
			"API.CreateLead: fields: "+fieldsStr,
		)
	case 4, 5, 7, 8, 9, 11, 24:
		return nil, errors.Wrap(errMap[response.ReturnCode], "API.CreateLead")
	default:
		return nil, errors.Wrap(errors.New("Unknown Return Code"), "API.CreateLead")
	}
}
