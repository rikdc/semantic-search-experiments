# Part 2: Checkpoint questions

Test your understanding of per-category centroids and nearest-centroid classification.

**Q1**: Why do you need per-category centroids instead of a single global centroid for classification?

<details>
<summary>Answer</summary>

A global centroid represents the average of *all* transactions. It can tell you how typical something is overall, but it can't distinguish between categories. Comparing against a single point gives you one similarity score — no basis for choosing coffee over gas.

Per-category centroids give you one reference point per category. You compare the query to each and pick the highest. The global centroid is useful for anomaly detection (Part 3), but classification requires multiple centroids.

</details>

**Q2**: What determines the quality of a category centroid?

<details>
<summary>Answer</summary>

Three factors:

1. **How many examples you have.** A centroid built from twenty transactions is more stable than one built from two. A single outlier barely moves the average when there's enough else to balance it.

2. **How representative they are.** A coffee centroid built entirely from "Starbucks latte" entries won't match "Tim Hortons double-double" as well as one built from a range of purchases.

3. **How coherent the category is.** If "shopping" covers electronics and clothing, the centroid ends up somewhere in between that isn't close to anything specific.

</details>

**Q3**: You have 3 coffee examples and 30 grocery examples. How does this imbalance affect classification?

<details>
<summary>Answer</summary>

The grocery centroid is more stable, built from 10x more data. The coffee centroid might be skewed by any single unusual example.

That said, centroid *position* isn't biased by count. Three examples averaging to a coffee-like point still works. The risk is that those three examples don't cover the full range of coffee transactions, so edge cases ("coffee beans from grocery store") may misclassify.

Nearest-centroid classification just picks the closest centroid, regardless of category size. That simplicity is the appeal, but it means there's no prior probability. More sophisticated classifiers account for class frequency.

</details>

**Q4**: "Starbucks breakfast sandwich" could be coffee or restaurant. How does nearest-centroid handle this?

<details>
<summary>Answer</summary>

It picks whichever centroid is closer and reports the similarity score. There's no "uncertain" output — you always get a winner.

The scores themselves tell you how confident to be. If coffee scores 0.82 and restaurant scores 0.80, that's a near-tie and the result is unreliable. If coffee scores 0.92 and restaurant scores 0.65, you can trust it.

To handle this programmatically, you could:
- Report the margin between the top two scores
- Set a minimum margin threshold below which you flag "uncertain"
- Return the top-N categories with scores instead of just the winner

</details>

**Q5**: How would you add a new category (e.g., "subscription") at runtime without restarting?

<details>
<summary>Answer</summary>

Add entries for the new category to `inputs.json` and re-run the command. The classify command loads inputs.json fresh on every invocation, so new categories appear automatically.

In a production system, you'd:
1. Store centroids in a database or cache rather than recomputing every time
2. Recompute a category's centroid when new training data arrives
3. Add a new category by computing its centroid and inserting it into the lookup

Worth noting: centroids are independent. Adding "subscription" doesn't invalidate the coffee or gas centroids. Just compute the new one and add it to the set.

</details>
