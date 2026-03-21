# Part 2: Classification - per-category centroids

## Learning objectives

By the end of this part, you will:

- Understand the difference between a global centroid and per-category centroids
- Group embedding vectors by category
- Compute a centroid for each category
- Classify new text by nearest centroid
- Create the `embed classify` command

## Part 2A: Per-category centroids vs global centroid

### Why one centroid isn't enough

In Part 1, you computed a single centroid across all transactions. That centroid represented the "average transaction" and was useful for spotting outliers — anything far from the center of the whole group.

But a global centroid can't classify. If you ask "is this a coffee purchase or a gas purchase?", comparing to the global centroid just tells you how transaction-like something is. You need separate centroids — one per category — so you can ask "which category does this look most like?"

### Per-category centroids

The idea is straightforward: group your training data by category, compute a centroid for each group, then classify new items by finding the nearest centroid.

```
Coffee transactions:
  "Starbucks coffee grande latte $5.75"          -> [0.12, -0.45, 0.78, ...]
  "Tim Hortons double-double and muffin $4.20"   -> [0.13, -0.44, 0.80, ...]
  "Peet's Coffee espresso and pastry $6.50"       -> [0.11, -0.46, 0.77, ...]

Coffee centroid (average):                        -> [0.12, -0.45, 0.78, ...]

Gas transactions:
  "Shell gas station fuel fill-up $52.40"         -> [-0.23, 0.67, -0.12, ...]
  "Chevron premium gasoline $61.25"               -> [-0.22, 0.68, -0.11, ...]

Gas centroid (average):                           -> [-0.225, 0.675, -0.115, ...]
```

New transaction: "Tim Hortons coffee and donut"
- Similarity to coffee centroid: 0.92
- Similarity to gas centroid: 0.45

Classification: **coffee** (highest similarity wins)

This is nearest-centroid classification — no training loop, no gradient descent, no hyperparameters. Just centroids and cosine similarity. It works surprisingly well when categories have distinct semantic signatures.

### When it breaks down

Nearest-centroid classification assumes each category forms a roughly spherical cluster in the embedding space. It struggles when:

- Categories overlap heavily (e.g., "Starbucks breakfast sandwich" sits between coffee and restaurant)
- One category is much more spread out than another
- The training data doesn't represent the category well (too few examples, or biased examples)

For transaction classification at scale, you'd eventually want something more sophisticated. But centroids are a solid baseline and the concepts transfer directly to more advanced approaches.

## Part 2B: Implement grouping logic and test with known vectors

### Step 1: Understand the grouping

Before touching embeddings, think about what the grouping step produces. You have a flat list of transactions, each with a category label. You need a `map[string][][]float32` where keys are category names and values are slices of embedding vectors.

```go
// Input: flat list with categories
transactions := []Transaction{
    {Label: "Coffee_A", Category: "coffee", Text: "..."},
    {Label: "Coffee_B", Category: "coffee", Text: "..."},
    {Label: "Gas_A",    Category: "gas",    Text: "..."},
}

// After embedding and grouping:
groups := map[string][][]float32{
    "coffee": {coffeeAVec, coffeeBVec},
    "gas":    {gasAVec},
}
```

The grouping itself is a simple append loop — the interesting part is what happens next when you compute centroids.

### Step 2: Test centroid logic with known vectors

You already have `CalculateCentroid` from Part 1. Verify it works for the grouping scenario with a quick test:

```go
// Three "coffee" vectors
coffeeVectors := [][]float32{
    {2, 4},
    {4, 6},
    {6, 8},
}

// Expected centroid: [4, 6]
centroid, _ := math.CalculateCentroid(coffeeVectors)
fmt.Printf("Coffee centroid: %v\n", centroid)
```

Then verify classification logic with simple vectors:

```go
coffeeCentroid := []float32{0.9, 0.1, 0.0}
groceryCentroid := []float32{0.1, 0.9, 0.0}
gasCentroid := []float32{0.0, 0.1, 0.9}

query := []float32{0.85, 0.15, 0.05}  // Looks like coffee

// Compare query to each centroid
for name, centroid := range centroids {
    sim, _ := math.CosineSimilarity(query, centroid)
    fmt.Printf("%s: %.4f\n", name, sim)
}
// coffee should win
```

