package data

import (
	"encoding/json"
	"os"
	"strconv"
)

// Config json key
type Config struct {
	Port   Value
	DBName Value
	AppID  Value
}

// NewConf return new parser object
func NewConf() *Config {
	return &Config{}
}

// SetPort --
func (c *Config) SetPort(n int) {
	c.Port = Value(strconv.Itoa(n))
}

// Open the json file
func (c *Config) Open(filename string) error {
	in, err := os.Open(filename)
	if err != nil {
		return err
	}

	defer in.Close()

	decodeJSON := json.NewDecoder(in)
	err = decodeJSON.Decode(c)
	if err != nil {
		return err
	}

	return nil
}
