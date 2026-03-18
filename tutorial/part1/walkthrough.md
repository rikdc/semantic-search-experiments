# Part 1: Foundations - from vectors to embeddings

## Learning objectives

By the end of this part, you will:

- Understand vectors as numerical representations of data
- Implement cosine similarity to measure vector similarity
- Apply similarity calculations to real text using embeddings
- Calculate centroids (average vectors) to represent groups
- Create the `embed analyze` command

## Part 1A: Understanding vectors and similarity

### What are vectors?

A vector is a list of numbers:

```
vec1 = [1, 2, 3]
vec2 = [4, 5, 6]
vec3 = [0.5, 1.2, 3.8]
```

Each number is a coordinate in multi-dimensional space:

- A 2D vector `[3, 5]` is a point in 2D space (x=3, y=5)
- A 3D vector `[1, 2, 3]` is a point in 3D space
- A 1536D vector has 1536 numbers — hard to visualize, but the math is the same

Represent things as vectors and you can compare them mathematically. That's the whole trick.

### Measuring similarity: cosine similarity

Think of two arrows in space. If they point in the same direction, they're similar (even if one is longer). If they point in different directions, they're not. Cosine similarity measures that angle. The result range is -1 to 1, though for embeddings you'll almost always see 0 to 1.

The formula:

```
cosine_similarity(A, B) = (A . B) / (||A|| x ||B||)
```

Where:

- `A . B` = dot product (sum of `a[0]*b[0] + a[1]*b[1] + ...`)
- `||A||` = magnitude of A (its length: `sqrt(a[0]^2 + a[1]^2 + ...)`)
- `||B||` = magnitude of B

In plain terms:

1. Multiply matching positions: `a[0] * b[0]`, then `a[1] * b[1]`, and so on
2. Add those up — that's the dot product
3. Divide by each vector's magnitude so longer vectors don't automatically score higher
4. Result is between -1 and 1

Here it is step by step on two simple vectors:

```
Vector A = [3, 4]
Vector B = [6, 8]   <- Same direction as A, just 2x longer

Step 1: Multiply positions:  3x6=18,  4x8=32
Step 2: Add them up:         18 + 32 = 50
Step 3: Adjust for lengths:  50 / (5 x 10) = 50/50 = 1.0
Result: 1.0 = Perfect similarity
```

B is twice as long as A but points in the same direction, so the score is perfect. This is why cosine works for semantic similarity — it measures direction, not distance. Scale-invariant. A short sentence and a long paragraph on the same topic still score close to 1.0.

| Score | Meaning |
|-------|---------|
| 1.0 | Identical direction |
| 0.9+ | Very similar |
| 0.7-0.9 | Related |
| 0.5-0.7 | Somewhat related |
| < 0.5 | Different |
| 0.0 | Orthogonal (completely independent) |
| -1.0 | Opposite directions |

## Part 1B: Your first implementation - test the math

Implement cosine similarity and test it with simple vectors before introducing embeddings. Get the math right first.

### Step 1: Implement CosineSimilarity

**File**: `math/similarity.go`

```go
func CosineSimilarity(a, b []float32) (float32, error)
```

Algorithm:

1. Validate inputs: check that vectors have the same length and aren't empty
2. Calculate dot product: `dotProduct = sum(a[i] * b[i]) for all i`
3. Calculate magnitude of A: `magA = sqrt(sum(a[i]^2) for all i)`
4. Calculate magnitude of B: `magB = sqrt(sum(b[i]^2) for all i)`
5. Check for zero magnitudes: return error if either is zero (division by zero)
6. Return: `dotProduct / (magA * magB)`

Hints:

- Use `math.Sqrt()` for square roots
- Use `float64` for intermediate calculations to avoid precision loss
- Convert back to `float32` for the return value
- The file already has detailed TODOs with pseudocode

Go to [`start/math/similarity.go`](start/math/similarity.go) and implement `CosineSimilarity` now.

Unit tests are provided in [`start/math/similarity_test.go`](start/math/similarity_test.go). Run them as you go:

```bash
go test ./math/ -v
```

The tests cover edge cases (zero vectors, mismatched dimensions, empty input) and precision — if your diagonal isn't hitting 1.0, the `float32_precision` test will tell you.

### Step 2: Test with simple product vectors

Once you've implemented `CosineSimilarity` and the unit tests pass, try it with a concrete example to build intuition. Add this to `main.go` temporarily or just read through it:

