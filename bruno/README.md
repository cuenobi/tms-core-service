# Bruno API Collection

This directory contains Bruno API requests for testing the TMS Core Service.

## What is Bruno?

Bruno is an open-source API client that stores collections in your filesystem as plain text files (`.bru` format). Unlike Postman, Bruno doesn't require an account and all collections are stored locally in your project.

## Installation

Install Bruno from: https://www.usebruno.com/downloads

Or via Homebrew:
```bash
brew install bruno
```

## Usage

1. Open Bruno
2. Click "Open Collection"
3. Select this `bruno` directory
4. You'll see all the requests organized in folders

## Collection Structure

```
bruno/
├── bruno.json              # Collection metadata
├── environments/
│   └── Local.bru          # Local environment variables
├── Health Check.bru       # Health check endpoint
├── Auth/
│   ├── Register User.bru  # User registration
│   └── Login.bru          # User login
└── Protected/
    └── Example.bru        # Protected endpoint template
```

## Quick Start

1. **Start the services**:
   ```bash
   docker-compose up -d
   go run main.go migrate up
   go run main.go serve
   ```

2. **Test the endpoints in order**:
   - Health Check → Verify service is running
   - Register User → Create a new account (token auto-saved)
   - Login → Authenticate (token auto-saved)
   - Protected endpoints → Use saved token automatically

## Environment Variables

- `base_url`: API base URL (default: http://localhost:8080)
- `access_token`: JWT access token (auto-populated after login/register)

## Auto Token Management

The Login and Register requests automatically extract and save the `access_token` to the environment variables. Protected endpoints will use this token automatically.

## Adding New Requests

Create new `.bru` files following the existing patterns:

```
meta {
  name: Your Request Name
  type: http
  seq: 5
}

get {
  url: {{base_url}}/your/endpoint
  body: none
  auth: bearer
}

auth:bearer {
  token: {{access_token}}
}
```

## Tips

- Use `{{variable}}` syntax for environment variables
- Scripts in `script:post-response` blocks run after each request
- Organize related requests in folders
- All changes are tracked in git automatically
