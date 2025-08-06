package cmd

import (
	"encoding/gob"
	"fmt"
	"os"
	"path"
	"strings"

	intonation "github.com/mikowitz/intonation/pkg"
	"github.com/spf13/cobra"
)

var (
	latticeIndices  []int
	latticeSaveName string
)

var latticeListCmd = &cobra.Command{
	Use:   "list",
	Short: "List saved lattices",
	Long:  `Return a list of saved lattices`,
	Args:  cobra.ExactArgs(0),
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := os.UserConfigDir()
		if err != nil {
			return err
		}
		latticePath := path.Join(cfg, "go-intonation", "lattices")

		items, _ := os.ReadDir(latticePath)
		for _, item := range items {
			filePath := path.Join(latticePath, item.Name())
			file, err := os.Open(filePath)
			if err != nil {
				return err
			}
			defer file.Close()
			decoder := gob.NewDecoder(file)
			l := intonation.Lattice{}
			decoder.Decode(&l)

			fmt.Printf("%s\t\t%s\n", item.Name(), l)
		}

		return nil
	},
}

var latticeSaveCmd = &cobra.Command{
	Use:   "save <ratios>",
	Short: "Construct and store a lattice for future reference",
	Long:  `Construct a just intonation lattice from comma-separated ratios and persist it for future reference`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := os.UserConfigDir()
		if err != nil {
			return err
		}
		latticePath := path.Join(cfg, "go-intonation", "lattices")
		if err = os.MkdirAll(latticePath, 0755); err != nil {
			return err
		}

		var fileName string
		if latticeSaveName == "" {
			fileName = args[0]
		} else {
			fileName = latticeSaveName
		}
		fileName = strings.ReplaceAll(fileName, "/", "-")

		filePath := path.Join(latticePath, fileName)
		file, err := os.Create(filePath)
		if err != nil {
			return err
		}
		defer file.Close()

		latticeRatios := strings.Split(args[0], ",")
		ratios, err := parseRatios(latticeRatios)
		if err != nil {
			return err
		}
		lattice := intonation.NewLattice(ratios...)
		encoder := gob.NewEncoder(file)

		err = encoder.Encode(lattice)
		if err != nil {
			return err
		}
		return nil
	},
}

var latticeLookupCmd = &cobra.Command{
	Use:   "lookup <ratios or name>",
	Short: "Construct a just intonation lattice, or read a saved one, and index into it",
	Long:  `Construct a just intonation lattice from comma-separated ratios, or read a lattice saved via "lattice save", and index into it using comma-separated indices.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := os.UserConfigDir()
		if err != nil {
			return err
		}
		latticePath := path.Join(cfg, "go-intonation", "lattices")

		filePath := path.Join(latticePath, strings.ReplaceAll(args[0], "/", "-"))
		file, err := os.Open(filePath)

		var lattice intonation.Lattice

		if err == nil {
			decoder := gob.NewDecoder(file)
			decoder.Decode(&lattice)
			file.Close()

		} else {
			latticeRatios := strings.Split(args[0], ",")

			ratios, err := parseRatios(latticeRatios)
			if err != nil {
				return err
			}
			lattice = intonation.NewLattice(ratios...)
		}

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

var latticeCmd = &cobra.Command{
	Use:   "lattice",
	Short: "Actions performed with JI ratio lattices",
	Long:  `Actions for working with just intontation ratio lattices`,
}

func init() {
	latticeLookupCmd.Flags().IntSliceVarP(&latticeIndices, "indices", "i", []int{}, "the indices used to access the lattice (required)")
	latticeLookupCmd.MarkFlagRequired("indices")
	latticeCmd.AddCommand(latticeLookupCmd)

	latticeSaveCmd.Flags().StringVarP(&latticeSaveName, "name", "n", "", "the name under which to save the lattice (defaults to a random, readable string otherwise)")
	latticeCmd.AddCommand(latticeSaveCmd)

	latticeCmd.AddCommand(latticeListCmd)
}
