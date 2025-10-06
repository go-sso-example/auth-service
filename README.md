# Auth Service — SSO API Gateway

Auth Service is a Single Sign-On (SSO) service that functions as an API gateway for all registered services and their endpoints. It centralizes authentication, authorization, and access management, providing a secure and consistent mechanism for controlling access across a microservices architecture.

Auth Service implements Single Sign-On (SSO) using OpenID Connect (OIDC) and OAuth 2.0 protocols. This allows users to authenticate once and gain access to all integrated services without repeatedly entering credentials.


## Key Components

- Identity Provider (IdP): Auth Service itself, responsible for authenticating users.
- Service Providers (SPs): Other services, such as User Service and Resource Service, which rely on Auth Service for user authentication.

## Protocols

- OAuth 2.0: Used for authorization and delegated access.
- OpenID Connect (OIDC): An extension of OAuth 2.0 that provides authentication and user information.

## Features

### Core Functionality
The core functionality of Auth Service revolves around centralized authentication, authorization, and access control for all registered services and their endpoints. It ensures that only properly authenticated and authorized requests are forwarded to downstream services. Key components include:

- **Request Validation and Forwarding**:
  - Every incoming request is checked against the registered services and their resources in the database.
  - If the requested service or resource does not exist, the request is rejected immediately.
  - If all checks pass, the request is forwarded to the corresponding service.
- **Access Token Handling**: 
  - Access tokens are extracted from cookies and validated. 
  - Token expiration is checked, and, if required, the access token is rotated using a refresh token. 
  - Invalid or expired tokens result in an unauthorized response.
- **Role-Based Authorization**: 
  - User roles are verified against the roles required for the requested resource. 
  - Only users with appropriate permissions can access specific services or endpoints.
- **Login and Logout Management**:
  - Login requests are proxied to the User Service via gRPC, which stores user credentials (login/password). 
  - Successful logins generate access and refresh tokens that are returned in cookies. 
  - Logout requests invalidate the tokens locally to terminate the user session.
- **Token Refresh Flow**: 
  - Refresh tokens are validated and used to issue new access tokens without requiring a new login. 
  - The system ensures token rotation for enhanced security.
- **CRUD Operations for Services and Resources**:
  - The service provides endpoints to create, read, update, and delete registered services and their resources. 
  - All changes are persisted in the database and reflected in the in-memory cache for faster access.

### API Endpoints
- `POST /login` - Handles user login. Credentials are forwarded to the User Service via gRPC. Returns access and refresh tokens in cookies.
- `POST /logout` - Logs the user out by invalidating access and refresh tokens.
- `POST /refresh-token` - Refreshes the access token using a valid refresh token. Handles token rotation.
- `GET /services` - Retrieves a list of all registered services.
- `POST /services` - Creates a new service and registers it in the database.
- `PUT /services/{id}` - Updates an existing service by ID.
- `DELETE /services/{id}` - Deletes a service by ID.
- `GET /resources` - Retrieves all registered resources/endpoints for services.
- `POST /resources` - Creates a new resource for a service.
- `PUT /resources/{id}` - Updates an existing resource by ID.
- `DELETE/resources/{id}` - Deletes a resource by ID.
- `ANY /{service}/{resource}` - Proxy endpoint: validates access tokens and user roles, then forwards requests to the corresponding service if authorized.

## Architecture

### Project Structure
```
auth-service/
├── api/                   # OpenAPI specifications
├── cmd/                   # Main executable
│   └── auth-service/      # Entry point
├── config/                # Configuration files
│   └── config.yaml        # Main configuration
├── internal/              # Internal packages
├── migrations/            # Database migrations
├── pkg/                   # Packages with business logic
│   └── api/               # Generated API packages
├── .gitignore             # Ignored files
├── Dockerfile             # Docker image
├── Makefile               # Automation scripts
├── README.md              # This file
├── docker-compose.yaml    # Local development setup
├── go.mod                 # Go module
├── go.sum                 # Dependency hashes
└── migrations.sh          # Script to run migrations
```

### Backend Services
- **gRPC Server** with HTTP/REST gateway
- **PostgreSQL** database with connection pooling
- **Modular Service Architecture**:
    - Auth Service
    - Resource Service
    - User Client (gRPC)
    - Token Management
    - Proxy Handler
    - Worker


## Technology Stack

### Backend
- **Go 1.24.2**
- **PostgreSQL** with pgx driver
- **OpenAPI** for REST API specs.
- **oapi-codegen** for generating Go types and server code from OpenAPI specs
- **Redis** for caching services and resources
- **Chi** for HTTP routing
- **JWT** for access and refresh tokens
- **Viper** for configuration management
- **Zap** for structured logging

### Infrastructure
- **Docker & Docker Compose**
- **OpenAPI & oapi-codegen**
- **PostgreSQL 16**
- **Redis Cache**
- **Database migrations** with Goose

### Configuration

The Auth Service can be configured using environment variables or configuration files located in the /config directory.
The main configuration file is config/values.yaml.

Before running the service, create these files in the ./config directory:
- values.yaml — see values_example.yaml for reference.

#### Example `./config/values.yaml` :

