package cmd

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	intonation "github.com/mikowitz/intonation/pkg"
)

var format string

func DiamondCommand() {
	diamondCmd := flag.NewFlagSet("diamond", flag.ExitOnError)
	diamondCmd.StringVar(&format, "format", "diamond", "printing format for the diamond: diamond or square")

	if len(os.Args) < 3 {
		fmt.Println("  limits\n        the limits used to construct the diamond")
		diamondCmd.PrintDefaults()
		os.Exit(1)
	}

	diamondCmd.Parse(os.Args[3:])

	limits := []uint{}

	for _, l := range strings.Split(os.Args[2], ",") {
		d, err := strconv.Atoi(l)
		if err != nil {
			log.Fatalf("Could not convert `%s` to an int", l)
		}
		limits = append(limits, uint(d))
	}

	d := intonation.NewDiamond(limits...)

	fmt.Println(d.String(intonation.DiamondStringFormat(format)))
}
