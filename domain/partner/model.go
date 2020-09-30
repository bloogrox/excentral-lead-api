package partner

// Partner ...
type Partner struct {
	ID       uint   `gorm:"primaryKey"`
	DealType string `gorm:"size:3"`
}
