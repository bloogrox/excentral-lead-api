package conversion

import "time"

// Conversion ...
type Conversion struct {
	LeadID        uint `gorm:"primaryKey"`
	DateConverted time.Time
	Status        string `gorm:"size:30"`
}
