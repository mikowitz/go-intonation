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

// `
// cmd lattice new 3/2,5/4,7/4 --name=trio
// cmd lattice get 1,1,1 --name=trio
// `

var indicesStr string

func LatticeCommand() {
	latticeCmd := flag.NewFlagSet("lattice", flag.ExitOnError)
	latticeCmd.StringVar(&indicesStr, "indices", "", "the indices to access the lattice at")

	if len(os.Args) < 3 {
		fmt.Println("  ratios\n        the ratios used to construct the lattice")
		latticeCmd.PrintDefaults()
		os.Exit(1)
	}

	ratiosStr := os.Args[2]
	latticeCmd.Parse(os.Args[3:])

	ratios := []intonation.Ratio{}
	for _, r := range strings.Split(ratiosStr, ",") {
		ratio, err := intonation.NewRatioFromString(r)
		if err != nil {
			log.Fatal(err)
		}
		ratios = append(ratios, ratio)
	}

	lattice := intonation.NewLattice(ratios...)

	indices := []int{}
	for _, i := range strings.Split(indicesStr, ",") {
		index, err := strconv.Atoi(i)
		if err != nil {
			log.Fatal(err)
		}
		indices = append(indices, index)

	}

	ratio, err := lattice.At(indices...)
	if err != nil {
		log.Fatal(err)
	}
	interval := ratio.Approximate12EDOInterval()
	fmt.Println(ratio, "\t", interval)
}
