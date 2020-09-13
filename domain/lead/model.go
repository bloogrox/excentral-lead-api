package lead

import "gorm.io/gorm"

// Lead ...
type Lead struct {
	gorm.Model
	RemoteID  uint   `gorm:"uniqueIndex"`
	Email     string `gorm:"size:50"`
	Phone     string `gorm:"size:50"`
	FirstName string `gorm:"size:50"`
	LastName  string `gorm:"size:50"`
	Language  string `gorm:"size:2"`
	Country   string `gorm:"size:2"`
	Source    string `gorm:"size:300"`
	Sub1      string `gorm:"size:100"`
	PID       int
}
