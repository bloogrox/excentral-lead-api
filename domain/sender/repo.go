package sender

// Repo ...
type Repo interface {
	Get(uint) (*Sender, error)
}
