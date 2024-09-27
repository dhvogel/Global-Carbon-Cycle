package main

import (
	"fmt"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
)

// MichaelisMenten calculates the reaction rate V for a given substrate concentration [S].
func MichaelisMenten(S, Vmax, Km float64) float64 {
	return (Vmax * S) / (Km + S)
}

func main() {
	// Parameters for the Michaelis-Menten equation
	Vmax := 1.0 // Assume VMax is 1
	Km := 0.5 // Michaelis constant is half VMax

	// Create a plot
	p := plot.New()

	p.Title.Text = "Michaelis-Menten"
	p.X.Label.Text = "Substrate Concentration [S]"
	p.Y.Label.Text = "Reaction Rate (V)"

	// Create a line plot for Michaelis-Menten
	points := make(plotter.XYs, 100)
	for i := 0; i < 100; i++ {
		S := float64(i) * 0.05 // Varying substrate concentration
		points[i].X = S
		points[i].Y = MichaelisMenten(S, Vmax, Km)
	}

	// Create a line to represent the data
	line, err := plotter.NewLine(points)
	if err != nil {
		panic(err)
	}

	p.Add(line)

	// Save the plot to a PNG file
	if err := p.Save(6*vg.Inch, 6*vg.Inch, "michaelis_menten.png"); err != nil {
		panic(err)
	}
}
