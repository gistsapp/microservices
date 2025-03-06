# Auth Microservice

This microservice is responsible for authenticating users and managing their identities.

## Structure

### Config

The config package contains the configuration for the microservice using the viper library.

### Core

The core package contains the business logic for the microservice, including the authentication logic and the identity management.

### HTTP

The http package contains the HTTP handlers for the microservice, including the authentication endpoints. So the business logic is agnostic from the deliverer.

## Configuration

The configuration is loaded from a JSON file using the viper library.

```json
{
    "port": "string",
    "keycloak": {
        "realm": "string",
        "url": "string",
        "client_id": "string",
        "client_secret": "string"
    },
    "database": {
        "host": "string",
        "port": "string",
        "user": "string",
        "password": "string",
        "name": "string"
    }
}
```


