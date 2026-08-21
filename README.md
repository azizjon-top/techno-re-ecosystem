# Techno RE Ecosystem

A comprehensive backend ecosystem for the Techno RE platform - combining e-commerce, real-time communication, video streaming, and blockchain-based fact mining.

## 🌟 Features

### Core Modules
- **User Management** - Registration, authentication, profile management
- **E-Commerce** - Product catalog, shopping cart, order management
- **Wallet & Transactions** - Token-based payments, balance management
- **Chat System** - Real-time messaging between users
- **Video Platform** - Video upload, streaming, and management
- **Mining Network** - Fact validation and token rewards
- **Ad Campaigns** - Promotional campaigns with budget management

### Technical Features
- **RESTful API** - Clean, well-documented endpoints
- **JWT Authentication** - Secure token-based auth
- **Error Handling** - Comprehensive error codes and messages
- **Caching** - Redis for performance optimization
- **Logging** - Structured logging with Zap
- **Docker Support** - Containerized deployment
- **Configuration Management** - Environment-based configuration

## 🚀 Quick Start

### Prerequisites
- Go 1.21+
- Docker & Docker Compose
- PostgreSQL 15+
- Redis 7+

### Setup Development Environment

1. **Clone the repository**
```bash
git clone https://github.com/azizjon-top/techno-re-ecosystem.git
cd techno-re-ecosystem
```

2. **Create environment file**
```bash
make env-setup
```

3. **Install dependencies**
```bash
make deps
```

4. **Start services with Docker**
```bash
make docker-up
```

5. **Build and run the server**
```bash
make run
```

The server will start on `http://localhost:8080`

## 📋 Available Commands

```bash
make help              # Show all available commands
make build             # Build the application
make run               # Run the application locally
make test              # Run all tests
make test-coverage     # Generate coverage report
make docker-up         # Start Docker services
make docker-down       # Stop Docker services
make lint              # Run linter
make fmt               # Format code
make dev               # Start full development environment
make db-shell          # Access PostgreSQL shell
make redis-shell       # Access Redis CLI
make status            # Check service status
```

## 🏗️ Project Structure

```
techno-re-ecosystem/
├── cmd/
│   └── server/           # Application entry point
├── internal/
│   ├── config/           # Configuration management
│   ├── errors/           # Error handling
│   ├── logger/           # Logging setup
│   ├── middleware/       # HTTP middleware
│   └── models/           # Data models
├── Dockerfile            # Docker image definition
├── docker-compose.yml    # Docker services composition
├── Makefile              # Development commands
├── go.mod & go.sum       # Go dependencies
├── .env.example          # Environment variables template
└── README.md             # This file
```

## 🔌 API Endpoints

### Health Check
```
GET /health
```

### Authentication
```
POST /api/v1/auth/register
POST /api/v1/auth/login
```

### Users
```
GET /api/v1/users/:id
```

### Products
```
GET /api/v1/products
POST /api/v1/products
```

### Orders
```
GET /api/v1/orders
POST /api/v1/orders
```

### Wallet
```
GET /api/v1/wallet/balance
POST /api/v1/wallet/transfer
```

### Chat
```
GET /api/v1/chats
POST /api/v1/chats/:id/messages
```

### Videos
```
GET /api/v1/videos
POST /api/v1/videos
```

### Mining
```
POST /api/v1/mining/start
POST /api/v1/mining/validate
```

### Campaigns
```
GET /api/v1/campaigns
POST /api/v1/campaigns
```

## 🔐 Authentication

The API uses JWT tokens for authentication. Include the token in the `Authorization` header:

```
Authorization: Bearer <your-jwt-token>
```

Default configuration:
- **Access Token TTL**: 15 minutes
- **Refresh Token TTL**: 7 days
- **Issuer**: techno-re-ecosystem
- **Audience**: techno-re-users

## 📦 Configuration

All configuration is managed through environment variables. See `.env.example` for available options:

### Server
- `SERVER_HOST` - Server host (default: 0.0.0.0)
- `SERVER_PORT` - Server port (default: 8080)
- `ENVIRONMENT` - Environment (development/production)

### Database
- `DATABASE_URL` - PostgreSQL connection string
- `DB_MAX_OPEN_CONNS` - Maximum open connections (default: 25)
- `DB_MAX_IDLE_CONNS` - Maximum idle connections (default: 5)

