package graph

import "rider-service/internal/core/ports"

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	RiderService ports.RiderService
}
