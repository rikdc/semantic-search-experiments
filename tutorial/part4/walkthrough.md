# Part 4: Anomaly detection — `embed detect`

## Learning objectives

By the end of this part, you will:

- Understand why embedding a raw amount field produces weaker signal than embedding a descriptive sentence
- Implement `FormatTransaction` to build context-rich synthetic strings
- Build a Normal centroid from labeled embeddings and use it as a reference point
- Flag transactions that fall below a fixed similarity threshold

---

## Part 4A: What synthetic strings are doing

### The framing problem

An embedding model encodes meaning from context. `45.00` is a number. The model knows it sits near other small prices, but it has no way to know that this particular 45.00 was spent at a grocery store on a Tuesday morning. Strip context and you get a representation of the number, not the transaction.

A synthetic string assembles the fields the model needs:

```
"Transaction: 45.00 USD at Whole Foods Market at 11:00 AM"
```

Now the model can use everything it knows about Whole Foods, morning grocery shopping, and typical amounts for that context. The representation is grounded in what a $45 Whole Foods transaction means, not in the number 45.

### Field order matters

Models weight earlier tokens more strongly. Put the semantically heavy fields first: amount and merchant define what the transaction is. Time is a modifier: it changes the interpretation of the merchant-amount pair but doesn't define it independently.

`"Transaction: 4999.00 USD at High-End Jewelry Store at 03:45 AM"`: the late-night signal lands on top of a high-value, unusual-merchant foundation.

`"At 03:45 AM, a transaction of 4999.00 USD occurred at High-End Jewelry Store"`: grammatically equivalent, but `03:45 AM` gets the first prominent token position. The representation shifts toward time-of-day first rather than merchant-and-amount first.

### Why the anomaly sits far from the Normal centroid

The four Normal transactions embed near each other: grocery, gas, coffee, retail (everyday merchants, plausible amounts, daytime times). Their centroid sits in the region of embedding space those inputs share.

`"Transaction: 4999.00 USD at High-End Jewelry Store at 03:45 AM"` lands somewhere else entirely. The model has seen jewelry store transactions before, but rarely at 3:45 AM for nearly $5,000. That sentence maps to a different region. Cosine similarity to the Normal centroid reflects that: 0.66 vs the normal cluster's 0.83–0.87.

---

## Part 4B: Implement `FormatTransaction`

Open `math/format.go`. `FormatTransaction` takes a `Transaction` struct and returns a single descriptive sentence for embedding.

The target output for this input:
```go
Transaction{Amount: 45.00, Merchant: "Whole Foods Market", Time: "11:00 AM"}
```
is:
```
"Transaction: 45.00 USD at Whole Foods Market at 11:00 AM"
```

Use `fmt.Sprintf` with the format string `"Transaction: %.2f USD at %s at %s"`. Arguments in order: `t.Amount`, `t.Merchant`, `t.Time`.

`TestFormatTransaction` and `TestFormatTransaction_HighValueLateNight` verify the exact output string, including the two decimal places on the amount. Run them:

```bash
cd tutorial/part4/start/
go test ./math/ -run TestFormatTransaction -v
```

---

## Part 4C: Implement `FormatRaw`

Still in `math/format.go`. `FormatRaw` returns only the dollar amount as a decimal string, with no merchant or time context.

The target output for `Transaction{Amount: 4999.00}` is `"4999.00"`.

Use `fmt.Sprintf("%.2f", t.Amount)`. The two decimal places matter: the tests check the exact string, and consistent formatting makes the raw-vs-synthetic comparison clean.

`TestFormatRaw` and `TestFormatRaw_SmallAmount` verify both a large and a small amount:

```bash
go test ./math/ -run TestFormatRaw -v
```

---

## Part 4D: Run all math tests

```bash
go test ./math/ -v
```

All 7 should pass:

- `TestFormatTransaction`: exact string for a normal transaction
- `TestFormatTransaction_HighValueLateNight`: exact string for the fraud candidate
- `TestFormatRaw`: large amount, two decimal places
- `TestFormatRaw_SmallAmount`: small amount, two decimal places
- `TestDetection_NormalsPassThreshold`: hand-crafted vectors, normals above 0.80, outlier below
- `TestDetection_SelfSimilarityIsOne`: a vector compared to itself returns exactly 1.0
- `TestDetection_CentroidFromOneVector`: centroid of one vector is that vector

`TestDetection_NormalsPassThreshold` is the most important one for this part. It builds a centroid from three hand-crafted normal vectors, then checks that a fourth normal-direction vector passes the threshold and an opposite-direction outlier fails it. No API required. If this passes, the detection pipeline is correct.

---

## Part 4E: Build and run

```bash
go build -o embed .

# Synthetic strings (full context)
./embed detect --provider ollama --model embeddinggemma
# or
./embed detect --provider openai --model text-embedding-3-small

# Raw amounts only (context-stripped)
./embed detect --provider ollama --model embeddinggemma --raw
```

In synthetic mode, the Anomaly row should be flagged. In raw mode, its similarity score should be close to the normals'. The threshold no longer separates them.

If the Anomaly isn't flagged in synthetic mode, check two things: that `FormatTransaction` is producing the correct string (the `[N/N]` output lines show the formatted text), and that the threshold is set to 0.80 (the default).

---

## Verifying against the reference

```bash
cd tutorial/part4/complete/
go build -o embed .
./embed detect --provider ollama --model embeddinggemma
```

Expected output (scores vary slightly by platform and model version):

```
Normal centroid built from 4 transactions. Threshold: 0.80

Label          Category           Similarity
----------------------------------------------
Normal_1       (Normal          ) sim: 0.8420  ✓ normal
Normal_2       (Normal          ) sim: 0.8505  ✓ normal
Normal_3       (Normal          ) sim: 0.8626  ✓ normal
Normal_4       (Normal          ) sim: 0.8307  ✓ normal
Anomaly        (Fraud_Candidate ) sim: 0.6593  ✗ FLAGGED
```

---

## Common issues

**All transactions flagged.** The threshold may be too high for your model. Check the raw similarity scores first with `--threshold 0.0` to see the actual distribution, then set the threshold below the lowest normal score.

**Anomaly not flagged in synthetic mode.** Verify `FormatTransaction` output. The `[N/N]` lines printed during embedding show exactly what string was sent to the API. If the string looks wrong, the test already caught it.

**`Anomaly` flagged in raw mode too.** This occasionally happens with very large amounts ($4,999) and some models. The gap narrows significantly compared to synthetic mode, but may not fully disappear. The point is the narrowing, not a binary flip.

**`no transactions labeled 'Normal' found`.** The `category` field in `inputs.json` must be exactly `"Normal"` (capital N). Check the JSON.

---

## What's next

**[Part 5: Z-scores and statistical thresholds](../part5/walkthrough.md)**

The 0.80 threshold is a guess calibrated against one dataset with one model. Part 5 replaces it with z-scores: flag anything more than 2 standard deviations below the mean similarity of Normal transactions. The threshold becomes a property of your data, not a number you chose by hand.
