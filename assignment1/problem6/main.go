package main

import (
	"encoding/csv"
	"fmt"
	"image/color"
	"log"
	"os"
	"strconv"

	"gonum.org/v1/gonum/stat"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
)

type DataPoint struct {
	U           float64
	V           float64
	W           float64
	Temperature float64
	CO2Density  float64
	H2ODensity  float64
}

func readData(fileName string) ([]DataPoint, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.FieldsPerRecord = -1 // this allows variable-length records

	rows, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	var data []DataPoint

	// skip header row
	for _, row := range rows[1:] {
		if len(row) < 6 {
			continue
		}
		u, _ := strconv.ParseFloat(row[0], 64)
		v, _ := strconv.ParseFloat(row[1], 64)
		w, _ := strconv.ParseFloat(row[2], 64)
		temp, _ := strconv.ParseFloat(row[3], 64)
		co2, _ := strconv.ParseFloat(row[4], 64)
		h2o, _ := strconv.ParseFloat(row[5], 64)

		data = append(data, DataPoint{u, v, w, temp, co2, h2o})
	}

	return data, nil
}

func calculateMean(data []DataPoint) (meanW, meanCO2 float64) {
	var sumW, sumCO2 float64

	// for each row sum up the W and CO2Density values
	for _, point := range data {
		sumW += point.W
		sumCO2 += point.CO2Density
	}

	// then divide by how many rows there are
	meanW = sumW / float64(len(data))
	meanCO2 = sumCO2 / float64(len(data))

	return
}

func calculatePerturbations(data []DataPoint, meanW, meanCO2 float64) ([]float64, []float64) {
	var wPrime, cPrime []float64

	// for each row, add the difference in observed w,c minus mean w,c
	for _, point := range data {
		wPrime = append(wPrime, point.W-meanW)
		cPrime = append(cPrime, point.CO2Density-meanCO2)
	}

	// return a slice (array) of all the differences
	return wPrime, cPrime
}

func plotAndFit(wPrime, cPrime []float64, outputFileName string) (slope, intercept, rValue, covariance float64) {
	// Create plot
	p := plot.New()

	p.Title.Text = "CO2 vs Vertical Velocity"
	p.Y.Label.Text = "Vertical Velocity Perturbation (w')"
	p.X.Label.Text = "CO2 Perturbation (C')"

	pts := make(plotter.XYs, len(wPrime))
	for i := range wPrime {
		pts[i].X = cPrime[i]
		pts[i].Y = wPrime[i]
	}

	// Scatter plot
	s, err := plotter.NewScatter(pts)
	if err != nil {
		panic(err)
	}
	p.Add(s)

	// Fit regression line
	slope, intercept = stat.LinearRegression(wPrime, cPrime, nil, false)
	rValue = stat.Correlation(wPrime, cPrime, nil)
	covariance = stat.Covariance(wPrime, cPrime, nil)

	// Add regression line to the plot
	line := plotter.NewFunction(func(x float64) float64 {
		return slope*x + intercept
	})
	line.Color = color.RGBA{R: 255, G: 0, B: 0, A: 255} // set line color to red
	p.Add(line)

	// Save plot to PNG
	if err := p.Save(6*vg.Inch, 6*vg.Inch, outputFileName); err != nil {
		log.Fatal(err)
	}

	return slope, intercept, rValue, covariance
}

func calculateCO2Flux(covariance float64, molarMassCO2 float64) float64 {
	// Flux = covariance * air density * molar mass of CO2 / volume (conversion factor)
	flux := covariance * molarMassCO2 * 1e6
	return flux
}

func createPlotForFile(inputFileName string, outputFileName string) (*float64, error) {
	data, err := readData(inputFileName)
	if err != nil {
		return nil, err
	}
	meanW, meanCO2 := calculateMean(data)
	wPrime, cPrime := calculatePerturbations(data, meanW, meanCO2)

	slope, intercept, rValue, covariance := plotAndFit(wPrime, cPrime, outputFileName)
	fmt.Printf("Stats for %s: \nSlope: %.4f, Intercept: %.4f, R^2: %.4f, Covariance: %.4f\n",
		inputFileName, slope, intercept, rValue, covariance)
	return &covariance, nil
}

func main() {
	// Constants
	molarMassCO2 := float64(1000 / 44)

	daytimeCovariance, err := createPlotForFile("daytime.eddies.csv", "daytime_plot.png")
	if err != nil {
		log.Fatal(err)
	}

	daytimeFlux := calculateCO2Flux(*daytimeCovariance, molarMassCO2)
	fmt.Printf("Daytime CO2 Flux: %.4f micromoles/m^2 s\n\n", daytimeFlux)

	nighttimeCovariance, err := createPlotForFile("nighttime.eddies.csv", "nighttime_plot.png")
	if err != nil {
		log.Fatal(err)
	}

	nighttimeFlux := calculateCO2Flux(*nighttimeCovariance, molarMassCO2)
	fmt.Printf("Nighttime CO2 Flux: %.4f micromoles/m^2 s\n\n", nighttimeFlux)

}
