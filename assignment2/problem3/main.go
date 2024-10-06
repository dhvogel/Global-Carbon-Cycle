package main

import (
	"fmt"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/plotutil"
	"gonum.org/v1/plot/vg"
)

// Function to plot the graph
func plotGraph(dataSets [][]float64, graphName string) {
	p := plot.New()

	for idx, data := range dataSets {
		pts := make(plotter.XYs, len(data))
		for i := range pts {
			pts[i].X = float64(i)
			pts[i].Y = data[i]
		}

		line, err := plotter.NewLine(pts)
		if err != nil {
			panic(err)
		}
		line.Color = plotutil.Color(idx)
		p.Add(line)
		p.Legend.Add(fmt.Sprintf("P0 = %d", (idx+1)*100), line)
	}

	p.Title.Text = "Plant Carbon Over Time"
	p.X.Label.Text = "Time (years)"
	p.Y.Label.Text = "P (Plant Carbon) (in GtC)"
	// To see 500 on the Y axis
	p.Y.Min = 0
	p.Y.Max = 700

	// See every 100 units on Y axis
	p.Y.Tick.Marker = plot.ConstantTicks([]plot.Tick{
		{Value: 0, Label: "0"},
		{Value: 100, Label: "100"},
		{Value: 200, Label: "200"},
		{Value: 300, Label: "300"},
		{Value: 400, Label: "400"},
		{Value: 500, Label: "500"},
		{Value: 600, Label: "600"},
		{Value: 700, Label: "700"},
	})

	if err := p.Save(8*vg.Inch, 4*vg.Inch, graphName); err != nil {
		panic(err)
	}

	fmt.Printf("\nPlot saved to '%s'\n", graphName)
}

func calculatePoints(g float64, K float64, L float64, initialCarbonStock float64, numPoints int) []float64 {
	P := make([]float64, numPoints)
	P[0] = initialCarbonStock

	for i := 1; i < numPoints; i++ {
		P[i] = P[i-1] + g*(1-P[i-1]/K-1/L)*P[i-1]
	}
	return P
}

func main() {
	g := 0.36
	K := 750.0
	L := 3.0
	numPoints := 100

	P1 := calculatePoints(g, K, L, 100, numPoints)
	P2 := calculatePoints(g, K, L, 200, numPoints)
	P3 := calculatePoints(g, K, L, 300, numPoints)
	P4 := calculatePoints(g, K, L, 400, numPoints)
	P5 := calculatePoints(g, K, L, 500, numPoints)
	P6 := calculatePoints(g, K, L, 600, numPoints)
	P7 := calculatePoints(g, K, L, 700, numPoints)

	plotGraph([][]float64{P1, P2, P3, P4, P5, P6, P7}, "plant_carbon.png")
}
