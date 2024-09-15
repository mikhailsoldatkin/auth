# Auth Server

![workflow](https://github.com/mikhailsoldatkin/auth/actions/workflows/go.yaml/badge.svg)

- **[Remote Server](http://51.250.32.78)**
- **[Local Server](http://127.0.0.1)**

## Features

- **gRPC Support**: Secure and scalable communication with your services via gRPC.
- **HTTP/REST Support with gRPC-Gateway**: Expose gRPC services as RESTful endpoints for broader compatibility.
- **Postgres**: Relational database for persistent storage, with support for pgx, goose, and squirrel for advanced database operations.
- **Redis (Cache)**: In-memory data store for caching and fast data retrieval.
- **Docker Compose**: Simplify the setup and deployment of the server and its dependencies using Docker Compose.
- **Kafka**: Message broker for scalable and reliable messaging between services.
- **Swagger API Documentation**: Explore and test RESTful endpoints through Swagger UI.
- **Zap Logger**: High-performance structured logging with Zap for better traceability and debugging.
- **Prometheus Monitoring**: Integrated metrics for real-time performance and health monitoring.
- **Distributed Tracing with Jaeger**: Monitor and trace the behavior of distributed systems using Jaeger.
- **GitHub Actions**: Continuous integration and deployment pipelines for automated testing and deployment.
- **TLS**: Secure communication between services using TLS.
- **JWT (User Authentication)**: JSON Web Tokens for secure user authentication and authorization.
- **Rate Limiting**: Configurable rate limiting per request using a token bucket algorithm.
- **Circuit Breaker**: Mechanism to prevent system failures by handling service disruptions gracefully.

## Local Deployment

To deploy the server locally using Docker Compose, follow these steps:

1. **Ensure Docker and Docker Compose are installed** on your machine.

2. **Clone the repository**:

    ```bash
    git clone https://github.com/mikhailsoldatkin/auth.git
    cd auth
    ```

3. **Create a `.env` file** in the root directory of the project with necessary environment variables. Example:

4. **Start the services**:

    ```bash
    docker compose -f docker-compose-dev.yaml up -d
    ```

## Monitoring and Observability

### Prometheus
The server exposes metrics compatible with Prometheus, available at:

- **Remote:** [Prometheus Metrics (Remote)](http://51.250.32.78:9090)
- **Local:** [Prometheus Metrics (Local)](http://127.0.0.1:9090)

### Grafana
Visualize your metrics with Grafana dashboards:

- **Remote:** [Grafana (Remote)](http://51.250.32.78:3000)
- **Local:** [Grafana (Local)](http://127.0.0.1:3000)

### Jaeger
For distributed tracing, access Jaeger here:

- **Remote:** [Jaeger (Remote)](http://51.250.32.78:16686)
- **Local:** [Jaeger (Local)](http://127.0.0.1:16686)

## API Documentation

Interactive API documentation is available via Swagger:

- **Remote:** [Swagger UI (Remote)](http://51.250.32.78:8090/swagger/index.html)
- **Local:** [Swagger UI (Local)](http://127.0.0.1:8090/swagger/index.html)
