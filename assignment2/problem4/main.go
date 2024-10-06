package main

import (
	"fmt"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/plotutil"
	"gonum.org/v1/plot/vg"
)

func main() {
	// Define constants
	g := 0.36
	K := 750.0
	L := 3.0
	numPoints := 100

	// Initialize the plant carbon array
	P := make([]float64, numPoints)
	P[0] = 100 // Initial plant carbon value

	// Loop to calculate plant carbon over time
	for i := 1; i < numPoints; i++ {
		P[i] = P[i-1] + g*(1-P[i-1]/K-1/L)*P[i-1]
	}

	// Print the results (optional)
	for i, val := range P {
		fmt.Printf("P[%d] = %f\n", i, val)
	}

	// Plotting the values
	plotGraph(P)
}

// Function to plot the graph
func plotGraph(data []float64) {
	pts := make(plotter.XYs, len(data))
	for i := range pts {
		pts[i].X = float64(i)
		pts[i].Y = data[i]
	}

	p := plot.New()

	p.Title.Text = "Plant Carbon Over Time"
	p.X.Label.Text = "Time"
	p.Y.Label.Text = "P (Plant Carbon)"
	// To see 500 on the Y axis
	p.Y.Min = 0
	p.Y.Max = 600

	// See every 100 units on Y axis
	// Add labels at each 100 units
	p.Y.Tick.Marker = plot.ConstantTicks([]plot.Tick{
		{Value: 0, Label: "0"},
		{Value: 100, Label: "100"},
		{Value: 200, Label: "200"},
		{Value: 300, Label: "300"},
		{Value: 400, Label: "400"},
		{Value: 500, Label: "500"},
		{Value: 600, Label: "600"},
	})

	err := plotutil.AddLinePoints(p, "Plant Carbon", pts)
	if err != nil {
		panic(err)
	}

	// Save the plot to a PNG file
	if err := p.Save(8*vg.Inch, 4*vg.Inch, "plant_carbon.png"); err != nil {
		panic(err)
	}

	fmt.Println("Plot saved to 'plant_carbon.png'")
}
