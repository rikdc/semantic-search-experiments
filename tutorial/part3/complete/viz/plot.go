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
func SaveEmbeddingPlot(filename string, labels, categories []string, points [][2]float64) error {
	p := plot.New()
	p.Title.Text = "Embedding Space Visualization (PCA)"
	p.X.Label.Text = "Principal Component 1"
	p.Y.Label.Text = "Principal Component 2"

	// Assign a color to each category in the order they are first seen.
	categoryColors := map[string]color.RGBA{}
	colorIdx := 0
	for _, cat := range categories {
		if _, seen := categoryColors[cat]; !seen {
			if cat == "[query]" {
				categoryColors[cat] = color.RGBA{R: 0, G: 0, B: 0, A: 255}
			} else {
				categoryColors[cat] = palette[colorIdx%len(palette)]
				colorIdx++
			}
		}
	}

	// Add one scatter point per item so each gets its own color and glyph style.
	for i, pt := range points {
		c := categoryColors[categories[i]]
		s, err := plotter.NewScatter(plotter.XYs{{X: pt[0], Y: pt[1]}})
		if err != nil {
			return err
		}
		s.GlyphStyle.Color = c
		s.GlyphStyle.Shape = draw.CircleGlyph{}
		s.GlyphStyle.Radius = vg.Points(5)
		if categories[i] == "[query]" {
			s.GlyphStyle.Shape = draw.CrossGlyph{}
			s.GlyphStyle.Radius = vg.Points(7)
		}
		p.Add(s)
	}

	// Add text labels offset slightly from each point.
	labelsPlotter, err := plotter.NewLabels(plotter.XYLabels{
		XYs:    xyFromPoints(points),
		Labels: labels,
	})
	if err != nil {
		return err
	}
	labelsPlotter.Offset = vg.Point{X: 5, Y: 5}
	p.Add(labelsPlotter)

	return p.Save(6*vg.Inch, 6*vg.Inch, filename)
}

func xyFromPoints(points [][2]float64) plotter.XYs {
	xys := make(plotter.XYs, len(points))
	for i, p := range points {
		xys[i].X = p[0]
		xys[i].Y = p[1]
	}
	return xys
}
