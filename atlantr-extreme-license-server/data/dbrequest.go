package parser

import "encoding/json"

// RequestJSON object from client
type Pair struct {
	ID  value
	KEY value
}

// NewRequestJSON returns a new RequestJSON pointer
func NewPair() *Pair {
	return &Pair{}
}

func (rj *Pair) Parse(jsonstring string) error {
	if err := json.Unmarshal([]byte(jsonstring), rj); err != nil {
		return err
	}
	return nil
}
