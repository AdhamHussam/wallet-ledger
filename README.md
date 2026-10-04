### 6. List Transfers (Paginated)
```http
GET /transfers?page_id=1&page_size=20
```
**Response (200 OK):**
```json
[
  {
    "id": 1,
    "from_account_id": 1,
    "to_account_id": 2,
    "amount": 250,
    "created_at": "2026-10-04T15:16:11.527005+03:00"
  }
]
```

### 7. Get Transfer by ID
```http
GET /transfers/1
```
**Response (200 OK):**
```json
{
  "id": 1,
  "from_account_id": 1,
  "to_account_id": 2,
  "amount": 250,
  "created_at": "2026-10-04T15:16:11.527005+03:00"
}
```
