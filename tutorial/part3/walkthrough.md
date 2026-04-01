# Part 3: Visualization — PCA and the `embed visualize` command

## Learning objectives

By the end of this part, you will:

- Understand why PCA can compress 1536 dimensions to 2 without destroying cluster structure
- Implement data centering as a required preprocessing step for PCA
- Use `gonum`'s SVD to compute principal components
- Correctly include a query vector in the PCA input (not project it separately)
- Render a category-colored scatter plot with `gonum/plot`

---

## Part 3A: What PCA is doing

### The core idea

PCA finds the directions in which your data varies most. Projected onto those two directions, a 1536-dimensional embedding vector becomes a 2D point. The axes have no inherent meaning — "Principal Component 1" isn't dollars or category type, it's the direction in which your 12 transaction embeddings differ from each other the most.

Two properties make this work for embeddings:

1. **Cluster structure is large-scale structure.** The differences that separate coffee from gas from grocery dominate the variance. The first two components capture most of that signal.
2. **The remaining 1534 components are small.** PCA discards the directions where transactions look most alike, keeping the ones where they diverge.

The result isn't perfect. You lose information. But the clusters your classifier depends on are usually visible.

### Why centering is required

Before computing principal components, subtract the mean of each dimension from all vectors. Skip this and PCA computes axes relative to the origin of the coordinate system — wherever the embedding model places its zero point — rather than relative to the centre of your data. The clusters will appear, but offset and distorted.

Centering is a precondition, not an optimisation. It moves the data from "positioned somewhere in 1536D space" to "centred at the origin," which is the input PCA is designed to receive.

### The index ordering trap

`gonum`'s `svd.VTo(&v)` returns V, a matrix of shape `d × min(n, d)`. The k-th principal component is the k-th **column** of V, so in code: `v.At(j, k)` — row is the dimension index, column is the component index.

The easy mistake is `v.At(k, j)`. That transposes the indices. When `d > n` (which is always true here: 1536 dimensions, 12 training examples), the matrix has 12 columns and `j` goes up to 1535, so it panics immediately. When `n > d`, it produces silently wrong results instead. The access order is always `v.At(j, k)`.

---

## Part 3B: Implement `centerMatrix`

**File**: `math/pca.go`

`centerMatrix` takes `n` vectors of dimension `d` and returns a new `[][]float64` with the per-dimension mean subtracted across all rows.

1. Allocate a `means` slice of length `d` (float64).
2. Loop over all vectors and accumulate dimension-wise sums into `means`.
3. Divide each element of `means` by `n`.
4. Allocate the `centered` output slice.
5. For each vector, for each dimension: `centered[i][j] = float64(data[i][j]) - means[j]`.

`TestCenterMatrix_ZeroMean` checks that summing any dimension across all rows after centering gives zero (within floating-point tolerance). `TestCenterMatrix_KnownValues` checks specific values against hand-calculated output. Run both before moving on:

```bash
cd tutorial/part3/start/
go test ./math/ -run TestCenterMatrix -v
```

---

## Part 3C: Implement `ProjectTo2D`

**File**: `math/pca.go`

`ProjectTo2D` takes `n` vectors of dimension `d`, runs PCA, and returns `n` 2D points as `[][2]float64`. Each projected point is the dot product of the centered input vector with the first two principal components — column `k` of V:

```
points[i][k] = Σ_j  centered[i][j] × V[j][k]
```

1. Call `centerMatrix(data)` to get zero-mean float64 data.
2. Flatten the centered data into a single `[]float64` in row-major order. Row `i` starts at index `i*d`.
3. Build a `*mat.Dense`: `m := mat.NewDense(n, d, flat)`.
4. Declare `var svd mat.SVD` and call `svd.Factorize(m, mat.SVDThin)`. Check the bool return — false means factorization failed.
5. Call `svd.VTo(&v)` to get V.
6. Check V has at least 2 columns: `_, nComponents := v.Dims()`.
7. Allocate `points := make([][2]float64, n)`.
8. Fill points: for each `i` in `0..n`, for each `k` in `0..1`, sum `centered[i][j] * v.At(j, k)` over all `j`.

Three tests cover this. `TestProjectTo2D_OutputDimensions` confirms output length equals input length with no NaN values. `TestProjectTo2D_PreservesClusterSeparation` confirms that two clearly separated 3D clusters stay separated after projection — a wrong implementation collapses or flips them. `TestProjectTo2D_QueryAppendedLast` confirms that appending a query before calling `ProjectTo2D` places it at the last output index, which is the pattern the command uses.

```bash
go test ./math/ -v
```

All six tests should pass.

---

