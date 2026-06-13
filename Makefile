.PHONY: migrate test ask ingest seed run db-create
# Reuses the shared pgvector server on :5433 (container llm-slm-postgres-1).
db-create:    ; @docker exec llm-slm-postgres-1 psql -U solar -d solar -tc "SELECT 1 FROM pg_database WHERE datname='french'" | grep -q 1 || docker exec llm-slm-postgres-1 createdb -U solar french
migrate:      ; go run ./cmd/server -migrate
seed:         ; python3 data/generate_seed.py
ingest:       ; go run ./cmd/ingest
run:          ; go run ./cmd/server
ask:          ; go run ./cmd/ask
test:         ; go test ./...
