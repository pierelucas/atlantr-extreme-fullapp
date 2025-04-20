package license

import (
	"encoding/json"

	"github.com/pierelucas/atlantr-extreme-license-server/data"
)

// Pair --
type Pair struct {
	ID      data.Value
	LICENSE data.Value
	APPID   data.Value
}

// NewPair generates a new license pair
func NewPair() (*Pair, error) {
	return &Pair{}, nil
}

// Unmarshal the json string
func (p *Pair) Unmarshal(jsonstring string) error {
	if err := json.Unmarshal([]byte(jsonstring), p); err != nil {
		return err
	}
	return nil
}
