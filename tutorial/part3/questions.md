# Part 3: Checkpoint questions

Test your understanding of PCA, dimensionality reduction, and embedding visualization.

---

**Q1**: Before running `embed visualize`, predict: will the coffee and gas clusters overlap or sit in clearly separate regions? Why?

<details>
<summary>Answer</summary>

They should sit in clearly separate regions. Coffee transactions share semantic content (café names, drink orders, small dollar amounts) that is quite different from gas transactions (fuel, stations, larger amounts). The embedding model encodes those differences across many dimensions, and PCA — which keeps the directions of maximum variance — will reflect that separation in the 2D projection.

If they overlapped, it would mean the model doesn't distinguish strongly between "Starbucks latte" and "Shell fuel fill-up," which would also show up as near-identical cosine similarity scores in `embed classify`. The visualization and the classifier should agree.

</details>

---

**Q2**: You run `ProjectTo2D` on 12 training vectors. Then you compute the query's embedding separately and try to plot it on the same axes. The query point appears in a completely different region of the plot from where `embed classify` suggests it should land. What went wrong?

<details>
<summary>Answer</summary>

The query was projected using axes computed from a different dataset — the 12 training vectors alone. PCA axes are properties of the specific data matrix you factorize. When you add even one new point, the axes shift.

The fix: append the query vector to the training slice *before* calling `ProjectTo2D`. All 13 vectors participate in computing the axes together, so the query's 2D coordinates are on the same coordinate system as the training points.

This is why `cmd/visualize.go` builds a single `allVectors` slice before calling `ProjectTo2D`, not after.

</details>

---

**Q3**: Your `centerMatrix` function has a bug: it subtracts the mean of the *first* vector from every vector instead of the mean across all vectors. What will the PCA output look like, and will the bug be obvious from the scatter plot?

<details>
<summary>Answer</summary>

The clusters will likely still appear, but in the wrong position and with distorted spacing. Subtracting one vector's values instead of the true per-dimension mean shifts the data to be centred on a specific point (the first vector) rather than on the geometric centre of the dataset.

The bug might not be immediately obvious from the plot — clusters can still look separated. But the axes will describe variance relative to that one vector, not relative to the dataset's centre, so the projection is geometrically incorrect. The `TestCenterMatrix_ZeroMean` test catches this: after correct centering, the sum of each dimension across all rows must be zero.

</details>

---

**Q4**: Looking at the scatter plot for the `inputs.json` from Parts 1 and 2, the restaurant cluster sits between coffee and grocery rather than in its own distant region. Does this mean the classifier is unreliable for restaurant transactions? What would you do to investigate?

<details>
<summary>Answer</summary>

Not necessarily unreliable — "between" in 2D doesn't mean equidistant in 1536D. PCA discards 1534 dimensions; the separating signal for restaurant vs coffee may live in components beyond PC1 and PC2.

To investigate:
1. Run `embed classify` with several restaurant queries and check the margin between the top score and the runner-up. A consistent margin of 0.05+ suggests the classifier is still separating them, even if the 2D plot looks crowded.
2. Check how many training examples each category has. Two restaurant examples produce a weaker centroid than four coffee examples.
3. Add more varied restaurant training data and re-run the visualization. If the restaurant cluster tightens and separates, the original data was the problem. If it stays mixed with coffee, you may have a genuine semantic overlap (breakfast sandwiches, café meals) that a single category can't handle cleanly.

</details>

---

**Q5**: You add a 13th transaction to `inputs.json` — a clear outlier, something like "international wire transfer $9,500." You regenerate the plot. All existing clusters shift slightly from where they were before. Is this a bug?

<details>
<summary>Answer</summary>

No. It's PCA behaving correctly.

PCA axes are computed from the full dataset. Adding one new point changes the covariance structure, which shifts the principal components, which changes the 2D coordinates of every point. This is the same mechanism that requires the query vector to be appended before calling `ProjectTo2D`.

The shift will be small if the new point is far from the existing data (the outlier's direction of variance doesn't dominate) and larger if it's close to the boundary between clusters. The *topology* — which clusters are near or far from each other — should remain stable. If clusters that were separated before are now overlapping, something else is wrong.

This is also why PCA plots shouldn't be compared across different datasets. "Coffee sits at (-0.2, 0.3)" means nothing — it's the relative position of coffee vs gas vs grocery that matters.

</details>

---

**Q6**: The scatter plot colors points by category. If you ran `embed visualize` without the category information — just labels and no colors — what would you lose, and could you recover the clustering information another way?

<details>
<summary>Answer</summary>

You'd lose the ability to see at a glance whether the geometric clusters in the plot correspond to your category labels. Points might form visual groups in the 2D projection, but you'd have no way to confirm those groups match coffee, gas, grocery, etc. without tracing each label manually.

You could partially recover it by cross-referencing labels against the plot by eye for small datasets, or by computing cluster assignments algorithmically (k-means on the 2D coordinates, then comparing to known categories). But the color-by-category mapping is doing real cognitive work — it's the link between the geometry and the meaning.

A related failure mode: if you color by category but the colors don't match the actual category of each point (an off-by-one in the `labels` or `categories` slice), the plot looks wrong in a confusing way. The geometry is correct but the interpretation is inverted. This is why the walkthrough checks that `labels`, `categories`, and `points` all have the same length and are built from the same loop.

</details>
