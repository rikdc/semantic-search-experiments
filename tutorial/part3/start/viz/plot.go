package viz

import (
	"image/color"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
	"gonum.org/v1/plot/vg/draw"
)

var palette = []color.RGBA{
	{R: 220, G: 57, B: 57, A: 255},  // red
	{R: 57, G: 120, B: 220, A: 255}, // blue
	{R: 57, G: 180, B: 80, A: 255},  // green
	{R: 220, G: 150, B: 30, A: 255}, // orange
	{R: 150, G: 57, B: 220, A: 255}, // purple
}

// SaveEmbeddingPlot writes a scatter plot of 2D projected embeddings to filename.
// Labels and categories must be the same length as points.
// Points with category "[query]" are rendered with a distinct cross marker in black.
//
// TODO: implement SaveEmbeddingPlot
// 1. Create a new plot with p := plot.New() and set Title, X.Label, Y.Label.
// 2. Build a categoryColors map: assign palette colors in order of first appearance;
//    the "[query]" category gets black (color.RGBA{A: 255}).
// 3. For each point, create a single-point plotter.NewScatter, set GlyphStyle.Color,
//    GlyphStyle.Shape = draw.CircleGlyph{} (filled), and GlyphStyle.Radius.
//    For "[query]", override Shape = draw.CrossGlyph{} and use a larger radius.
// 4. Add text labels with plotter.NewLabels(plotter.XYLabels{XYs: ..., Labels: labels})
//    and set labelsPlotter.Offset = vg.Point{X: 5, Y: 5}.
// 5. Save with p.Save(6*vg.Inch, 6*vg.Inch, filename).
// Hint: xyFromPoints converts [][2]float64 to plotter.XYs — use it for the labels step.
func SaveEmbeddingPlot(filename string, labels, categories []string, points [][2]float64) error {
	// Suppress unused imports — remove these lines as you implement above.
	_ = plot.New
	_ = plotter.NewScatter
	_ = plotter.NewLabels
	_ = vg.Inch
	_ = draw.CrossGlyph{}
	_ = palette

	panic("not implemented")
}

func xyFromPoints(points [][2]float64) plotter.XYs {
	xys := make(plotter.XYs, len(points))
	for i, p := range points {
		xys[i].X = p[0]
		xys[i].Y = p[1]
	}
	return xys
}
