package excentral

import (
	"crypto/md5"
	"encoding/hex"
	"net/url"
	"strings"
)

const baseURL = "https://api-partners.excentral-int.com/v1/"

// API ...
type API struct {
	affiliateID int
	partnerID   string
}

// New ...
func New(affiliateID int, partnerID string) *API {
	return &API{affiliateID, partnerID}
}

func (a *API) generateChecksum(data url.Values) string {
	hash := md5.Sum(
		[]byte(
			a.partnerID + data.Encode(),
		),
	)
	return strings.ToUpper(
		hex.EncodeToString(hash[:]),
	)

}