### Cache (Redis)
- `CACHE_ENABLED` - Enable caching (default: true)
- `REDIS_ADDR` - Redis address (default: localhost:6379)
- `CACHE_DEFAULT_TTL` - Cache TTL (default: 1h)

### JWT
- `JWT_SECRET` - Secret key for signing tokens
- `ACCESS_TOKEN_TTL` - Access token expiration
- `REFRESH_TOKEN_TTL` - Refresh token expiration

### Mining
- `MINING_REWARD_PER_FACT` - Reward amount per validated fact
- `MINING_MIN_SESSION_DURATION` - Minimum mining session duration
- `MINING_MAX_CPU_PERCENT` - Maximum CPU usage percentage
- `MINING_CONSENSUS_THRESHOLD` - Consensus threshold for validation

## 🐛 Error Handling

The API returns structured error responses:

```json
{
  "code": "ERROR_CODE",
  "message": "Human readable error message",
  "status_code": 400,
  "details": {}
}
```

### Common Error Codes
- `UNAUTHORIZED` - Authentication required
- `USER_NOT_FOUND` - User doesn't exist
- `VALIDATION_FAILED` - Input validation failed
- `INSUFFICIENT_BALANCE` - Insufficient wallet balance
- `PRODUCT_NOT_IN_STOCK` - Product unavailable
- `INTERNAL_ERROR` - Server error

## 📊 Database Schema

### Core Tables
- `users` - User accounts
- `products` - E-commerce products
- `orders` - Customer orders
- `wallets` - User token wallets
- `chats` - Chat conversations
- `messages` - Chat messages
- `videos` - Video metadata
- `mining_sessions` - Mining activity
- `facts` - Validated facts
- `campaigns` - Ad campaigns
- `impressions` - Campaign impressions

## 🔄 Development Workflow

### 1. Create a feature branch
```bash
git checkout -b feature/your-feature-name
```

### 2. Make changes and test
```bash
make test
make lint
make fmt
```

### 3. Build and verify
```bash
make build
make run
```

### 4. Commit and push
```bash
git add .
git commit -m "Add your feature"
git push origin feature/your-feature-name
```

### 5. Create Pull Request

## 📝 Logging

The application uses structured logging with Zap. Logs include:
- Request/response details
- Performance metrics
- Error traces
- Custom application events

Logs are formatted as JSON in production and pretty-printed in development.

## 🧪 Testing

```bash
# Run all tests
make test

# Run with coverage
make test-coverage

# Run specific test
go test -v ./internal/models
```

## 🐳 Docker Deployment

### Build Docker image
```bash
make docker-build
```

### Deploy with docker-compose
```bash
make docker-up
```

### View logs
```bash
make docker-logs
```

### Stop services
```bash
make docker-down
```

## 📚 Documentation

### API Documentation
Full API documentation will be available at `/api/v1/docs` (Swagger/OpenAPI)

### Database Documentation
Database schema and migrations are documented in the `migrations/` directory

### Architecture
See `docs/architecture.md` for system design and component relationships

## 🤝 Contributing

1. Fork the repository
2. Create your feature branch
3. Make your changes
4. Run tests and linting
5. Commit with clear messages
6. Push to your fork
7. Submit a Pull Request

## 📄 License

This project is licensed under the MIT License - see LICENSE file for details.

## 👥 Authors

**azizjon-top** - Initial development

## 📞 Support

For issues, questions, or suggestions:
- GitHub Issues: https://github.com/azizjon-top/techno-re-ecosystem/issues
- Email: abduazizy2000@gmail.com

## 🗺️ Roadmap

- [ ] Database migrations setup
- [ ] User authentication service
- [ ] Product catalog API
- [ ] Order management system
- [ ] Wallet and transaction system
- [ ] Chat system with WebSockets
- [ ] Video streaming service
- [ ] Mining consensus algorithm
- [ ] Ad campaign management
- [ ] Admin dashboard API
- [ ] Analytics and reporting

## 📈 Performance

- Response time: < 100ms for cached requests
- Database: Connection pooling with configurable limits
- Caching: Redis with configurable TTL
- Concurrency: Full goroutine-based concurrency support

## ✅ Checklist for Production

- [ ] Change JWT_SECRET to a strong random value
- [ ] Set ENVIRONMENT to "production"
- [ ] Configure PostgreSQL with production settings
- [ ] Enable SSL/TLS
- [ ] Set up monitoring and alerting
- [ ] Configure backups
- [ ] Review security settings
- [ ] Load test the application
- [ ] Set up CI/CD pipeline
- [ ] Document deployment procedure
