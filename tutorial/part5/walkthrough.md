# Part 5: Z-scores (`embed detect --zscores`)

## Learning objectives

By the end of this part, you will:

- Understand why a fixed similarity threshold doesn't adapt to data spread
- Implement `CalculateMeanStdDev` to compute population mean and standard deviation
- Interpret z-scores as a model-independent measure of anomaly distance
- Use `--zscores` alongside `--zscore-threshold` to replace the hardcoded 0.80 cutoff

---

## Part 5A: What the fixed threshold gets wrong

### The problem with 0.80

The threshold from Part 4 worked because the fraud case was extreme: $4,999 at a jewelry store at 3:45 AM produced a similarity of 0.66, well clear of 0.80. But 0.80 was a judgment call about one dataset with one model. Change either and it may stop working.

**Model shift.** Swap `text-embedding-3-small` for a local embedder and absolute similarity scores shift. Normals that scored 0.88–0.95 now score 0.83–0.87. Your fraud transaction, previously at 0.66, might now score 0.75. Still below your cluster but above 0.80, so no longer flagged. The threshold stayed fixed; the scale moved.

**Data composition shift.** Add a large but legitimate transaction to your Normal training set (say, a $3,200 Apple Store purchase). The Normal centroid shifts toward high-value retail. Now your 0.66 fraud score might become 0.72 because the centroid is closer to it. Same problem: the raw number means nothing without context.

### What z-scores give you

A z-score asks: how many standard deviations is this value from the Normal cluster mean? It adapts as the cluster moves. A transaction 16 standard deviations below normal is anomalous on any model, against any cluster position.

```
z = (s - μ) / σ
```

`μ` and `σ` are the mean and standard deviation of the Normal transactions' similarities. Both are computed from Normal transactions only, so the anomaly is evaluated against the reference distribution, not folded into it.

---

## Part 5B: Implement `CalculateMeanStdDev`

**File**: `math/zscores.go`

`CalculateMeanStdDev` takes a slice of similarity scores (the Normal transactions' similarities to the centroid) and returns `(mean float64, stddev float64)`.

**Step 1: compute the mean:**
Sum all values (converting each `float32` to `float64`), then divide by `n`.

**Step 2: compute the standard deviation:**
Loop again. For each value, compute `diff = float64(s) - mean`, then accumulate `diff * diff`.
Divide the accumulated sum by `n` and take `math.Sqrt`.

This is population standard deviation (divide by `n`). Two passes keeps the implementation clear and avoids floating-point cancellation from single-pass algorithms.

Edge cases to handle:
- Empty slice: return `(0, 0)` immediately.
- Single value: mean = that value, stddev = 0 (the second loop produces 0 and `sqrt(0/1) = 0`).

Run the tests:

```bash
cd tutorial/part5/start/
go test ./math/ -run TestCalculateMeanStdDev -v
```

Four tests should pass: `KnownValues`, `IdenticalValues`, `SingleValue`, and `EmptySlice`.

---

## Part 5C: Verify the z-score pipeline

With `CalculateMeanStdDev` implemented, run the full math suite:

```bash
go test ./math/ -v
```

Nine tests should pass. The two new pipeline tests (`TestZScore_OutlierFlagged` and `TestZScore_NormalsNotFlagged`) use the same hand-crafted similarity values as the post's example:

```go
normalSims := []float32{0.80, 0.84, 0.86, 0.90}  // mean ≈ 0.85, stddev ≈ 0.036
outlierSim := float32(0.60)                         // z ≈ -4.2 → flagged at -2.0
```

These tests verify the math before you touch the API. If they pass, the z-score comparisons in the command loop will be correct.

---

## Part 5D: Build and run

```bash
go build -o embed .

# Fixed-threshold mode (same as Part 4)
./embed detect --provider ollama --model embeddinggemma

# Z-score mode
./embed detect --provider ollama --model embeddinggemma --zscores

# Tighten the threshold
./embed detect --provider ollama --model embeddinggemma --zscores --zscore-threshold -3.0
```

In z-score mode, the output adds a `Z-Score` column. Normal transactions should cluster between -2 and +2; the Anomaly should be sharply negative.

Run fixed-threshold mode first, then z-score mode. The Anomaly row will be flagged in both. The z-score column makes the scale explicit: not just "below threshold" but "16 standard deviations below normal."

---

## Verifying against the reference

```bash
cd tutorial/part5/complete/
go build -o embed .
./embed detect --provider ollama --model embeddinggemma --zscores
```

Expected output (computed from Part 4 Ollama data; verify against your run):

```
Normal centroid built from 4 transactions. Z-score threshold: -2.00

Label          Category           Similarity   Z-Score
------------------------------------------------------
Normal_1       (Normal          ) sim: 0.8420  z:  -0.38  ✓ normal
Normal_2       (Normal          ) sim: 0.8505  z:   0.35  ✓ normal
Normal_3       (Normal          ) sim: 0.8626  z:   1.38  ✓ normal
Normal_4       (Normal          ) sim: 0.8307  z:  -1.35  ✓ normal
Anomaly        (Fraud_Candidate ) sim: 0.6593  z: -16.03  ✗ FLAGGED
```

---

## Common issues

**`stddev = 0` and a division-by-zero z-score** — Happens when all Normal transactions produce identical similarity scores. This won't occur with real embeddings, but add a guard anyway: if `stddev == 0`, return z = 0 or skip the comparison. The `IdenticalValues` test covers this path.

**Normal transactions flagged at -2.0** — With only 4 training examples, the Normal cluster's z-scores span roughly -1.4 to +1.4 (see the expected output above). If your normals are hitting below -2.0, check that `CalculateMeanStdDev` is receiving only the Normal similarities, not all similarities including the anomaly.

**Z-scores look right but the anomaly isn't flagged** — Check that the `--zscore-threshold` default is -2.0. Pass it explicitly to confirm: `--zscore-threshold -2.0`.

---

## What's next

**[Part 6: RAG Baseline](../part6/walkthrough.md)**

Parts 4 and 5 worked with structured transaction data. Part 6 shifts to a codebase: embed Go functions, store them, search with a natural-language query. Part 6 ends at the baseline wall and names it honestly.
