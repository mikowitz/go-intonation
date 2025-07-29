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

var (
	play    bool
	noPlay  bool
	compare bool
)

func RatioCommand() {
	ratioCmd := flag.NewFlagSet("ratio", flag.ExitOnError)
	ratioCmd.BoolVar(&play, "play", true, "play the ratio")
	ratioCmd.BoolVar(&noPlay, "no-play", false, "don't play the ratio (shortcut for `--play=false`")
	ratioCmd.BoolVar(&compare, "compare", false, "play the nearest 12-EDO interval as a comparison")

	if len(os.Args) < 3 {
		fmt.Println("  ratio\n        the ratio to calculate")
		ratioCmd.PrintDefaults()
		os.Exit(1)
	}

	ratioStr := os.Args[2]
	ratioParts := strings.Split(ratioStr, "/")

	n, err := strconv.Atoi(ratioParts[0])
	if err != nil {
		log.Fatalf("error parsing `%s` as int", ratioParts[0])
	}
	d, err := strconv.Atoi(ratioParts[1])
	if err != nil {
		log.Fatalf("error parsing `%s` as int", ratioParts[1])
	}

	ratio := intonation.NewRatio(uint(n), uint(d))
	interval := ratio.Approximate12EDOInterval()
	fmt.Println(ratio, "\t", ratio.Approximate12EDOInterval())

	ratioCmd.Parse(os.Args[3:])

	shouldPlay := !noPlay && play
	if shouldPlay {
		ratio.Play()
		if compare {
			interval.Interval().PlayChord()
		}
	}
}
