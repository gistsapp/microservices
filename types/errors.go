package types

import "errors"

// ErrNotFound is returned when an entity is not found, should be mapped to 404 status code or equivalent status in other protocol than HTTP
var ErrNotFound error = errors.New("Entity not found")