```
# PostgreSQL database
database:
  name: ""                  # Database name
  user: ""                  # Database user
  password: ""              # Database password
  host: ""                  # Database host
  port: ""                  # Database port
  ssl_mode: ""              # SSL mode (disable, require, etc.)
  max_cons:                 # Maximum connections in pool
  min_cons:                 # Minimum connections in pool
  max_conn_lifetime: ""     # Maximum connection lifetime

# JWT Token
jwt:
  token_secret: ""          # Secret key for signing tokens
  access_token_ttl: ""      # Access token TTL (e.g., 15m)
  refresh_token_ttl: ""     # Refresh token TTL (e.g., 7d)

# Redis
redis:
  addr: ""                  # Redis host:port
  password: ""              # Redis password
  bd: 0                     # Redis database number
  key: ""                   # Key prefix for caching

# Worker
refresh_resources_ttl: ""       # TTL for cached resources
refresh_resources_interval: ""  # Interval for resource cache refresh

# Swagger UI
swagger:
  url: ""                       # Base URL for Swagger UI
  auth_spec_url: ""              # URL for auth API spec
  auth_spec_path: ""             # Local path for auth API spec
  auth_url: ""                   # Auth endpoint URL
  resource_spec_url: ""          # URL for resource API spec
  resource_spec_path: ""         # Local path for resource API spec
  resource_url: ""               # Resource endpoint URL

# HTTP server
server:
  http:
    port: ""                     # Port for HTTP server
    proxy_prefix: ""             # Proxy URL prefix
    auth_prefix: ""              # Auth endpoint prefix
    resource_prefix: ""          # Resource endpoint prefix

# gRPC client
grpc_client:
  max_retry:                     # Maximum number of retries
  timeout_per_retry: ""          # Timeout per retry

user_service:
  host: ""                       # User Service host URL

# Worker modes
worker_mode:
  refresh_routes: ""             # Enable/disable automatic route refresh
```

## Database Schema

Main entities used by Auth Service:

- **Users**: system users managed via User Service.
- **Services**: registered backend services for proxying.
- **Resources**: endpoints/resources associated with services.
- **Roles**: user roles and permissions for access control.

Migrations are located in `/migrations` and handled automatically on startup.

## Getting Started

### Prerequisites
- Go 1.24.2+
- Docker & Docker Compose

### Installation

1. **Clone the repository**
   ```bash
   git clone https://github.com/handmade-jewelry/auth-service.git
   cd auth-service
   ```

2. **Configure the service using environment variables or the config files in `/config`**

3. Generate Go code from OpenAPI specifications
    ```bash
    make generate-resource
    make generate-auth
    ```
    Or use the corresponding Makefile targets if defined.

4. **Start the services**
   ```bash
   docker-compose up -d
   ```

## Usage

After installing and starting the Auth Service, you can use it to handle authentication, authorization, and request proxying for all registered services.

### Login
`POST /login`
- Forwarded to User Service via gRPC. 
- Returns access token and refresh token in cookies. 
- Access token is used for subsequent API requests

### Logout
`POST /logout`
- Invalidates access and refresh tokens to terminate the session.

### Refresh Token
`POST /refresh-token`
- Uses the refresh token to generate a new access token.
- Ensures continuous authentication without requiring login.

### Accessing Protected Services
Proxy Endpoint
`{service}/{resource}`
- Auth Service validates:
    - Service and resource existence in the database. 
    - Access token validity. 
    - User roles against required permissions.
- If all checks pass, the request is forwarded to the target service.

### Managing Services and Resources (CRUD)
Manage registered services.
Service Endpoints:
- `GET /services`
- `POST /services`
- `PUT /services/{id}`
- `DELETE /services/{id}`

### Resource Endpoints
Manage endpoints/resources associated with services.
Resources Endpoints:
- `GET /resources`
- `POST /resources`
- `PUT /resources/{id}`
- `DELETE /resources/{id}`


## Proxy Logic

### Service & Resource Validation
- Checks if the requested {service} exists in the database. 
- Checks if the requested {resource} (endpoint) is registered for that service. 
- Rejects the request with 404 Not Found if either the service or resource is missing.

### Access Token Verification
- Extracts the access token from cookies. 
- Validates the token's signature and expiration. 
- If the access token is missing or expired, the request can optionally use a refresh token to generate a new access token.

### Role-Based Authorization
- Checks user roles embedded in the access token against roles required for the requested resource (stored in the database).
- If the user lacks the required roles, the request is rejected with 403 Forbidden.

### Request Forwarding
- After passing all checks, the request is forwarded to the appropriate backend service.
- Auth Service preserves headers and request body, acting transparently as a gateway.

### Error Handling
- Returns 401 Unauthorized if token validation fails.
- Returns 403 Forbidden if role validation fails.
- Returns 404 Not Found if service/resource is not registered.

### Caching
- Services and resources are cached in Redis for faster validation.
- Cache is refreshed periodically according to the worker configuration (refresh_resources_interval).

### Example Flow

- User sends a request to GET /bank/getAccounts.
- Auth Service:
  - Checks that bank service and getAccounts resource exist.
  - Validates the access token from the cookie.
  - Verifies that the user has the bank_read role. 
  - Forwards the request to the Bank Service if authorized. 
- The Bank Service processes the request and returns a response. 
- Auth Service returns the response to the client.

> Note: `Bank Service` is used here as an illustrative example. Auth Service can proxy any registered service.


