package sender

// Sender ...
type Sender struct {
	ID  uint `gorm:"primaryKey"`
	PID int
}
