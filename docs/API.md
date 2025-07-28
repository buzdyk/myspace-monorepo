# API Documentation

## Base URL
`http://localhost:8080`

## Endpoints

### Root Redirect
**GET /** 
- Redirects to `/today`

### Today View
**GET /today**
- Redirects to current date: `/:year/:month/:day`

**GET /:year/:month/:day**
- Returns daily time tracking data for specified date
- **Parameters:**
  - `year` (int): Year (e.g., 2024)  
  - `month` (int): Month (1-12)
  - `day` (int): Day (1-31)
- **Response:**
```json
{
  "date": "2024-01-15",
  "hours": 8.75,
  "running_hours": 0.5
}
```

### Monthly Views
**GET /month**
- Redirects to current month projects: `/:year/:month/projects`

**GET /:year/:month/projects**
- Returns monthly project breakdown
- **Parameters:**
  - `year` (int): Year
  - `month` (int): Month (1-12)
- **Response:**
```json
{
  "year": 2024,
  "month": 1,
  "total_hours": 168.5,
  "projects": [
    {
      "source": "clockify",
      "project_id": "proj123",
      "project_title": "MySpace Development",
      "seconds": 28800,
      "hours": 8.0
    }
  ]
}
```

**GET /:year/:month/calendar**
- Returns calendar view with daily hours
- **Parameters:**
  - `year` (int): Year
  - `month` (int): Month (1-12)
- **Response:**
```json
{
  "year": 2024,
  "month": 1,
  "days": [
    {
      "day": null,
      "hours": null
    },
    {
      "day": 1,
      "hours": 8.5
    }
  ]
}
```

## Data Types

### ProjectTime
```json
{
  "source": "string",        // Time tracker source (clockify, everhour, mayven)
  "project_id": "string",    // Project identifier
  "project_title": "string", // Human-readable project name
  "seconds": 0,              // Time in seconds
  "hours": 0.0,              // Time in hours (calculated)
  "datetime": "2024-01-15T10:00:00Z" // Optional timestamp
}
```

## Error Responses

All endpoints may return error responses in the following format:

```json
{
  "error": "Error message description"
}
```

### HTTP Status Codes
- `200` - Success
- `302` - Redirect
- `400` - Bad Request (invalid parameters)
- `500` - Internal Server Error

## CORS
The API includes CORS headers to allow frontend access:
- `Access-Control-Allow-Origin: *`
- `Access-Control-Allow-Methods: GET, POST, PUT, DELETE, OPTIONS`
- `Access-Control-Allow-Headers: Content-Type, Authorization`