# Inventory Tracker API

A REST API for tracking IT assets — hardware inventory, assignment to users,
categories and category groups — built with Go, Gin and PostgreSQL.

## Stack

Go 1.26.2 · Gin · PostgreSQL (pgx/v5) · JWT auth (golang-jwt/v5) · bcrypt

## Architecture

The asset endpoints are organised in three layers, each depending only on the
one below it through an interface:

```
handler/     HTTP only — parse the request, map errors to status codes
   ↓         (service.AssetService)
services/    business logic — validation, query construction, pagination
   ↓         (repository.AssetRepository)
repository/  database only — execute queries, scan rows
```

Each layer is injected with an interface rather than a concrete type, so each
one can be tested against a fake of the layer beneath it: handler tests use a
fake service, repository tests use `pgxmock` in place of a real connection pool.

`middleware/` holds authentication (`AuthMiddleware`) and authorisation
(`RoleMiddleware`) as separate, stackable concerns — routes can require a valid
token, or a valid token *and* a specific role.

## Running locally

Requires PostgreSQL and two environment variables:

```bash
export DATABASE_URL="postgres://user:pass@localhost:5432/inventory"
export JWT_SECRET="any-long-random-string"

go run main.go        # listens on :8080
```

`JWT_SECRET` is read at package initialisation and the app panics without it —
deliberately, so a missing secret fails loudly at startup instead of silently
signing tokens with an empty key.

## Authentication

`POST /api/auth/register` and `POST /api/auth/login` return a JWT valid for 24
hours. Every other endpoint requires it:

```bash
curl -X POST localhost:8080/api/auth/login \
  -H 'Content-Type: application/json' \
  -d '{"username":"admin","password":"secret"}'
# => {"token":"eyJhbGc...","user":{...}}

curl localhost:8080/api/getAssetList?page=1&pageSize=10 \
  -H "Authorization: Bearer eyJhbGc..."
```

## Endpoints

| Method | Path | Auth |
|---|---|---|
| POST | `/api/auth/register` | — |
| POST | `/api/auth/login` | — |
| GET | `/api/me` | token |
| GET | `/api/getAssetList` | token |
| GET | `/api/getAssetByAssetCode/:assetCode` | token |
| POST | `/api/addAsset` | token |
| PUT | `/api/updateAsset` | token |
| DELETE | `/api/deleteAssetByAssetCode` | token |
| GET | `/api/getAssetCategoryList` | token |
| POST | `/api/addAssetCategory` | token |
| DELETE | `/api/deleteAssetCategoryById` | token |
| GET | `/api/getCategoryGroup` | token |
| POST | `/api/addCategoryGroup` | token |
| DELETE | `/api/deleteCategoryGroup` | token + `IT` role |

`getAssetList` supports `page`, `pageSize`, `search` (across code, name,
category, brand, serial, status, location, user), `status`, `type`, `sortField`
and `sortOrder`.

## Tests

```bash
JWT_SECRET=test-secret go test ./... -v
```

Covered: handler-level tests for `AddAsset` and `GetAssetList` against a fake
service, and middleware tests covering missing, malformed, expired,
wrong-key and `alg: none` tokens.

## Known gaps

- Endpoint naming is RPC-style (`/addAsset`) rather than resource-oriented
  (`POST /assets`); a REST-conventional rewrite is planned.
- `auth.go` and the category handlers still query the database directly and
  have not been moved into the service/repository layers.
- No container or CI pipeline yet — both in progress.
- Tokens are stateless and cannot be revoked before expiry.