```go
package main

import (
    "embedtutorial/math"
    "fmt"
)

func main() {
    // Imagine these vectors represent products:
    // High values in dimension 0 = "electronics"
    // High values in dimension 1 = "food/organic"
    // High values in dimension 2 = "computing"

    wirelessMouse := []float32{8.5, 1.2, 9.2}     // Electronics, food, computing
    bluetoothKeyboard := []float32{8.3, 1.5, 9.0} // Electronics, food, computing
    organicApple := []float32{1.0, 9.5, 0.5}      // Not electronics!

    fmt.Println("Testing Cosine Similarity")
    fmt.Println("=" + strings.Repeat("=", 50))

    // Compare similar products
    sim1, _ := math.CosineSimilarity(wirelessMouse, bluetoothKeyboard)
    fmt.Printf("Wireless Mouse vs Bluetooth Keyboard: %.4f (similar!)\n", sim1)

    // Compare different products
    sim2, _ := math.CosineSimilarity(wirelessMouse, organicApple)
    fmt.Printf("Wireless Mouse vs Organic Apple:      %.4f (different!)\n", sim2)

    sim3, _ := math.CosineSimilarity(bluetoothKeyboard, organicApple)
    fmt.Printf("Bluetooth Keyboard vs Organic Apple:  %.4f (different!)\n", sim3)
}
```

Expected output:

```
Testing Cosine Similarity
==================================================
Wireless Mouse vs Bluetooth Keyboard: 0.9997 (similar!)
Wireless Mouse vs Organic Apple:      0.2036 (different!)
Bluetooth Keyboard vs Organic Apple:  0.2293 (different!)
```

The two electronics products score ~1.0 because their vectors point in nearly the same direction. The organic apple points somewhere else entirely (~0.20). If those numbers come out right, the function is correct.

## Part 1C: From vectors to embeddings

So far we've been hand-crafting vectors with three dimensions where we decided what each dimension means. Embeddings are vectors that a model produces from text — 1536 dimensions where the model decides what each dimension captures during training.

### What are embeddings?

An embedding is a fixed-size vector generated from text. Similar meanings land near each other in the vector space. OpenAI's `text-embedding-3-small` turns any string into a 1536-element `[]float32`.

```
Text: "coffee"        -> Embedding: [0.123, -0.456, 0.789, ... ] (1536 dimensions)
Text: "espresso"      -> Embedding: [0.134, -0.442, 0.801, ... ] (similar to coffee!)
Text: "automobile"    -> Embedding: [-0.234, 0.678, -0.123, ... ] (very different!)
```

Your `CosineSimilarity` function works on these without modification. The vectors are longer (1536 dimensions instead of 3). The formula is the same.

### OpenAI embeddings API (provided)

