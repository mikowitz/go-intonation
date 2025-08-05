package intonation

import (
	"strings"
)

type Diamond [][]Ratio

type DiamondStringFormat string

var (
	FormatDiamond DiamondStringFormat = "diamond"
	FormatSquare  DiamondStringFormat = "square"
)

func NewDiamond(limits ...uint) Diamond {
	d := Diamond{}

	for _, u := range limits {
		row := []Ratio{}
		for _, o := range limits {
			row = append(row, NewRatio(o, u))
		}
		d = append(d, row)
	}

	return d
}

func (d Diamond) String(format DiamondStringFormat) string {
	if format == FormatSquare {
		return d.asSquareString()
	}
	return d.asDiamondString()
}

func (d Diamond) asSquareString() string {
	ret := []string{}
	for _, row := range d {
		r := []string{}
		for _, col := range row {
			r = append(r, col.String())
		}
		ret = append(ret, strings.Join(r, "\t"))
	}
	return strings.Join(ret, "\n\n")
}

func (d Diamond) asDiamondString() string {
	l := len(d) - 1
	coordRows := [][][2]uint{}

	for i := l; i >= 0; i-- {
		row := [][2]uint{}
		x := 0
		for j := i; j <= l; j++ {
			row = append(row, [2]uint{uint(x), uint(j)})
			x++
		}
		coordRows = append(coordRows, row)
	}

	for i := 1; i <= l; i++ {
		row := [][2]uint{}
		x := 0
		for j := i; j <= l; j++ {
			row = append(row, [2]uint{uint(j), uint(x)})
			x++
		}
		coordRows = append(coordRows, row)
	}

	ret := []string{}

	for _, row := range coordRows {
		r := []string{}
		prefix := strings.Repeat("\t", len(d)-len(row))

		for _, pair := range row {
			ratio := d[pair[0]][pair[1]]
			r = append(r, ratio.String())
		}

		ret = append(ret, prefix+strings.Join(r, "\t\t"))
	}

	return strings.Join(ret, "\n\n")
}
