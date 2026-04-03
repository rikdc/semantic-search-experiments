# Part 4: Checkpoint questions

Test your understanding of synthetic strings, centroid-based anomaly detection, and embedding context.

---

**Q1**: Before running `embed detect --raw`, predict: will the Anomaly's similarity score be higher or lower than in synthetic mode, and will it still be flagged? Why?

<details>
<summary>Answer</summary>

Higher, and probably not flagged.

In raw mode, `4999.00` is just a number — unusual in magnitude but not dramatically different from other prices in the model's training data (laptops, flights, furniture all cost thousands). The Normal transactions also encode only their amounts: `45.00`, `32.50`, `8.75`, `120.00`. The centroid is built from those numeric representations, not from the full transactional context.

Without merchant names and times, the similarity scores compress toward each other. The anomaly's embedding is no longer anchored to "high-end jewelry store at 3:45 AM" — it's anchored to "a price in the low thousands," which isn't far from "a price in the tens or hundreds." The gap that was 0.17 in synthetic mode (0.83 normal vs 0.66 anomaly) narrows toward noise.

</details>

---

**Q2**: The `FormatTransaction` function puts amount before merchant: `"Transaction: 45.00 USD at Whole Foods Market at 11:00 AM"`. If you reversed the order to merchant-first — `"Whole Foods Market transaction: 45.00 USD at 11:00 AM"` — would you expect the anomaly detection to improve, degrade, or stay roughly the same? Why?

<details>
<summary>Answer</summary>

Likely degrade slightly, though the effect depends on the model.

Embedding models weight earlier tokens more heavily (attention mechanisms can attend to earlier context when processing later tokens). Putting merchant first makes the merchant name the dominant signal, which is useful when merchants are distinctive. But "High-End Jewelry Store" is already distinctive — the problem is that the late-night time and large amount together push the embedding away from normal.

The current order (amount + merchant + time) gives the model "a large amount at an unusual place at an unusual time" as a coherent package. Merchant-first makes the merchant the anchor and treats amount and time as annotations on it. Either order works; the current order puts the quantitative and contextual signals in a sequence that's more consistent with how financial transaction descriptions typically appear in training data.

The safest way to know for certain is to try both and compare the similarity gap between normals and the anomaly.

</details>

---

**Q3**: You add a fifth Normal transaction: `{ "amount": 3200.00, "merchant": "Apple Store", "time": "01:00 PM" }`. The formatted string is `"Transaction: 3200.00 USD at Apple Store at 01:00 PM"`. Predict: will this pull the Normal centroid closer to the Anomaly or push it further away? Will the Anomaly still be flagged at 0.80?

<details>
<summary>Answer</summary>

It will pull the Normal centroid toward the Anomaly — a large amount at a retail electronics store is semantically closer to "high-end jewelry store" than a Starbucks or a gas station is. The centroid will shift in the direction of that large-purchase, retail-store region of embedding space.

Whether the Anomaly is still flagged depends on how much the centroid moves. If the Apple Store transaction is similar enough to the Anomaly that the centroid's similarity to the Anomaly exceeds 0.80, it won't be flagged anymore. This is a real failure mode of centroid-based detection: adding edge-case normals can widen the "normal" region to include actual anomalies.

This is also why Part 5's z-score approach is more robust — it adapts the threshold to the spread of the Normal cluster rather than using a fixed cutoff. If the Apple Store transaction is genuinely normal and the cluster widens, z-scores will widen the acceptable range proportionally. A fixed 0.80 doesn't.

</details>

---

**Q4**: The `TestDetection_SelfSimilarityIsOne` test checks that a vector compared to itself returns exactly 1.0. The implementation uses `float64` for intermediate accumulation. What would happen if `CosineSimilarity` accumulated in `float32` instead, and why does it matter for this part?

<details>
<summary>Answer</summary>

With `float32` accumulation, the dot product and magnitude calculations on a 1536-dimensional vector accumulate rounding error. The diagonal of any similarity matrix (a vector compared to itself) would drift from 1.0 — not by much, but measurably: values like 0.9999997 or 1.0000003 instead of exactly 1.0.

For anomaly detection specifically, this matters because the Normal centroid is compared against the Normal training vectors themselves. If rounding error causes normals to score below 1.0 by enough to creep toward the threshold, you get false positives. With a 0.80 threshold the margin is large enough that it's unlikely to matter in practice, but with a tighter threshold (say, 0.95) float32 drift could flag your own training data as anomalous.

The `float64` pattern is a cheap fix and standard practice in the series. The test exists to make this guarantee explicit and catch any regression if the accumulation type is ever changed.

</details>

---

**Q5**: `CalculateCentroid` computes a simple average of all Normal embeddings. If one of your Normal training examples is actually mislabeled — it's a fraud case that got tagged as Normal — how will it affect the centroid, and what will you observe in the output?

<details>
<summary>Answer</summary>

The centroid will be pulled in the direction of the mislabeled fraud case's embedding. The "normal" reference point moves away from the genuine normal cluster and toward the fraud region.

Two observable effects:
1. The four genuine Normal transactions will score slightly lower against the shifted centroid — their similarity drops because the centroid is no longer squarely in their region.
2. The actual Anomaly may score higher against the centroid, because the centroid has moved toward the fraud region the Anomaly occupies.

In the worst case, the Anomaly's score crosses above 0.80 and it stops being flagged — the mislabeled example has effectively taught the centroid that fraud looks normal.

This is the core limitation of unsupervised centroid-based detection: it trusts its labels completely. One bad label contaminates the reference point. Adding more Normal training examples dilutes the contamination (the centroid is an average over all of them), but a well-labeled dataset is the foundation the whole approach depends on.

</details>

---

**Q6**: The `embed detect` command prints the formatted text it sends to the embedder for each transaction: `[1/5] Normal_1 → "Transaction: 45.00 USD at Whole Foods Market at 11:00 AM"`. Why is printing the formatted string directly useful here, and when would you want to suppress it?

<details>
<summary>Answer</summary>

It makes `FormatTransaction` observable without a debugger. If the similarity scores look wrong — normals scoring low, anomaly scoring high — the first thing to check is whether the strings being embedded match your expectations. If you see `"Transaction: 45.00 USD at  at "` (empty merchant and time), you know immediately that the JSON fields aren't mapping to the struct fields correctly, without having to add print statements or step through a debugger.

You'd want to suppress it when running against large datasets — printing hundreds of formatted strings is noise. In production usage, a `--verbose` flag (or the absence of one) would gate this output. For a tutorial with 5 transactions, always printing it is the right default: the formatted strings are the central concept of the part, and seeing them explicitly confirms that `FormatTransaction` is doing what the post describes.

</details>
