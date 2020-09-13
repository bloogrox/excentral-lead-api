package conversion

// Repo ...
type Repo interface {
	Get(uint) (*Conversion, error)
	Insert(*Conversion) (uint, error)
}
