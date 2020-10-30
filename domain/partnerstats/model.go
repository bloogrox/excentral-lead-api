package partnerstats

// DailyReport ...
type DailyReport struct {
	Date  string `gorm:"column:day" json:"date"`
	Leads int    `json:"leads"`
}
