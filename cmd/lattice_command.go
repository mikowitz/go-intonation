package cmd

import (
	"errors"
	"fmt"

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
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(latticeIndices) > len(latticeRatios) {
			return errors.New("more indices than dimensions to index into")
		}

		ratios, err := parseRatios(latticeRatios)
		if err != nil {
			return err
		}
		lattice := intonation.NewLattice(ratios...)

		ratio, err := lattice.At(latticeIndices...)
		if err != nil {
			return err
		}
		interval := ratio.Approximate12EDOInterval()
		fmt.Println(ratio, "\t", interval)
		return nil
	},
}

func parseRatios(input []string) ([]intonation.Ratio, error) {
	output := []intonation.Ratio{}
	for _, i := range input {
		n, err := intonation.NewRatioFromString(i)
		if err != nil {
			return output, err
		}
		output = append(output, n)
	}
	return output, nil
}

func init() {
	latticeCmd.Flags().StringSliceVarP(&latticeRatios, "ratios", "r", []string{}, "the ratios used to construct the lattice dimensions (required)")
	latticeCmd.Flags().IntSliceVarP(&latticeIndices, "indices", "i", []int{}, "the indices used to access the lattice (required)")

	latticeCmd.MarkFlagRequired("ratios")
	latticeCmd.MarkFlagRequired("indices")
}
