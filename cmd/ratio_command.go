package cmd

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/gopxl/beep/v2"
	"github.com/mikowitz/intonation/cmd/ratio"
	"github.com/mikowitz/intonation/internal"
	"github.com/spf13/cobra"
)

var (
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
	RunE: func(cmd *cobra.Command, args []string) error {
		f, err := tea.LogToFile("debug.log", "ratio")
		if err != nil {
			fmt.Println("fatal:", err)
			os.Exit(1)
		}
		defer f.Close()
		output := internal.BeepAudioOutput{SampleRate: beep.SampleRate(48000)}
		p := tea.NewProgram(ratio.NewRatioUI(output))
		if _, err := p.Run(); err != nil {
			fmt.Printf("alas, there's been an error: %v", err)
			os.Exit(1)
		}
		return nil
		// 	ratio, err := intonation.NewRatioFromString(args[0])
		// 	if err != nil {
		// 		return err
		// 	}
		// 	interval := ratio.Approximate12EDOInterval()
		// 	fmt.Println(ratio, "\t", interval)
		//
		// 	if ratioPlay && !ratioQuiet {
		// 		ctx, ctxCancel := context.WithCancel(context.Background())
		// 		defer ctxCancel()
		// 		output := internal.BeepAudioOutput{SampleRate: beep.SampleRate(48000)}
		// 		if ratioInterval {
		// 			intonation.PlayInterval(ratio, ctx, output)
		// 		}
		// 		intonation.PlayChord(ratio, ctx, output)
		// 		if ratioCompare {
		// 			intonation.PlayChord(interval, ctx, output)
		// 		}
		// 	}
		// 	return nil
		// },
	},
}

func init() {
	// ratioCmd.Flags().BoolVarP(&ratioPlay, "play", "p", true, "play the ratio")
	// ratioCmd.Flags().BoolVarP(&ratioCompare, "compare", "c", false, "play the nearest 12-EDO interval as a comparison")
	// ratioCmd.Flags().BoolVarP(&ratioInterval, "interval", "i", false, "play the ratio as a split dyad before playing it as a chord")
	// ratioCmd.Flags().BoolVarP(&ratioQuiet, "quiet", "q", false, "compare ratio with no audio output")
	//
	// ratioCmd.MarkFlagsMutuallyExclusive("play", "quiet")
}
