package cmd

import (
	"context"
	"fmt"
	"log"

	"github.com/gopxl/beep/v2"
	"github.com/mikowitz/intonation/internal"
	intonation "github.com/mikowitz/intonation/pkg"
	"github.com/spf13/cobra"
)

var (
	ratioRatio    string
	ratioPlay     bool
	ratioCompare  bool
	ratioInterval bool
	ratioQuiet    bool
)

var ratioCmd = &cobra.Command{
	Use:   "ratio <ratio>",
	Short: "Compare a ratio to its closest 12-EDO interval",
	Long:  `Compare a ratio to its closest 12-EDO interval, optionally playing both intervals for audio comparison.`,
	Args:  cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		ratioStr := ratioRatio
		ratio, err := intonation.NewRatioFromString(ratioStr)
		if err != nil {
			log.Fatal(err)
		}
		interval := ratio.Approximate12EDOInterval()
		fmt.Println(ratio, "\t", interval)

		if ratioPlay && !ratioQuiet {
			ctx, ctxCancel := context.WithCancel(context.Background())
			defer ctxCancel()
			output := internal.BeepAudioOutput{SampleRate: beep.SampleRate(48000)}
			if ratioInterval {
				ratio.PlayInterval(ctx, output)
			}
			ratio.PlayChord(ctx, output)
			if ratioCompare {
				interval.Interval().PlayChord(ctx, output)
			}
		}
	},
}

func init() {
	ratioCmd.Flags().StringVarP(&ratioRatio, "ratio", "r", "1/1", "the ratio to calculate")
	ratioCmd.Flags().BoolVarP(&ratioPlay, "play", "p", true, "play the ratio")
	ratioCmd.Flags().BoolVarP(&ratioCompare, "compare", "c", false, "play the nearest 12-EDO interval as a comparison")
	ratioCmd.Flags().BoolVarP(&ratioInterval, "interval", "i", false, "play the ratio as a split dyad before playing it as a chord")
	ratioCmd.Flags().BoolVarP(&ratioQuiet, "quiet", "q", false, "compare ratio with no audio output")
	ratioCmd.MarkFlagRequired("ratio")

	ratioCmd.MarkFlagsMutuallyExclusive("play", "quiet")
}
