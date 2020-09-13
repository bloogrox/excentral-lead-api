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

// RequestGetDeposits ...
type RequestGetDeposits struct {
	dateFrom string
	dateTo   string
	limit    int
	offset   int
}

// ResponseGetDeposits ...
type ResponseGetDeposits struct {
	ReturnCode         int      `json:"returnCode"`
	Description        string   `json:"description"`
	InvalidFields      []string `json:"invalidFields"`
	TimestampGenerated string   `json:"timestampGenerated"`
	Data               struct {
		AllDepositsCount      string `json:"allDepositsCount"`
		SelectedDepositsCount string `json:"selectedDepositsCount"`
		Limit                 string `json:"limit"`
		Offset                string `json:"offset"`
		Deposits              []struct {
			LeadID             string `json:"leadID"`
			DepositID          string `json:"depositID"`
			DateDeposited      string `json:"dateDeposited"`
			Amount             string `json:"amount"`
			AmountUSD          string `json:"amountUSD"`
			IsFirstTimeDeposit string `json:"isFirstTimeDeposit"`
			IsValid            string `json:"cuisValidrrency"`
		} `json:"deposits"`
	} `json:"data"`
}

// GetDeposits ...
func (a *API) GetDeposits(r RequestGetDeposits) (*ResponseGetDeposits, error) {
	data := url.Values{}
	data.Set("affiliateID", fmt.Sprintf("%d", a.affiliateID))
	data.Set("dateFrom", r.dateFrom)
	data.Set("dateTo", r.dateTo)
	data.Set("limit", fmt.Sprintf("%d", r.limit))
	data.Set("offset", fmt.Sprintf("%d", r.offset))

	checksum := a.generateChecksum(data)
	data.Set("checksum", checksum)

	req, err := http.NewRequest(
		"POST",
		baseURL+"lead/getDeposits",
		bytes.NewBuffer([]byte(data.Encode())),
	)
	if err != nil {
		return nil, errors.Wrap(err, "API.GetDeposits")
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "API.GetDeposits")
	}

	var response ResponseGetDeposits
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return nil, errors.Wrap(err, "API.GetDeposits")
	}

	switch response.ReturnCode {
	case 1:
		return &response, nil
	case 2, 3:
		fieldsStr := strings.Join(response.InvalidFields, ", ")
		return nil, errors.Wrap(
			errMap[response.ReturnCode],
			"API.GetDeposits: fields: "+fieldsStr,
		)
	case 4, 5, 7, 8, 9, 11, 24:
		return nil, errors.Wrap(errMap[response.ReturnCode], "API.GetDeposits")
	default:
		return nil, errors.Wrap(errors.New("Unknown Return Code"), "API.GetDeposits")
	}
}
