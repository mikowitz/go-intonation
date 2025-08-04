package cmd

import (
	"fmt"
	"log"

	intonation "github.com/mikowitz/intonation/pkg"
	"github.com/spf13/cobra"
)

var (
	latticeRatios  []string
	latticeIndices []int
)

var latticeCmd = &cobra.Command{
	Use:   "lattice <ratios>",
	Short: "Construct a just intonation lattice and index into it",
	Long:  `Construct a just intonation lattice from comma-separated ratios and index into it using comma-separated indices.`,
	Args:  cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		if len(latticeIndices) > len(latticeRatios) {
			log.Fatal("more indices than dimensions to index into")
		}

		ratios := parseRatios(latticeRatios)
		lattice := intonation.NewLattice(ratios...)

		ratio, err := lattice.At(latticeIndices...)
		if err != nil {
			log.Fatal(err)
		}
		interval := ratio.Approximate12EDOInterval()
		fmt.Println(ratio, "\t", interval)
	},
}

func parseRatios(input []string) []intonation.Ratio {
	output := []intonation.Ratio{}
	for _, i := range input {
		n, err := intonation.NewRatioFromString(i)
		if err != nil {
			log.Fatal(err)
		}
		output = append(output, n)
	}
	return output
}

func init() {
	latticeCmd.Flags().StringSliceVarP(&latticeRatios, "ratios", "r", []string{}, "the ratios used to construct the lattice dimensions (required)")
	latticeCmd.Flags().IntSliceVarP(&latticeIndices, "indices", "i", []int{}, "the indices used to access the lattice (required)")

	latticeCmd.MarkFlagRequired("ratios")
	latticeCmd.MarkFlagRequired("indices")
}
