# Merch Shop API

A merchandise management system with virtual currency transactions, built with Go, Gin, and GORM.

## Features

- 🔐 JWT Authentication
- 📜 Transaction history tracking
- 🔄 Coin transfers between users
- 📦 Item purchase system
- 🧪 Comprehensive test coverage

## Tech Stack
- **Language** Golang
- **Framework**: [Gin](https://gin-gonic.com/)
- **ORM**: [GORM](https://gorm.io/)
- **Database**: PostgreSQL
- **Authentication**: JWT
- **Testing**: testify

## Prerequisites
- Docker
- Docker Compose

## Setup and Running

1. Clone the repository:
   ```
   git clone https://github.com/VIXI1703/merch_shop.git
   cd merch_shop
   ```

2. Build and run the application using Docker Compose:
   ```
   docker-compose up --build
   ```

3. The application will be available at `http://localhost:8080`

## API Endpoints

### Authentication
#### POST `/api/auth`
```json
{
  "username": "alice",
  "password": "password"
}
```

### Get User Info
#### GET `/api/info`
Requires JWT in Authorization header.

### Send Coins
#### POST `/api/sendCoin`
```json
{
  "toUser": "bob",
  "amount": 150
}
```

### Buy Item
#### POST `/api/buy/{item-name}`

---

## Configuration Options
Environment variables for `docker-compose.yml`:

| Variable          | Default | Description             |
|-------------------|---------|-------------------------|
| `JWT_SIGNING_KEY` | ~       | JWT encryption secret   |
| `JWT_DURATION`    | 24h     | Token validity duration |
| `DB_HOST`         | ~       | Database host           |
| `DB_PORT`         | ~       | Database port           |
| `DB_USER`         | ~       | Database username       |
| `DB_PASSWORD`     | ~       | Database password       |
| `DB_NAME`         | ~       | Database table name     |
| `HTTP_PORT`       | ~       | Http server port        |

---

## Tests:

### Unit Tests:
```bash
go test ./...
```

### Integration Tests:
```bash
go test ./tests/integration/...
```

### E2E Tests:
```bash
go test ./tests/e2e/...
```

---


## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.



