package main

import (
	"github.com/liuuner/temp-container/colors"
	"github.com/liuuner/temp-container/selector"
)

var col = colors.CreateColors(true)

func main() {
	cfg := selector.Config{}

	items := []selector.Item{
		{
			Display: "GoLang",
			Color:   col.BlueBright,
		},
		{
			Display: "Rust",
			Color:   col.Red,
		},
		{
			Display: "Java",
			Color:   col.Yellow,
		},
		{
			Display: "Python",
			Color:   col.Green,
		},
	}

	sel := selector.New(items, cfg)

	sel.Open()

}
