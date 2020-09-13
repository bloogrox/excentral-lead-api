package postback

// Postback ...
type Postback struct {
	ID  uint `gorm:"primaryKey"`
	PID int
	URL string `gorm:"size:300"`
}
