// Look here: https://github.com/golang-standards/project-layout

package main

import (
	"fmt"

	"github.com/FelixDux/imposcg/dynamics"
	"github.com/FelixDux/imposcg/dynamics/parameters"

	"image/color"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/font"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
)

func main() {
	fmt.Println("Impact Oscillator")

	params, errParams := parameters.NewParameters(2.8, 0.0, 0.8, 100)

	if len(errParams) > 0 {
		for _, err := range(errParams) {

			fmt.Println(err.Error())

		}

		return
	}
	
	impactMap, errMap := dynamics.NewImpactMap(*params)

	if errMap != nil {
		fmt.Println(errMap.Error())

		return
	}

	phi := 0.0
	v := 0.0
	data := impactMap.IterateFromPoint(phi, v, 1000)

	scatterData := make(plotter.XYs, len(data.Impacts))

	for i, point := range(data.Impacts) {
		scatterData[i].X = point.Phase
		scatterData[i].Y = point.Velocity
	}

	// Create a new plot, set its title and
	// axis labels.

	// ω σ π

	p := plot.New()
	
	p.Title.Text = fmt.Sprintf("ω = %g, σ = %g, r = %g, Initial impact at (φ = %g, v = %g)", 
		params.ForcingFrequency, params.ObstacleOffset, params.CoefficientOfRestitution, phi, v)
	p.X.Label.Text = "φ"
	p.Y.Label.Text = "v"

	p.X.Max = 1.0

	// Draw a grid behind the data
	p.Add(plotter.NewGrid())

	// Make a scatter plotter and set its style.
	s, err := plotter.NewScatter(scatterData)
	if err != nil {
		panic(err)
	}
	s.GlyphStyle.Color = color.RGBA{R: 0, B: 0, A: 255}
	s.GlyphStyle.Radius = 0.5
	
	p.Add(s)

	// Save the plot to a PNG file.
	scale := font.Length(8)
	if err := p.Save(scale*vg.Inch, scale*vg.Inch, "scatter.png"); err != nil {
		panic(err)
	}
}