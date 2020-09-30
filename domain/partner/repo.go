package partner

// Repo ...
type Repo interface {
	Get(int) (*Partner, error)
}
