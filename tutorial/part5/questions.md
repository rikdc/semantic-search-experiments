# Part 5: Checkpoint questions

Test your understanding of z-scores, statistical thresholds, and their advantages over fixed cutoffs.

---

**Q1**: The Normal similarities in Part 5 are `0.8420, 0.8505, 0.8626, 0.8307`. Before running the code, compute the z-score of a new transaction with similarity 0.81. Is it flagged at a threshold of -2.0?

<details>
<summary>Answer</summary>

Mean = (0.8420 + 0.8505 + 0.8626 + 0.8307) / 4 = 0.8465

Deviations squared: (0.8420−0.8465)² + (0.8505−0.8465)² + (0.8626−0.8465)² + (0.8307−0.8465)²
= 0.00002025 + 0.00001600 + 0.00026082 + 0.00024964 = 0.00054671

Population stddev = sqrt(0.00054671 / 4) = sqrt(0.00013668) ≈ 0.01169

z = (0.81 − 0.8465) / 0.01169 = −0.0365 / 0.01169 ≈ −3.12

Yes, flagged. 0.81 looks close to the Normal cluster (0.83–0.86), but this cluster is tight — 3 standard deviations is a lot of distance in a small space. A fixed threshold of 0.80 would let it through.

</details>

---

**Q2**: You switch the Normal training set from 4 examples to 20, all with similarities ranging 0.82–0.87. The mean stays roughly the same but the standard deviation drops because the cluster is larger and more consistent. What happens to the Anomaly's z-score?

<details>
<summary>Answer</summary>

More negative. Z = (s − μ) / σ — if σ shrinks while the Anomaly's similarity doesn't move, the magnitude grows.

More Normal examples mean a tighter reference distribution. An outlier sitting 16 standard deviations below a well-established cluster is less ambiguous than one sitting 16 standard deviations below four noisy data points. A fixed threshold of 0.80 wouldn't change at all with more training data.

</details>

---

**Q3**: `CalculateMeanStdDev` uses population standard deviation (divide by `n`). If you changed it to sample standard deviation (divide by `n-1`), would the Anomaly still be flagged at -2.0 with 4 training examples?

<details>
<summary>Answer</summary>

Yes. Sample stddev with n=4 is larger than population stddev by sqrt(4/3) ≈ 1.155. A larger denominator produces a less extreme z-score: roughly −16 becomes roughly −13.9. Still nowhere near −2.0.

The choice matters when the anomaly is borderline, the standard deviation is small, or the threshold is tight. For extreme cases like this one it makes no difference. With datasets in the hundreds, the n vs n-1 distinction is negligible anyway.

</details>

---

**Q4**: What happens to the z-score output if you accidentally pass all transaction similarities (including the Anomaly) to `CalculateMeanStdDev` instead of only the Normal ones?

<details>
<summary>Answer</summary>

The mean drops and stddev inflates. The Anomaly's z-score drifts toward zero — possibly above −2.0, no longer flagged.

Normal transactions shift too: with the mean pulled toward the Anomaly, they score above-average relative to the contaminated distribution.

The code compiles fine and gives no error. `TestZScore_NormalsNotFlagged` still passes. `TestZScore_OutlierFlagged` fails if the inflated stddev is large enough. That test exists precisely to catch this class of bug.

</details>

---

**Q5**: The `--zscore-threshold` flag defaults to -2.0. A colleague suggests setting it to -1.0 to catch more potential fraud. What's the tradeoff, and how would you evaluate whether -1.0 is better than -2.0 for a given dataset?

<details>
<summary>Answer</summary>

At −1.0 roughly 1 in 6 normally-distributed transactions falls below the threshold. With a tight Normal cluster that means genuine transactions get flagged. At −2.0 it's about 2.3%. More sensitive, more noise.

To know which is better you need labelled data — transactions where you know which are actually fraud. Run detection at both thresholds, count true positives, false positives, false negatives. The right answer depends on what errors cost you: missed fraud vs flagging a legitimate customer.

Without labelled data, keep it tight. Loosening to −1.0 without anything to validate against produces flags you can't trust.

</details>

---

**Q6**: The walkthrough notes that `CalculateMeanStdDev` should return `(0, 0)` for an empty slice rather than panicking. In what realistic scenario would an empty slice reach this function, and is `(0, 0)` actually a safe return value in that case?

<details>
<summary>Answer</summary>

An empty slice gets there if `inputs.json` has no `Normal` transactions — wrong category label (`"normal"` vs `"Normal"`) being the most likely culprit.

`(0, 0)` isn't safe in the detection loop. Stddev of 0 means every z-score divides by zero: `+Inf` or `NaN`. But `runDetect` checks `len(normalVecs) == 0` and returns an error before anything is computed, so the bad return value is never used.

It's a defensive contract: the function won't panic, but the caller is responsible for not using the output in a broken state. Same pattern as `CalculateCentroid` in Part 2.

</details>
