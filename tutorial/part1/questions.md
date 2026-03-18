# Part 1: Checkpoint questions

Test your understanding of vectors, cosine similarity, and centroids.

## Understanding questions

**Q1**: What is a vector embedding?

<details>
<summary>Answer</summary>

A fixed-size array of numbers that a model generates from text. OpenAI's `text-embedding-3-small` produces a 1536-element `[]float32` for any input string.

The useful property: similar meanings produce similar vectors. "Coffee" and "espresso" land near each other in 1536-dimensional space. "Automobile" does not.

</details>

**Q2**: Why cosine similarity instead of Euclidean distance?

<details>
<summary>Answer</summary>

Cosine measures direction. Euclidean measures distance. For embeddings, meaning is encoded in direction.

```
vec1 = [1, 2, 3]        (length ~ 3.74)
vec2 = [2, 4, 6]        (length ~ 7.48)
```

Euclidean distance: 3.74 (looks different).
Cosine similarity: 1.0 (identical direction).

These vectors point the same way — one is just twice as long. Cosine catches that. Euclidean doesn't. Scale-invariance matters when comparing short text to long text.

</details>

**Q3**: What is a centroid and what is it good for?

<details>
<summary>Answer</summary>

The average of a group of vectors — the centre of mass. Calculate it by averaging each dimension across all vectors:

```
centroid[i] = (v1[i] + v2[i] + ... + vN[i]) / N
```

Ten coffee transaction embeddings produce a centroid that represents a "typical coffee purchase." Compare a new transaction to it with cosine similarity: high score means it looks like the others.

Uses: classification (nearest centroid), anomaly detection (far from centroid = unusual), and semantic profiles (Part 6).

</details>

## Implementation questions

**Q4**: Why check for zero magnitude before dividing in `CosineSimilarity`?

<details>
<summary>Answer</summary>

A zero-magnitude vector is `[0, 0, 0, ..., 0]` — it has no direction. Dividing by zero magnitude either panics (Go integer division) or produces infinity/NaN (floating point). Cosine similarity is mathematically undefined for zero vectors because there's no angle to measure.

In practice, embedding APIs never return zero vectors for real text. But the check costs nothing and prevents a confusing crash if someone passes bad input.

```go
if magnitudeA == 0 || magnitudeB == 0 {
    return 0, fmt.Errorf("cannot calculate cosine similarity with zero-magnitude vector")
}
```

</details>

**Q5**: What happens if you calculate cosine similarity between vectors of different dimensions?

<details>
<summary>Answer</summary>

It's undefined. The dot product requires matching positions:

```
vec1 = [1, 2, 3]       (3 dimensions)
vec2 = [4, 5, 6, 7]    (4 dimensions)

Dot product = 1*4 + 2*5 + 3*6 + ???
```

The implementation should check `len(a) == len(b)` and return an error if they don't match. This comes up when mixing embeddings from different models — `text-embedding-3-small` (1536 dims) and `text-embedding-3-large` (3072 dims) produce incompatible vectors. All embeddings in your system need to come from the same model.

</details>

**Q6**: What should `CalculateCentroid` do with an empty input slice?

<details>
<summary>Answer</summary>

Return an error. You can't average zero vectors — that's division by zero.

```go
if len(vectors) == 0 {
    return nil, fmt.Errorf("cannot calculate centroid of empty vector set")
}
```

A single-element slice is fine: `centroid([v1]) = v1`. The centroid of one vector is itself.

</details>

## Conceptual questions

**Q7**: Two transactions have cosine similarity 0.0. Are they opposites?

<details>
<summary>Answer</summary>

No. 0.0 means orthogonal — perpendicular, independent, unrelated. Not opposite.

```
vec1 = [1, 0, 0]   (points along X-axis)
vec2 = [0, 1, 0]   (points along Y-axis)
cosine = 0.0       (perpendicular)

vec3 = [-1, 0, 0]  (points opposite to vec1)
cosine(vec1, vec3) = -1.0  (opposite)
```

With real embeddings, you'll rarely see negative similarities. Text embeddings tend to live in the positive region of the space, so most scores fall between 0 and 1.

</details>

**Q8**: Why is the diagonal of the similarity matrix always 1.0?

<details>
<summary>Answer</summary>

The diagonal compares each vector to itself. A vector points in exactly its own direction:

```
cosine_similarity(v, v) = (v . v) / (||v|| * ||v||)
                        = ||v||^2 / ||v||^2
                        = 1.0
```

If your diagonal isn't exactly 1.0000, the implementation has a precision bug. This is usually from accumulating `float32` products instead of promoting to `float64` for intermediate calculations.

</details>

## Challenge questions

**Q9**: How would you find the two most similar transactions in the matrix?

<details>
<summary>Answer</summary>

Iterate through the upper triangle (skip the diagonal), track the max:

```go
maxSim := float32(0)
var bestI, bestJ int

for i := 0; i < len(embeddings); i++ {
    for j := i+1; j < len(embeddings); j++ {
        sim, _ := math.CosineSimilarity(embeddings[i], embeddings[j])
        if sim > maxSim {
            maxSim = sim
            bestI, bestJ = i, j
        }
    }
}

fmt.Printf("Most similar: %s and %s (%.4f)\n",
    transactions[bestI].Label, transactions[bestJ].Label, maxSim)
```

With the sample data, expect two Coffee transactions to win.

</details>

**Q10**: How would you use centroids to classify a new transaction by category?

<details>
<summary>Answer</summary>

Calculate a centroid for each known category, then assign the new transaction to whichever centroid it's most similar to:

```go
func Classify(embedding []float32, centroids map[string][]float32) string {
    best := ""
    bestSim := float32(0)

    for category, centroid := range centroids {
        sim, _ := math.CosineSimilarity(embedding, centroid)
        if sim > bestSim {
            bestSim = sim
            best = category
        }
    }

    return best
}
```

```
New transaction: "Tim Hortons coffee and donut"

  Coffee centroid:   0.92  <- highest
  Grocery centroid:  0.68
  Gas centroid:      0.45

Classification: Coffee
```

This is a zero-shot classifier — no training loop, just centroids and a comparison. Part 2 builds on centroids for anomaly detection using z-scores.

</details>

---

## Additional resources

- [OpenAI Embeddings Guide](https://platform.openai.com/docs/guides/embeddings)
- [Cosine Similarity (Wikipedia)](https://en.wikipedia.org/wiki/Cosine_similarity)
- [Cobra CLI Framework](https://github.com/spf13/cobra)
