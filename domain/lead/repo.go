package lead

// ID ...
// type ID uint

// Repo ...
type Repo interface {
	Get(ID uint) (*Lead, error)
	Insert(*Lead) (uint, error)
}
