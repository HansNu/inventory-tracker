# Inventory Tracker API
![CI](https://github.com/HansNu/inventory-tracker/actions/workflows/ci.yml/badge.svg)

A REST API for tracking IT assets — hardware inventory, assignment to users, categories and category groups — built with Go, Gin and PostgreSQL.

Live: /swagger/index.html · /health Frontend: HansNu/inventory-tracker-fe

The Swagger page is interactive: call /auth/login, copy the token, paste it into Authorize as Bearer <token>, and every protected endpoint becomes clickable.

**Stack**

Go 1.26.2 · Gin · PostgreSQL (pgx/v5) · JWT (golang-jwt/v5) · bcrypt · Docker · GitHub Actions · swaggo

**Architecture**
handler/     HTTP only — parse the request, map errors to status codes
   ↓         (service.AssetService)
services/    business logic — validation, query construction, pagination
   ↓         (repository.AssetRepository)
repository/  database only — execute queries, scan rows

Each layer is injected with an interface rather than a concrete type, so each one can be tested against a fake of the layer beneath it: handler tests use a fake service, repository tests use pgxmock in place of a real connection pool. main.go is the only file that knows the concrete types.

middleware/ holds authentication (AuthMiddleware) and authorisation (RoleMiddleware) as separate, stackable concerns — routes can require a valid token, or a valid token and a specific role.

Running locally
With Docker
bash
docker compose up

Starts Postgres (schema loaded from schema.sql on first run) and the API on :8080. The API container waits for the database's healthcheck rather than just for the container to exist.

Without Docker
bash
export DATABASE_URL="postgres://user:pass@localhost:5432/inventory"
export JWT_SECRET="any-long-random-string"
export FRONTEND_ORIGIN="http://localhost:5173"   # optional, this is the default

go run .        # listens on $PORT, or :8080

JWT_SECRET is read at package initialisation and the app panics without it — deliberately, so a missing secret fails loudly at startup instead of silently signing every token with an empty key.

**Health check**
GET /health pings the database and returns 503 if it can't be reached. It deliberately checks the dependency rather than returning 200 for a live process: "running" and "working" are different claims.

**Authentication**

POST /api/auth/register and POST /api/auth/login return a JWT valid for 24 hours. Every other endpoint except /health and /api/getCategoryGroup requires it.

bash
curl -X POST https://inventory-tracker-production-a38f.up.railway.app/api/auth/login \
  -H 'Content-Type: application/json' \
  -d '{"username":"admin","password":"secret"}'
# => {"token":"eyJhbGc...","user":{...}}

curl "https://inventory-tracker-production-a38f.up.railway.app/api/getAssetList?page=1&pageSize=10" \
  -H "Authorization: Bearer eyJhbGc..."

The token is signed with HS256 and the parser is pinned to that algorithm via jwt.WithValidMethods — without it, a token carrying "alg": "none" would be accepted as valid. There is a test for exactly that.

**Endpoints**
Method	Path	Auth
GET	/health	—
POST	/api/auth/register	—
POST	/api/auth/login	—
GET	/api/getCategoryGroup	—
GET	/api/me	token
GET	/api/getAssetList	token
GET	/api/getAssetByAssetCode/{assetCode}	token
POST	/api/addAsset	token
PUT	/api/updateAsset/{id}	token
DELETE	/api/deleteAssetByAssetCode/{assetCode}	token
GET	/api/getAssetCategoryList	token
POST	/api/addAssetCategory	token
DELETE	/api/deleteAssetCategoryById/{id}	token
POST	/api/addCategoryGroup	token
DELETE	/api/deleteCategoryGroup	token + IT role

getAssetList supports page, pageSize, search (across code, name, category, brand, serial, status, location and user), status, type, sortField and sortOrder.

Endpoint names are verb-style by choice rather than resource-style, so every route has a unique path signature in the logs regardless of method.

**Tests**
bash
JWT_SECRET=test-secret go test ./... -race

Table-driven handler tests for AddAsset and GetAssetList against a fake service, and middleware tests covering missing, malformed, expired, wrong-key and alg: none tokens — plus a separate test that a valid token puts the right values into the request context.

JWT_SECRET has to be set before the test binary starts: the secret is read into a package-level variable, which initialises before TestMain.

CI

.github/workflows/ci.yml runs on every push: go vet, go test -race, go build, and a separate job that builds the Docker image. The Go version comes from go.mod so it can't drift.

Deployment

Deployed on Railway from this repo's Dockerfile, with Postgres on Supabase.

Two things that are easy to get wrong:

Supabase's direct connection is IPv6-only. Use the Supavisor session pooler (port 5432 on aws-N-<region>.pooler.supabase.com, username postgres.<project-ref>) — IPv4-reachable, and session mode keeps pgx's prepared-statement cache working. Transaction mode on 6543 breaks it.
The container must bind $PORT, not a hardcoded 8080, or the platform's health check never reaches it.
Known gaps
