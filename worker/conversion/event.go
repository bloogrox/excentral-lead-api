package persist

// Event ...
type Event struct {
	LeadID        uint   `json:"lead_id"`
	DateConverted string `json:"date_converted"`
	Status        string `json:"status"`
}
