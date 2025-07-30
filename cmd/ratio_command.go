package cmd

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/gopxl/beep/v2"
	"github.com/mikowitz/intonation/internal"
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
	ratio, err := intonation.NewRatioFromString(ratioStr)
	if err != nil {
		log.Fatal(err)
	}
	interval := ratio.Approximate12EDOInterval()
	fmt.Println(ratio, "\t", interval)

	ratioCmd.Parse(os.Args[3:])

	shouldPlay := !noPlay && play
	if shouldPlay {
		output := internal.BeepAudioOutput{SampleRate: beep.SampleRate(48000)}
		ratio.Play(output)
		if compare {
			interval.Interval().PlayChord(output)
		}
	}
}
