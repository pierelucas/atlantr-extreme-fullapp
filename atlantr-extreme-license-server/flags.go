package main

import "flag"

var (
	flagLOGOUTPUT = flag.String(
		"l",
		"",
		"Logoutput",
	)
)

func init() {
	flag.Parse() // Parse our command-line arguments
}
