# Event Sourcing & CQRS Learning Repository

A comprehensive Go implementation demonstrating Event Sourcing patterns with CQRS (Command Query Responsibility Segregation). This repository contains three progressive branches that build upon each other to teach core event sourcing concepts.

## 📚 What You'll Learn

- **Event Sourcing**: Store events as the source of truth rather than current state
- **CQRS**: Separate read and write models for better scalability
- **Event Streaming**: Real-time event processing with Apache Kafka
- **Event Store**: Using MinIO as an event store
- **Caching**: Redis for read model optimization
- **Database**: PostgreSQL for projections and read models

## 🌟 Repository Structure

This repository contains three branches, each building upon the previous:

### 1. `simple-eventsourcing-example`

**Basic Event Sourcing Implementation**

A minimal implementation showing core event sourcing concepts:

- Account creation, deposits, and withdrawals as events
- In-memory event store
- Basic event replay for state reconstruction
- Simple REST API for account operations

**Key Concepts:**

- Event-driven architecture
- State reconstruction from events
- Immutable event history

### 2. `eventsourcing-w-kafka-eventstreaming`

**Event Sourcing with Kafka Event Streaming**

Builds upon the simple example by adding:

- Apache Kafka for event streaming
- MinIO as persistent event store
- Redis for caching
- PostgreSQL for projections
- Distributed event processing

**Key Concepts:**

- Event streaming and pub/sub patterns
- Persistent event storage
- Distributed system considerations
- Performance optimization with caching

### 3. `CQRS-Pattern`

**Complete CQRS Implementation**

Full implementation with separated read/write models:

- Command handlers for write operations
- Query handlers for read operations
- Separate read and write databases
- Event projections and view models
- Advanced caching strategies

**Key Concepts:**

- Command/Query separation
- Read model projections
- Eventual consistency
- Scalability patterns

## 🚀 Getting Started

### Prerequisites

```bash
# Install Go (1.19+)
go version

# Install Docker and Docker Compose
docker --version
docker-compose --version
```

### Clone and Setup

```bash
# Clone the repository
git clone https://github.com/vakharwalad23/eventsource-starter-go.git
cd eventsource-starter-go

# Start with the simple example
git checkout simple-eventsourcing-example
```

### Running the Examples

#### Simple Event Sourcing Example

```bash
# Switch to simple example branch
git checkout simple-eventsourcing-example

# Run the application
go run main.go
```

#### Event Sourcing with Kafka

```bash
# Switch to Kafka example branch
git checkout eventsourcing-w-kafka-eventstreaming

# Start infrastructure services
docker-compose up -d

# Run the application
go run main.go

# Test the full flow
./test_api.sh (Not Available for now)
```

#### CQRS Pattern

```bash
# Switch to CQRS branch
git checkout CQRS-Pattern

# Start all services
docker-compose up -d

# Run the application
go run main.go
```

## 🏗️ Architecture Overview

### Event Sourcing Flow

```
[Command] → [Event] → [Event Store] → [Event Stream] → [Projections] → [Read Model]
```

### Components

- **Commands**: Write operations (CreateAccount, Deposit, Withdraw)
- **Events**: Immutable facts (AccountCreated, MoneyDeposited, MoneyWithdrawn)
- **Event Store**: MinIO for persistent event storage
- **Event Stream**: Kafka for real-time event processing
- **Projections**: Event handlers that build read models
- **Read Model**: Optimized views in PostgreSQL and Redis

## 📊 API Examples

### Account Management

```bash
# 1. Create account
curl -X POST http://localhost:8080/accounts \
  -H "Content-Type: application/json" \
  -d '{"ID": "account123"}' \
  -v

# 2. Deposit money
curl -X POST http://localhost:8080/accounts/account123/deposit \
  -H "Content-Type: application/json" \
  -d '{"Amount": 500.00}' \
  -v

# 3. Check balance
curl -X GET http://localhost:8080/accounts/account123/balance \
  -H "Content-Type: application/json" \
  -v

# 4. Withdraw money
curl -X POST http://localhost:8080/accounts/account123/withdraw \
  -H "Content-Type: application/json" \
  -d '{"Amount": 150.00}' \
  -v

# 5. Check final balance
curl -X GET http://localhost:8080/accounts/account123/balance \
  -H "Content-Type: application/json" \
  -v
```

## 🧪 Testing (Not Available for now)

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run specific test package
go test ./internal/app/
go test ./internal/api/
go test ./internal/domain/

# Run integration tests
go test -v integration_test.go
```

## 📋 Key Patterns Demonstrated

### 1. Event Sourcing

- Events as the source of truth
- Immutable event history
- State reconstruction through event replay
- Temporal queries and audit trails

### 2. CQRS (Command Query Responsibility Segregation)

- Separate models for reads and writes
- Command handlers for state changes
- Query handlers for data retrieval
- Independent scaling of read/write sides

### 3. Event Streaming

- Real-time event processing
- Pub/sub messaging patterns
- Event-driven microservices
- Distributed system coordination

### 4. Projections

- Building read models from events
- Eventual consistency
- View materialization
- Denormalization strategies

## 🔧 Technology Stack

- **Language**: Go 1.19+
- **Event Store**: MinIO
- **Message Broker**: Apache Kafka
- **Database**: PostgreSQL
- **Cache**: Redis
- **HTTP Framework**: Gorilla Mux
- **Testing**: Go standard library + sqlmock

## 🏃‍♂️ Learning Path

1. **Start Simple**: Begin with `simple-eventsourcing-example`

   - Understand basic event sourcing concepts
   - See how events rebuild state
   - Learn the core patterns

2. **Add Complexity**: Move to `eventsourcing-w-kafka-eventstreaming`

   - Introduce persistent storage
   - Add event streaming
   - Understand distributed considerations

3. **Full CQRS**: Finish with `CQRS-Pattern`
   - Implement complete separation
   - Build optimized read models
   - Scale read and write independently

## 📈 Benefits of Event Sourcing

- **Audit Trail**: Complete history of all changes
- **Temporal Queries**: Query state at any point in time
- **Debugging**: Replay events to understand system behavior
- **Scalability**: Independent scaling of read/write operations
- **Resilience**: Rebuild state from events if needed

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests
5. Submit a pull request

## 📖 Additional Resources

- [Event Sourcing Pattern](https://learn.microsoft.com/en-us/azure/architecture/patterns/event-sourcing)
- [CQRS Pattern](https://learn.microsoft.com/en-us/azure/architecture/patterns/cqrs)
- [Event Streaming with Kafka](https://kafka.apache.org/documentation/)
- [Go Best Practices](https://golang.org/doc/effective_go.html)

## 🎯 Next Steps

After completing all three branches, consider:

- Implementing event versioning
- Adding event encryption
- Building event-driven sagas
- Creating event sourcing frameworks
- Exploring event store alternatives

---

**Happy Learning!** 🚀

Feel free to open issues for questions or suggestions to improve this learning
