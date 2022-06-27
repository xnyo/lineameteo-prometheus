package config

import "os"

var (
	HTTPListen        = Value{"HTTP_LISTEN_ADDR", "0.0.0.0:3000"}
	WantedLocationIDs = Value{"WANTED_LOCATION_IDS", "1793,2166"}
)

type Value struct {
	EnvVar  string
	Default string
}

func (c *Value) Get() string {
	v := os.Getenv(c.EnvVar)
	if v == "" {
		return c.Default
	}
	return v
}
