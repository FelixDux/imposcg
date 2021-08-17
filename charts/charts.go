package charts

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/FelixDux/imposcg/dynamics/impact"
	"github.com/FelixDux/imposcg/dynamics/parameters"

	"image/color"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/font"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
	"gonum.org/v1/plot/vg/draw"
)

func ImageFile(prefix string) *os.File {
	// First argument ensures we use default tmp directory
	file, err := ioutil.TempFile("", fmt.Sprintf("%s.*.png", prefix))
	if err != nil {
		log.Fatal(err)
	}

	defer os.Remove(file.Name())

	return file
}

type phaseTicks struct{}
func (phaseTicks) Ticks(min, max float64) []plot.Tick {
	return []plot.Tick{{Label: "0", Value: min}, {Label: "π/ω", Value: 0.5*(max-min)}, {Label: "2π/ω", Value: max}}
}

func DOAPlotter(priorityColor color.RGBA, radius font.Length, impactData []impact.SimpleImpact) *plotter.Scatter {
	scatterData := make(plotter.XYs, len(impactData))

	for i, point := range(impactData) {
		scatterData[i].X = point.Phase
		scatterData[i].Y = point.Velocity
	}

	// Make a scatter plotter and set its style.
	s, err := plotter.NewScatter(scatterData)
	if err != nil {
		panic(err)
	}

	s.GlyphStyle.Color = priorityColor
	s.GlyphStyle.Radius = radius
	s.GlyphStyle.Shape = draw.BoxGlyph{}

	return s
}

func ScatterPlotter(priority int, impactData []impact.Impact) *plotter.Scatter {
	scatterData := make(plotter.XYs, len(impactData))

	for i, point := range(impactData) {
		scatterData[i].X = point.Phase
		scatterData[i].Y = point.Velocity
	}

	// Make a scatter plotter and set its style.
	s, err := plotter.NewScatter(scatterData)
	if err != nil {
		panic(err)
	}

	colorLevel := uint8(priority*25)

	s.GlyphStyle.Color = color.RGBA{R: colorLevel, B: colorLevel, A: 255}
	s.GlyphStyle.Radius = 0.5

	return s
}

func PreparePlot(parameters parameters.Parameters, info string) *plot.Plot {

	// Create a new plot, set its title and
	// axis labels.

	p := plot.New()
	
	p.Title.Text = fmt.Sprintf("ω = %g, σ = %g, r = %g, %s", 
		parameters.ForcingFrequency, parameters.ObstacleOffset, parameters.CoefficientOfRestitution, info)
	p.X.Label.Text = "φ"
	p.Y.Label.Text = "v"

	p.X.Max = 1.0

	p.X.Tick.Marker = phaseTicks{}

	return p
}

func ImageToFile(p *plot.Plot) *os.File {
	// Save the plot to a PNG file.
	imageFile := ImageFile("scatter")
	scale := font.Length(6)
	if err := p.Save(scale*vg.Inch, scale*vg.Inch, imageFile.Name()); err != nil {
		panic(err)
	}

	return imageFile
}

func ScatterPlot(parameters parameters.Parameters, impactData [][]impact.Impact, info string) *os.File {

	// Create a new plot, set its title and
	// axis labels.

	p := PreparePlot(parameters, info)
	
	for i, data := range (impactData) {

		s := ScatterPlotter(i, data)

		p.Add(s)
	}

	return ImageToFile(p)
}

func ImpactMapPlot(parameters parameters.Parameters, impactData [][]impact.Impact, phi, v float64) *os.File {
	return ScatterPlot(parameters, impactData, fmt.Sprintf("Initial impact at (φ = %g, v = %g)", phi, v))
}

func DOAPlot(parameters parameters.Parameters, radius font.Length, data *map[string][]impact.SimpleImpact) *os.File {
	p := PreparePlot(parameters, "Domains of attraction")

	p.Legend.ThumbnailWidth = 10
	p.Legend.YPosition = draw.PosCenter
	// p.Legend.XOffs = 60

	priority := int(0)

	picker := DOAColorPicker(len(*data))

	for label, impacts := range *data {
		plt := DOAPlotter(picker(priority), radius, impacts)

		p.Add(plt)

		p.Legend.Add(label, plt)

		priority++
	}

	return ImageToFile(p)
}

func DOAColorPicker(numPriorities int) func(int) color.RGBA {

	colorDelta := uint(3*255 / numPriorities)

	return func(priority int) color.RGBA {

		colorLevel := priority * int(colorDelta)

		OffsetColorLevel := func (offset int) uint8 {
			result := colorLevel - offset * 255
			if result > 255 {
				result = 255
			}
			if result < 0 {
				result = 128
			}
			return uint8(result)
		}

		return color.RGBA{R: OffsetColorLevel(0), B: OffsetColorLevel(1), G: OffsetColorLevel(2), A: 255}
	}
}

func DOARadius(numPhases uint, numVelocities uint) font.Length {
	multiplier := float32(300.0)
	if numPhases > numVelocities {
		return font.Length(multiplier / float32(numPhases))
	} else {
		return font.Length(multiplier / float32(numVelocities))
	}
}