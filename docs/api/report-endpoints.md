# Report API Endpoints

## Get Report

Retrieve sales report data with optional date range filtering.

### Endpoint

```
GET /api/v1/reports
```

### Query Parameters

| Parameter | Type | Required | Format | Description |
|-----------|------|----------|--------|-------------|
| start_date | string | conditional | YYYY-MM-DD | Start date for report range |
| end_date | string | conditional | YYYY-MM-DD | End date for report range |

**Note:** Both `start_date` and `end_date` must be provided together, or neither.

### Behavior

1. **No parameters**: Returns today's report
2. **Both parameters**: Returns report for specified date range
3. **One parameter**: Returns error (both required)

### Response Success (200 OK)

```json
{
  "data": {
    "total_revenue": 1500000.0,
    "total_transactions": 45,
    "best_seller": {
      "name": "Mie Ayam",
      "quantity_sold": 60
    }
  }
}
```

### Response Errors

#### Missing Parameters (400 Bad Request)

```json
{
  "error": {
    "code": "MISSING_PARAMETERS",
    "message": "Both start_date and end_date are required"
  }
}
```

#### Invalid Date Format (400 Bad Request)

```json
{
  "error": {
    "code": "INVALID_DATE_FORMAT",
    "message": "invalid date format: start_date must be YYYY-MM-DD format"
  }
}
```

#### Invalid Date Range (400 Bad Request)

```json
{
  "error": {
    "code": "INVALID_DATE_RANGE",
    "message": "end_date must be after start_date"
  }
}
```

### Examples

#### Get today's report
```bash
curl -X GET "http://localhost:8080/api/v1/reports"
```

#### Get report for specific date range
```bash
curl -X GET "http://localhost:8080/api/v1/reports?start_date=2026-01-01&end_date=2026-02-01"
```

## Get Report By Period

Retrieve sales report for predefined periods.

### Endpoint

```
GET /api/v1/reports/{period}
```

### Path Parameters

| Parameter | Type | Required | Values | Description |
|-----------|------|----------|--------|-------------|
| period | string | yes | hari-ini, today | Report period |

### Response

Same as Get Report endpoint.

### Example

```bash
curl -X GET "http://localhost:8080/api/v1/reports/hari-ini"
```
