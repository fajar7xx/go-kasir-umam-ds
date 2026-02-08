# Reports API Documentation

## Overview

The Reports API provides endpoints to retrieve sales reports and analytics for the cashier system.

## Endpoints

### GET /api/v1/reports

**Description:** Placeholder endpoint for general reports (future use)

**Response:**
```json
{
  "data": {
    "message": "Report fetched successfully"
  }
}
```

**Status Codes:**
- `200 OK` - Success
- `405 Method Not Allowed` - Invalid HTTP method

---

### GET /api/v1/reports/{period}

**Description:** Get report for a specific time period

**Supported Periods:**
- `hari-ini` - Today's report (Indonesian)
- `today` - Today's report (English)

**Response (Success):**
```json
{
  "data": {
    "total_revenue": 350000.0,
    "total_transactions": 15,
    "best_seller": {
      "name": "Kopi Susu",
      "quantity_sold": 25
    }
  }
}
```

**Response (No Products Sold):**
```json
{
  "data": {
    "total_revenue": 0.0,
    "total_transactions": 0,
    "best_seller": null
  }
}
```

**Status Codes:**
- `200 OK` - Success
- `404 Not Found` - Unsupported period
- `405 Method Not Allowed` - Invalid HTTP method
- `500 Internal Server Error` - Database or processing error

**Error Response:**
```json
{
  "error": {
    "code": "NOT_FOUND",
    "message": "Period not supported. Use: hari-ini, today"
  }
}
```

## Business Logic

### Revenue Calculation
- Sums all `total_amount` from transactions within the date range
- Uses `COALESCE` to return 0 if no transactions found

### Transaction Count
- Counts all transaction records within the date range

### Best Seller
- Aggregates quantity sold per product across all transactions
- Groups by product ID and name
- Orders by total quantity sold (descending)
- Returns top 1 product
- Returns `null` if no products were sold

## Date Handling

### Today's Report
- Start: `YYYY-MM-DD 00:00:00 UTC`
- End: `YYYY-MM-DD+1 00:00:00 UTC` (exclusive)

### Custom Date Range (Future)
- Format: `YYYY-MM-DD`
- End date is inclusive (adds 24 hours internally)
- Validation: end date must be after start date

## Examples

### Get Today's Report
```bash
curl http://localhost:8080/api/v1/reports/hari-ini
```

### Get Today's Report (English)
```bash
curl http://localhost:8080/api/v1/reports/today
```

## Future Enhancements

Planned support for additional periods:
- `week` / `minggu-ini` - This week's report
- `month` / `bulan-ini` - This month's report
- `year` / `tahun-ini` - This year's report
- `custom` with query params - Custom date range

## Testing

Comprehensive test coverage:
- ✅ Handler layer: 8 tests
- ✅ Service layer: 7 tests
- ✅ Repository layer: 5 tests

Test scenarios:
- Success cases (with/without best seller)
- Error handling (database errors, invalid periods)
- Date range validation
- Method not allowed
- Nil pointer safety

## Bug Fixes

### Fixed in Latest Release
1. **Nil Pointer Bug** - Fixed crash when fetching reports due to uninitialized `bestSeller` variable
2. **No Products Sold** - Added graceful handling when no products are sold (returns `null` for `best_seller`)
3. **Error Handling** - Improved error messages and proper HTTP status codes