The test file at `complete/math/centroid_test.go` covers these cases plus edge cases (empty input, mismatched dimensions, ties).

### Step 3: Run the tests

```bash
cd tutorial/part2/start/
go test ./math/ -v
```

All five tests should pass before you move on to wiring up the full command.

## Part 2C: Wire to embeddings and run classify

### Step 4: Implement the classify command

**File**: `cmd/classify.go` (stubbed in `start/`)

The stub has TODO markers A through H. Work through them in order:

**A. Load inputs.json** — Same as Part 1. `os.ReadFile` then `json.Unmarshal`.

**B. Initialise embedder client** — Same switch on `--provider` flag as Part 1.

**C. Generate embeddings** — Loop through transactions, call `client.CreateEmbedding` for each. Store the category alongside each embedding.

**D. Group by category** — Build `map[string][][]float32` by appending each embedding to its category's slice.

**E. Compute centroids** — Range over the groups map, call `math.CalculateCentroid` for each.

**F. Embed query** — Call `client.CreateEmbedding` on the positional argument.

**G. Compare to centroids** — Call `math.CosineSimilarity` between the query embedding and each centroid. Collect scores.

**H. Print results** — Sort scores descending, print the table, declare the winner.

### Step 5: Run it

```bash
cd tutorial/part2/start/
go build -o embed .
./embed classify --provider openai --model text-embedding-3-small "Tim Hortons coffee and donut"
```

Expected output:

```
Generating embeddings for training data via openai (text-embedding-3-small)...

Classifying: 'Tim Hortons coffee and donut'

Category Similarities:
----------------------
Category: coffee      | Similarity: 0.9187
Category: restaurant  | Similarity: 0.7312
Category: grocery     | Similarity: 0.7043
Category: gas         | Similarity: 0.6821

Result: 'Tim Hortons coffee and donut' is likely in the [coffee] category.
```

Try some ambiguous queries:

```bash
./embed classify "Starbucks breakfast sandwich"
./embed classify "Costco gas station fill-up"
./embed classify "convenience store snacks and energy drink"
```

### Step 6: Compare with reference

```bash
cd ../complete/
go run . classify "Tim Hortons coffee and donut"
```

## Testing your implementation

### Build and run

```bash
cd tutorial/part2/start/
go mod tidy
go build -o embed .
```

### Set up environment

```bash
export OPENAI_API_KEY="sk-your-key-here"
```

Or for Ollama:

```bash
ollama pull qwen2.5:latest
./embed classify --provider ollama --model qwen2.5:latest "Tim Hortons coffee and donut"
```

### Run unit tests

```bash
cd tutorial/part2/complete/
go test ./math/ -v
```

### Verify correctness

Key checks:

- Coffee queries classify as coffee
- Gas queries classify as gas
- Ambiguous queries (e.g., "Starbucks breakfast sandwich") still pick a reasonable category
- All similarity scores are between 0 and 1
- The highest score wins

### Common issues

`inputs.json not found` — Make sure you run the command from the directory containing inputs.json.

`no embeddings generated` — Check that the API key is set and the provider flag is correct.

`unexpected classification` — With only 3-4 examples per category, some ambiguous queries may misclassify. This is expected — more training data improves centroids.

## Checkpoint questions preview

See [`questions.md`](questions.md) for five checkpoint questions covering:

1. Why per-category centroids instead of a global centroid?
2. What determines the quality of a centroid?
3. How does the number of training examples affect classification?
4. When does nearest-centroid classification fail?
5. How would you handle a new category at runtime?

## Recap

You built a zero-shot classifier using nothing but centroids and cosine similarity. No model training, no labeling pipeline, no hyperparameter tuning. Group embeddings by category, average each group, compare new items to the averages.

The same `CalculateCentroid` and `CosineSimilarity` functions from Part 1 power the whole thing. The math doesn't change — you're just applying it to groups instead of the whole dataset.

## What's next

**[Part 3: Anomaly detection](../part3/walkthrough.md)**

Z-scores and statistical methods for flagging unusual transactions. A centroid tells you *which* category something belongs to — z-scores tell you *how confident* you should be about that assignment.
