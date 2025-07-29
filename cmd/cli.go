package cmd

import (
	"fmt"
	"os"
)

func Run() {
	if len(os.Args) < 2 {
		usage()
		os.Exit(1)
	}

	switch os.Args[1] {
	case "--help", "-h":
		usage()
	case "ratio":
		RatioCommand()
	case "diamond":
		DiamondCommand()
	case "lattice":
		LatticeCommand()
	default:
		usage()
	}
}

func usage() {
	fmt.Print(`Usage:

./go-intonation

  -h / --help
    Print this message
  ratio <ratio> [--[no]-play] [--compare]
    Compare a ratio to its closest 12-EDO interval,
    optionally playing both intervals for audio comparison
  diamond <limits> [--square]
    Create a otonality/utonality diamond from the provided limits
`)
}
