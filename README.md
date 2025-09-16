# Odoo Signup System

A secure web-based signup system for creating or cloning Odoo databases with a modern frontend interface.

## Features

- Modern, responsive UI built with HTML/CSS/JavaScript
- Secure Go backend using Gin framework
- Odoo integration for database creation/cloning and user setup
- Form validation (client/server-side)
- Rate limiting to prevent abuse
- Configurable via environment variables
- Docker support for easy deployment

## Prerequisites

- Go 1.21 or later
- Odoo instance with multi-tenant support (for database operations)
- Docker (optional, for containerized deployment)

## Installation

1. Clone the repository:
   ```
   git clone <repository-url>
   cd odoo-signup
   ```

2. Install dependencies:
   ```
   go mod download
   ```

3. Copy and configure environment:
   ```
   cp .env.example .env
   # Edit .env with your values
   ```

## Configuration

Edit `.env` with the following variables:

```
# Server Configuration
PORT=8080
ENVIRONMENT=development

# Domain Configuration
DOMAIN=yourdomain.com
ODOO_COMPANY=YourCompany

# Odoo Configuration
ODOO_URL=http://localhost:8069
ODOO_MASTER_PASSWORD=your_master_password
TEMPLATE_DATABASE=odoo-template
ADMIN_USER=admin
ADMIN_PASSWORD=your_admin_password
DEFAULT_DB_MODE=create  # or "clone"

# Rate Limiting (requests per second)
RATE_LIMIT=10
BURST_LIMIT=20

# Timeout
HTTP_TIMEOUT_SECONDS=600

# Logging
LOG_LEVEL=info
```

**Notes:**
- `ODOO_MASTER_PASSWORD` is required for database operations.
- For clone mode, ensure `TEMPLATE_DATABASE` exists in Odoo.
- Set `DOMAIN` for generating instance URLs (e.g., username.yourdomain.com).

## Running the Application

### Development
```
go run ./cmd/server/main.go
```

The server starts on `http://localhost:8080`.

### Production
```
go build -o odoo-signup ./cmd/server
./odoo-signup
```

## Usage

1. Access the signup form at `http://localhost:8080`.
2. Fill in the form: username (becomes subdomain), email, password, personal/company details, country, accept terms.
3. Submit: The system creates/clones an Odoo database, sets up the admin user, and returns the instance URL.

## API Endpoints

### POST `/api/signup`
Creates/clones Odoo database and user. Supports `?db_mode=clone` query param.

**Request Body:**
```json
{
  "username": "mycompany",
  "email": "admin@mycompany.com",
  "password": "SecurePass123",
  "firstName": "John",
  "lastName": "Doe",
  "phone": "+1234567890",
  "companyName": "My Company Ltd",
  "industry": "technology",
  "companySize": "11-50",
  "country": {
    "id": 234,
    "code": "US",
    "name": "United States"
  },
  "terms": true
}
```

**Response (Success):**
```json
{
  "success": true,
  "message": "Signup successful using create mode! Your Odoo instance is ready.",
  "data": {
    "instanceUrl": "mycompany.yourdomain.com",
    "email": "admin@mycompany.com",
    "database": "mycompany"
  }
}
```

### GET `/api/health`
Health check: Returns `{"status": "healthy", "timestamp": "..."}`.

## Deployment

### Docker
1. Build:
   ```
   docker build -t odoo-signup .
   ```

2. Run:
   ```
   docker run -p 8080:8080 --env-file .env odoo-signup
   ```

Configure your web server (e.g., Nginx) to proxy requests and route subdomains to Odoo. Use HTTPS in production.

## License

MIT License - see [LICENSE](LICENSE) for details.