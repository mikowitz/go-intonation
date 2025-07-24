package main

import intonation "github.com/mikowitz/intonation/pkg"

func main() {
	r := intonation.NewRatio(5, 4)
	r.Play()

	i := intonation.NewInterval(4, 12)
	i.Play()
}
