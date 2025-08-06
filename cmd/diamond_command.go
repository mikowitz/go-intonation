package cmd

import (
	"fmt"
	"strconv"
	"strings"

	intonation "github.com/mikowitz/intonation/pkg"
	"github.com/spf13/cobra"
)

var diamondFormat string

var diamondCmd = &cobra.Command{
	Use:   "diamond <limits>",
	Short: "Create a otonality/utonality diamond from the provided limits",
	Long:  `Create a otonality/utonality diamond from the provided limits. Limits should be comma-separated integers.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		limits := []uint{}
		limitStrs := strings.Split(args[0], ",")
		for _, ls := range limitStrs {
			l, err := strconv.Atoi(ls)
			if err != nil {
				return err
			}
			limits = append(limits, uint(l))
		}

		d := intonation.NewDiamond(limits...)

		fmt.Println(d.String(intonation.DiamondStringFormat(diamondFormat)))
		return nil
	},
}

func init() {
	diamondCmd.Flags().StringVarP(&diamondFormat, "format", "f", "diamond", "printing format for the diamond: diamond or square")
}
