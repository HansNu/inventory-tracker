package handler

import (
	"inventory-tracker/repository"
	service "inventory-tracker/services"
)

// Handler holds every boundary this package is allowed to call.
//
// Assets and auth go through services, because both carry rules worth keeping
// out of HTTP code — pagination and sort normalisation for one, uniqueness,
// hashing and token issuance for the other.
//
// Categories and category groups go straight to their repositories. They are
// pure CRUD with no rules to protect, and a service that only forwards a call is
// an extra hop pretending to be a boundary. The moment a rule appears there, a
// service is where it goes.
//
// There is deliberately no *pgxpool.Pool here. No handler in this package can
// reach a database except through one of these four fields.
type Handler struct {
	Service        service.AssetService
	Auth           service.AuthService
	Categories     repository.CategoryRepository
	CategoryGroups repository.CategoryGroupRepository
}
