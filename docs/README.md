# TMS Core Service - Documentation

## Overview

This directory contains documentation and generated API documentation.

### Swagger Documentation

The Swagger/OpenAPI documentation is auto-generated from code annotations using [swaggo](https://github.com/swaggo/swag).

To generate documentation:

```bash
make swagger
# or
swag init -g main.go -o ./docs
```

The generated files will be placed in this directory:
- `docs.go` - Go bindings
- `swagger.json` - OpenAPI JSON spec
- `swagger.yaml` - OpenAPI YAML spec

### Viewing Documentation

Once the server is running, visit:
```
http://localhost:8080/swagger/index.html
```

### Adding API Documentation

Add Swagger annotations to your handler functions. Example:

```go
// Login godoc
// @Summary Login
// @Description Authenticate user and get tokens
// @Tags auth
// @Accept json
// @Produce json
// @Param request body auth.LoginRequest true "Login credentials"
// @Success 200 {object} httpresponse.Response{data=auth.AuthResponse}
// @Failure 401 {object} httpresponse.Response
// @Router /auth/login [post]
func (h *Handler) Login(c *gin.Context) {
    // ...
}
```
