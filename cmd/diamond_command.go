package cmd

import (
	"fmt"

	intonation "github.com/mikowitz/intonation/pkg"
	"github.com/spf13/cobra"
)

var (
	diamondLimits []uint
	diamondFormat string
)

var diamondCmd = &cobra.Command{
	Use:   "diamond <limits>",
	Short: "Create a otonality/utonality diamond from the provided limits",
	Long:  `Create a otonality/utonality diamond from the provided limits. Limits should be comma-separated integers.`,
	Args:  cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		d := intonation.NewDiamond(diamondLimits...)

		fmt.Println(d.String(intonation.DiamondStringFormat(diamondFormat)))
	},
}

func init() {
	diamondCmd.Flags().UintSliceVarP(&diamondLimits, "limits", "l", []uint{1, 5, 3}, "the limits used to construct the diamond")
	diamondCmd.Flags().StringVarP(&diamondFormat, "format", "f", "diamond", "printing format for the diamond: diamond or square")

	diamondCmd.MarkFlagRequired("limits")
}