The tutorial provides a wrapper around the OpenAI API (you don't need to implement this):

Model: `text-embedding-3-small`

- Dimensions: 1536
- Max input: 8191 tokens (~6000 words)
- Cost: $0.02 per 1M tokens
- Speed: <100ms per request

Alternatives:

- `text-embedding-3-large`: 3072 dimensions, more accurate, more expensive
- Ollama (local, free): `qwen2.5:latest`, `nomic-embed-text`

Usage (already provided in [`../shared/embedder`](../shared/embedder)):

```go
import "embedtutorial/shared/embedder"

client := embedder.NewOpenAIClient("text-embedding-3-small")
vec, err := client.CreateEmbedding("Coffee shop purchase")
// vec is []float32 with 1536 dimensions
```

### Step 3: Apply to real text

Now use `CosineSimilarity` on real text by converting it to embeddings first.

**File**: `cmd/analyze.go` (already stubbed for you)

What to do:

1. Load transaction data from `inputs.json` (financial transactions like "Starbucks coffee $4.50")
2. Generate embeddings for each transaction using the provided API client
3. Calculate pairwise similarities using your `CosineSimilarity` function
4. Display a similarity matrix

The code is almost the same as the product vector example, but instead of hardcoded vectors, you get them from the API.

The stubbed file at [`start/cmd/analyze.go`](start/cmd/analyze.go) has TODOs for:

- Loading transactions from JSON
- Creating embedding client (OpenAI or Ollama)
- Generating embeddings for each transaction
- Calculating and displaying pairwise similarities

Implement the first three sections now (load transactions, create client, generate embeddings). We'll add centroids next.

Expected output:

```
Analyzing Transaction Embeddings
============================================================

Loaded 10 transactions from inputs.json

Using openai with model text-embedding-3-small

Generating embeddings...
  [1/10] Coffee_A -> [0.0123, -0.0456, ...] (1536 dims)
  [2/10] Coffee_B -> [0.0118, -0.0461, ...] (1536 dims)
  ...

=== Pairwise Similarity Matrix ===

                Coffee_A   Coffee_B   Grocery_A  ...
Coffee_A        1.0000     0.9234     0.7123
Coffee_B        0.9234     1.0000     0.7089
Grocery_A       0.7123     0.7089     1.0000
```

Coffee transactions cluster together (~0.92). Cross-category comparisons drop (~0.71). The math works on real text.

## Part 1D: Centroids

Comparing two vectors tells you about a pair. A centroid tells you about a group — it's the average point in the space the group occupies.

### What is a centroid?

Picture a scatter plot of points. The centroid is where you'd stick a pin to balance the whole cluster. For embeddings, it represents the "typical" meaning of the group. Ten coffee transaction embeddings produce a centroid that looks like an average coffee purchase.

How to calculate it: average each dimension across all vectors.

```
Vector 1: [2, 4]
Vector 2: [4, 6]
Vector 3: [6, 8]

Position 0: (2 + 4 + 6) / 3 = 12/3 = 4
Position 1: (4 + 6 + 8) / 3 = 18/3 = 6

Centroid: [4, 6] <- Right in the middle
```

Or more formally:

```
centroid = (v1 + v2 + ... + vN) / N
```

Once you have a centroid, classify new items by comparing them to it with cosine similarity. High similarity to the coffee centroid? Probably coffee.

```
Coffee transactions:
  "Starbucks coffee $4.50"       -> [0.12, -0.45, 0.78, ...]
  "Peet's Coffee latte"          -> [0.13, -0.44, 0.80, ...]
  "Tim Hortons coffee and donut" -> [0.14, -0.43, 0.82, ...]

Centroid (average):              -> [0.13, -0.44, 0.80, ...]

New transaction:
  "Coffee beans from store"      -> [0.12, -0.45, 0.79, ...]
  Similarity to centroid: 0.98   -> Very likely coffee
```

Centroids come back in Part 2 for anomaly detection, and again in Part 6 for semantic profiles.

### Step 4: Implement CalculateCentroid

**File**: `math/similarity.go` (same file as CosineSimilarity)

```go
func CalculateCentroid(vectors [][]float32) ([]float32, error)
```

Algorithm:

1. Validate inputs: check that vectors slice is not empty and all vectors have same length
2. Create result vector of same length as input vectors
3. For each dimension i: sum all `vectors[j][i]`, divide by number of vectors
4. Return the centroid vector

Hints:

- Get dimensions from first vector: `dimensions := len(vectors[0])`
- Create result: `centroid := make([]float32, dimensions)`
- Use nested loops: outer for dimensions, inner for vectors
- The file already has detailed TODOs with pseudocode

Go to [`start/math/similarity.go`](start/math/similarity.go) and implement `CalculateCentroid` now. The same test file has `TestCalculateCentroid` cases — run `go test ./math/ -v` again to verify.

### Step 5: Add centroid analysis to the analyze command

**File**: `cmd/analyze.go` (continue where you left off)

In the same `runAnalyze` function, after displaying the pairwise similarity matrix, add:

1. Calculate centroid from all embeddings
2. For each transaction, calculate similarity to centroid
3. Display results with interpretation (typical vs unusual)

The stubbed file has TODOs for these steps starting around line 150. Look for:

```go
// Step 6: Calculate centroid
// TODO: Calculate the centroid of all embeddings
```

and

```go
// Step 7: Calculate similarity to centroid
// TODO: For each transaction, calculate similarity to centroid
```

Expected output addition:

```text
=== Calculating Centroid ===

Centroid calculated: [0.0134, -0.0445, ...] (1536 dims)
This represents the 'average' or 'typical' transaction

=== Similarity to Centroid ===

Coffee_A:       0.8456 (typical)
Coffee_B:       0.8234 (typical)
Coffee_C:       0.8567 (typical)
Grocery_A:      0.7654 (typical)
Grocery_B:      0.7423 (typical)
Gas_A:          0.7234 (typical)
Gas_B:          0.7189 (typical)
Restaurant_A:   0.7845 (typical)
Restaurant_B:   0.7923 (typical)
Electronics_A:  0.7456 (typical)
```

All transactions land in the 0.7-0.8 range because they're all financial transactions — shared vocabulary pushes the baseline up. In Part 2, z-scores tell you which scores are actually unusual rather than just lower than the others.

## Project structure

```text
part1/start/
├── main.go                 <- CLI entry point (provided)
├── cmd/
│   └── analyze.go         <- Your implementation
├── math/
│   └── similarity.go      <- Your implementation
├── inputs.json            <- Sample data (provided)
└── go.mod                 <- Module definition (provided)
```

Provided: `main.go` (Cobra CLI setup), `../shared/embedder` (OpenAI and Ollama clients), `inputs.json` (sample transactions).

You implement: `math/similarity.go` (vector math) and `cmd/analyze.go` (analysis command).

## Testing your implementation

### 1. Build and run

```bash
cd tutorial/part1/start/
go mod tidy
go build -o embed .
```

### 2. Set up your environment

For OpenAI:

```bash
export OPENAI_API_KEY="sk-your-key-here"
```

For Ollama (local, free):

```bash
ollama pull qwen2.5:latest
export OLLAMA_HOST="http://localhost:11434"
```

### 3. Run the analyze command

```bash
./embed analyze --provider openai --model text-embedding-3-small
```

or

```bash
./embed analyze --provider ollama --model qwen2.5:latest
```

### 4. Expected output

```text
Analyzing Transaction Embeddings
============================================================

Loaded 10 transactions from inputs.json

Using openai with model text-embedding-3-small

Generating embeddings...
  [1/10] Coffee_A -> [0.0123, -0.0456, ...] (1536 dims)
  [2/10] Coffee_B -> [0.0118, -0.0461, ...] (1536 dims)
  ... (8 more)

=== Pairwise Similarity Matrix ===

                Coffee_A   Coffee_B   Grocery_A  Gas_A      ...
Coffee_A        1.0000     0.9234     0.7123     0.6845
Coffee_B        0.9234     1.0000     0.7089     0.6798
Grocery_A       0.7123     0.7089     1.0000     0.7234
Gas_A           0.6845     0.6798     0.7234     1.0000

=== Calculating Centroid ===

Centroid calculated: [0.0134, -0.0445, ...] (1536 dims)
This represents the 'average' or 'typical' transaction

=== Similarity to Centroid ===

Coffee_A:       0.8456 (typical)
Coffee_B:       0.8234 (typical)
Grocery_A:      0.7654 (typical)
Gas_A:          0.7234 (typical)
...
```

### 5. Verify correctness

Key checks:

- Diagonal is exactly 1.0000 (a vector compared to itself)
- Same-category pairs score > 0.90 (Coffee_A vs Coffee_B)
- Cross-category pairs score < 0.75 (Coffee_A vs Gas_A)
- All similarities between 0 and 1 (normalized vectors)
- Centroid similarities in the 0.70-0.90 range for typical data

### 6. Compare with reference solution

If you get stuck:

```bash
cd ../complete/
go run . analyze
```

### Common issues

`OPENAI_API_KEY not set` — Export the environment variable: `export OPENAI_API_KEY="sk-..."`

`dimension mismatch error` — Check that all vectors have same length in `CosineSimilarity()` validation.

`API rate limit exceeded` — Add small delays between API calls or switch to Ollama.

Negative similarities — Check your dot product calculation. With normalized embeddings, values should always be positive.

Diagonal not 1.0 — Your `CosineSimilarity()` implementation probably has rounding issues. Use float64 for intermediate calculations.

`division by zero` — Add magnitude validation in `CosineSimilarity()` before dividing.

## Recap

You implemented two functions (`CosineSimilarity` and `CalculateCentroid`) and wired them into a CLI command that runs against real embeddings. The math that works on three-element product vectors works identically on 1536-dimension text embeddings.

One thing worth internalising: use `float64` for intermediate calculations and `float32` for storage. The precision loss from accumulating `float32` products is small but annoying — your diagonal won't hit exactly 1.0000, and that's a confusing bug to chase later.

## What's next

Part 2 adds statistics to the mix.

**[Part 2: Anomaly detection](../part2/walkthrough.md)**

Z-scores and statistical methods for flagging unusual transactions. Centroid similarity alone can't tell you whether 0.72 is normal or suspicious — the distribution of all the other similarities gives that context.

If you can calculate cosine similarity by hand on small vectors and explain why the diagonal is always 1.0, you're ready for Part 2.

### Challenge exercises (optional)

1. Find most similar pair: modify the analyze command to find and highlight the two most similar transactions
2. Category centroids: calculate separate centroids for each category (Coffee, Grocery, Gas) and compare
3. Simple classifier: implement a function that takes a new transaction and predicts its category based on nearest centroid

See [`questions.md`](questions.md) Q9-Q10 for hints.

---

## Additional resources

- [OpenAI Embeddings Guide](https://platform.openai.com/docs/guides/embeddings)
- [Cosine Similarity Explained](https://en.wikipedia.org/wiki/Cosine_similarity)
- [Cobra CLI Framework](https://github.com/spf13/cobra)
- [Go Modules](https://go.dev/blog/using-go-modules)
