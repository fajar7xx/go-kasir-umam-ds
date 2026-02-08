# Transaction Performance Optimization

## Overview

The `CreateTransaction` method has been optimized to use batch INSERT operations for transaction details, significantly reducing database round-trips.

## Problem

**Before Optimization:**

```go
// Loop-based: N separate INSERT queries
for i := range details {
    tx.QueryRowContext(ctx,
        `INSERT INTO transaction_details (...) VALUES ($1, $2, $3, $4, $5)`,
        transactionID, productID, price, qty, subtotal)
}
```

**Performance:**
- 1 item: 1 INSERT query
- 10 items: 10 INSERT queries
- 100 items: 100 INSERT queries

Each INSERT is a separate database round-trip (network latency + query execution).

## Solution

**After Optimization:**

```go
// Batch INSERT: 1 query for all items
tx.QueryContext(ctx,
    `INSERT INTO transaction_details (...)
     SELECT * FROM UNNEST($1::int[], $2::int[], $3::float8[], $4::int[], $5::float8[])`,
    pq.Array(txIDs), pq.Array(productIDs), pq.Array(prices),
    pq.Array(quantities), pq.Array(subtotals))
```

**Performance:**
- 1 item: 1 INSERT query
- 10 items: 1 INSERT query ✅
- 100 items: 1 INSERT query ✅

## Implementation Details

### PostgreSQL UNNEST Function

UNNEST takes multiple arrays and combines them row-by-row:

```sql
SELECT * FROM UNNEST(
    ARRAY[1, 1, 1],              -- transaction_ids
    ARRAY[101, 102, 103],        -- product_ids
    ARRAY[5000, 10000, 15000],   -- prices
    ARRAY[2, 3, 1],              -- quantities
    ARRAY[10000, 30000, 15000]   -- subtotals
)
-- Produces 3 rows:
-- (1, 101, 5000, 2, 10000)
-- (1, 102, 10000, 3, 30000)
-- (1, 103, 15000, 1, 15000)
```

### Data Preservation

All database-generated values are preserved using `RETURNING`:

```sql
... RETURNING id, transaction_id, product_id, price, quantity,
              subtotal, created_at, updated_at
```

The implementation:
1. Scans all returned rows
2. Maps product names back to details
3. Returns complete `TransactionDetail` structs with DB values

### Code Structure

```go
// 1. Prepare arrays
txIDs := make([]int, len(details))
productIDs := make([]int, len(details))
// ... etc

// 2. Batch insert
rows, err := tx.QueryContext(ctx, insertQuery,
    pq.Array(txIDs), pq.Array(productIDs), ...)

// 3. Scan results
for rows.Next() {
    var d models.TransactionDetail
    rows.Scan(&d.ID, &d.TransactionID, ...)
    finalDetails = append(finalDetails, d)
}

// 4. Restore product names (not in DB)
d.ProductName = productNameMap[d.ProductID]
```

## Performance Impact

### Query Count Reduction

| Items | Queries Before | Queries After | Improvement |
|-------|---------------|---------------|-------------|
| 1     | 1             | 1             | 0%          |
| 10    | 10            | 1             | **90%**     |
| 50    | 50            | 1             | **98%**     |
| 100   | 100           | 1             | **99%**     |

### Real-World Impact

Assuming 5ms per query (network + execution):

| Items | Time Before | Time After | Saved   |
|-------|------------|------------|---------|
| 10    | 50ms       | 5ms        | 45ms    |
| 50    | 250ms      | 5ms        | 245ms   |
| 100   | 500ms      | 5ms        | 495ms   |

**For a transaction with 100 items, we go from 500ms to 5ms - a 100x improvement!**

## When to Use

This optimization is most beneficial when:
- ✅ Transactions have multiple items (5+)
- ✅ Database is remote (network latency matters)
- ✅ High transaction volume

Less impactful when:
- ⚠️ Transactions typically have 1-2 items
- ⚠️ Database is local (localhost)
- ⚠️ Low transaction volume

## Backward Compatibility

- ✅ All existing tests pass
- ✅ API contract unchanged
- ✅ All database values preserved
- ✅ No breaking changes

## Alternative: CreateTransactionOptimal

For even better performance, consider `CreateTransactionOptimal` which also batches:
- Product SELECT queries (single `IN` query)
- Stock UPDATE queries (single `UNNEST` update)
- Row locking with `FOR UPDATE`

See [Transaction Repository Design](transaction-repository-design.md) for details.

## Testing

Run benchmarks to measure performance:

```bash
go test ./internal/repositories -bench=BenchmarkCreateTransaction -benchmem
```

## References

- PostgreSQL UNNEST: https://www.postgresql.org/docs/current/functions-array.html
- lib/pq Arrays: https://pkg.go.dev/github.com/lib/pq#Array
