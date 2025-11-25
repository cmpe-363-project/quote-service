# Quote Service

A simple REST API service for managing quotes, built with Go. This service communicates with the author-service to fetch author information.

## Starting the Service

Start database:

```bash
docker compose up mysql -d
```

Create the database:

```bash
docker exec quote-service-mysql mysql -uroot -proot -e "CREATE DATABASE quote_service;"
```

Start the Author Service:

```bash
docker compose up restapi
```

Add sample authors: (TODO: this will be updated with the dummy quote databse being found)

```bash
docker exec quote-service-mysql mysql -uroot -proot quote_service -e "INSERT INTO quotes (id, message, author_id) VALUES (1, 'quote 1', 1), (2, 'quote 2', 2), (3, 'quote 3', 3);"
```


## Endpoints

### GET /api/version
Returns version information for both quote-service and author-service.

**Response:**
```json
{
  "quote-service": "v1.0.1",
  "author-service": "v1.0.0"
}
```

### GET /api/quote/{id}
Returns a specific quote by ID with author information.

**Response:**
```json
{
  "id": 123,
  "message": "lorem ipsum",
  "author": {
    "id": 1,
    "name": "John Doe"
  }
}
```

### GET /api/quote/random
Returns a random quote with author information.

**Response:**
```json
{
  "id": 123,
  "message": "lorem ipsum",
  "author": {
    "id": 1,
    "name": "John Doe"
  }
}
```