## Part 3D: Implement `SaveEmbeddingPlot`

**File**: `viz/plot.go`

`SaveEmbeddingPlot` takes 2D points, labels, and category strings, and writes a colored scatter plot to `filename`.

1. Create a plot: `p := plot.New()`. Set `p.Title.Text`, `p.X.Label.Text`, `p.Y.Label.Text`.
2. Build a `categoryColors map[string]color.RGBA`. For each unseen category, assign the next color from `palette`. The category `"[query]"` always gets black: `color.RGBA{A: 255}`.
3. For each point `i`:
   - Create `s, err := plotter.NewScatter(plotter.XYs{{X: points[i][0], Y: points[i][1]}})`.
   - Set `s.GlyphStyle.Color = categoryColors[categories[i]]`.
   - Set `s.GlyphStyle.Shape = draw.CircleGlyph{}` and `s.GlyphStyle.Radius = vg.Points(5)`.
   - If `categories[i] == "[query]"`, override to `draw.CrossGlyph{}` and `vg.Points(7)`.
   - Call `p.Add(s)`.
4. Create labels: `plotter.NewLabels(plotter.XYLabels{XYs: xyFromPoints(points), Labels: labels})`. Set `labelsPlotter.Offset = vg.Point{X: 5, Y: 5}`. Call `p.Add(labelsPlotter)`. (`xyFromPoints` is already implemented.)
5. Save: `return p.Save(6*vg.Inch, 6*vg.Inch, filename)`.

`gonum/plot` handles axis scaling automatically — no need to normalise projected coordinates to a canvas range.

---

## Part 3E: Wire up the command

**File**: `cmd/visualize.go`

Loading inputs, initialising the client, and generating embeddings are already done (TODOs A–C). Fill in D, E, and F.

**TODO D — append query vector**:

```go
if query != "" {
    fmt.Printf("\nEmbedding query: '%s'\n", query)
    queryVec, err := client.CreateEmbedding(query)
    if err != nil {
        return fmt.Errorf("failed to embed query: %w", err)
    }
    vectors = append(vectors, queryVec)
    labels = append(labels, "[query]")
    categories = append(categories, "[query]")
}
```

The query must be appended **before** calling `ProjectTo2D` — it needs to participate in axis computation alongside the training vectors, not get projected onto axes built without it.

**TODO E — project to 2D**:

```go
fmt.Println("\nRunning PCA...")
points, err := math.ProjectTo2D(vectors)
if err != nil {
    return fmt.Errorf("PCA failed: %w", err)
}
```

**TODO F — save the plot**:

```go
if err := viz.SaveEmbeddingPlot(output, labels, categories, points); err != nil {
    return fmt.Errorf("failed to save plot: %w", err)
}
fmt.Printf("Wrote %s (%d points)\n", output, len(points))
```

---

## Running it

```bash
cd tutorial/part3/start/
go mod tidy
go build -o embed .

export OPENAI_API_KEY="sk-your-key-here"
# or: export OLLAMA_HOST="http://localhost:11434"

./embed visualize --provider openai --model text-embedding-3-small
# writes output.svg

./embed visualize --provider openai --model text-embedding-3-small \
  --query "flat white from the local café"
# writes output.svg with a black cross marking the query point
```

Open `output.svg` in a browser. You should see four colored clusters. Compare against Part 2's `embed classify` on the same inputs — whatever the classifier separates should appear as spatially distinct groups.

---

## Verifying against the reference

```bash
cd tutorial/part3/complete/
go run . visualize --provider openai --model text-embedding-3-small
```

Point positions may differ slightly (floating-point projection can vary by platform). The cluster topology — which groups sit near or far from each other — should match.

---

## Common issues

**`SVD factorization failed`** — Usually means all vectors are identical or nearly so. Check that your embeddings are being generated correctly and that `centerMatrix` isn't returning all zeros.

**Plot has all points in one spot** — The projection loop may have the indices transposed. Check that you're using `v.At(j, k)` not `v.At(k, j)`.

**Query point not visible** — Confirm the query category string is exactly `"[query]"` (with brackets) in both `labels`, `categories`, and the color check in `SaveEmbeddingPlot`.

**`go mod tidy` fails** — `gonum.org/v1/gonum` and `gonum.org/v1/plot` are new dependencies in this part. Run `go mod tidy` before building.

---

## What's next

**[Part 4: Anomaly Detection](../part4/walkthrough.md)**

Some transactions in this visualization sit noticeably far from their category's centroid. Part 4 makes that precise — anomaly detection, z-scores, and a domain shift to fraud detection. It also forces a harder question: *what you embed* matters as much as how you compare it.
