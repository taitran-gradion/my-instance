# A Margince installation that requires core's Go modules instead of vendoring core.
#
# Go: go.mod requires github.com/margince/margince/backend, and this module
# builds core's cmd/api, cmd/worker and cmd/migrate itself, linked with the
# composition it generated (replace composition => ./build/composition/backend).
#
# Core is the fork that carries the extensions.mod POC (margince/margince#7727),
# at the version go.mod replaces core's modules with. Stand-ins, until core publishes what #7727 proposes:
#   - `go tool margince`    → a clone of the fork at that tag in .margince/core,
#                             whose make composition and Vite this file runs
#   - @margince/app (SPA)   → that clone's frontend/
#
# Infrastructure is the local dev Postgres, Redis and object store that
# `make dev` in core starts; this instance takes its own database, Redis
# database number, bucket and ports.

CORE      := .margince/core
# The core version is the one go.mod replaces the backend with, so composing and
# building can never use two different cores.
CORE_VERSION := $(shell awk '$$1 == "github.com/margince/margince/backend" && $$2 == "=>" {print $$NF}' go.mod)
CORE_URL  ?= https://github.com/taitran-gradion/margince
API_PORT  := 18190
WEB_PORT  := 8190
DB        := margince_instance
REDIS     := localhost:16379/50
BUCKET    := margince-instance
OWNER_DSN := postgres://margince_owner:dev@localhost:15432/$(DB)
APP_DSN   := postgres://margince_app:margince_app_dev@localhost:15432/$(DB)
PG        := docker exec -i margince-postgres-1 psql -q -U margince_owner
RUN       := .run
# What extensions compile over core (backend/overlays/), written by compose for
# `go build -overlay`. Go refuses to overlay files in the module cache, so with
# overlays the binaries are built from the core clone, at the same version.
OVERLAY   := $(abspath $(CORE))/build/overlay/overlay.json
REVISION  := instance-$(CORE_VERSION)

export GOWORK   := off
export GOFLAGS  := -mod=mod

export MARGINCE_ENV                  := dev
export MARGINCE_KEYVAULT_ROOT_KEY    := bWFyZ2luY2UtZGV2LW9ubHkta2V5dmF1bHQtcm9vdGs=
export MARGINCE_CONNECTOR_STATE_KEY  := margince-dev-connector-state-key-0001
export MARGINCE_BLOBSTORE_ENDPOINT   := localhost:29000
export MARGINCE_BLOBSTORE_ACCESS_KEY := minioadmin
export MARGINCE_BLOBSTORE_SECRET_KEY := minioadmin
export MARGINCE_BLOBSTORE_REGION     := us-east-1
export MARGINCE_BLOBSTORE_BUCKET     := $(BUCKET)

.PHONY: help core install up compose build db migrate run web stop logs clean

help: ## List the targets
	@grep -E '^[a-z-]+:.*## ' $(MAKEFILE_LIST) | awk -F':.*## ' '{printf "  %-8s %s\n", $$1, $$2}'

core: ## Fetch core at the version go.mod names (again only when it changes)
	@test -n "$(CORE_VERSION)" || { echo "go.mod replaces no github.com/margince/margince/backend version"; exit 1; }
	@if [ "$$(cat $(CORE)/.version 2>/dev/null)" != "$(CORE_VERSION)" ]; then \
	  rm -rf $(CORE) && git init -q $(CORE) \
	  && git -C $(CORE) fetch -q --depth 1 $(CORE_URL) refs/tags/backend/$(CORE_VERSION) \
	  && git -C $(CORE) -c advice.detachedHead=false checkout -q FETCH_HEAD \
	  && echo $(CORE_VERSION) > $(CORE)/.version; fi

install: core ## Once: core's frontend dependencies, for the web app
	@$(MAKE) -s -C $(CORE) fe-install GOWORK= GOFLAGS=

up: install compose build db migrate run web ## Everything, in order
	@echo; echo "  OPEN  http://localhost:$(WEB_PORT)   (api http://localhost:$(API_PORT))"
	@echo "  login admin@demo.test / $$(cat config/margince-admin-password)"; echo "  stop  make stop"

compose: core ## Generate this instance's composition from extensions.mod
	@$(MAKE) -s -C $(CORE) composition EXTENSIONS_MOD=$(abspath extensions.mod) GOWORK= GOFLAGS=
	@echo "composed: $$(grep -o 'mustBe("[^"]*"' $(CORE)/build/composition/backend/extensions_gen.go | cut -d'"' -f2 | paste -sd ' ' -)"

build: ## Build core's role binaries from the required module (from the core clone when extensions overlay it)
	@if [ -f $(OVERLAY) ]; then \
	  echo "extensions overlay core: building from $(CORE)"; \
	  for role in api worker migrate; do \
	    (cd $(CORE)/backend && GOWORK=$(abspath $(CORE))/build/composition/go.work GOFLAGS= go build -overlay $(OVERLAY) \
	      -ldflags "-X github.com/margince/margince/backend/internal/shared/buildinfo.Revision=$(REVISION)" \
	      -o $(abspath bin)/$$role ./cmd/$$role) || exit 1; done; \
	else \
	  for role in api worker migrate; do \
	    go build -ldflags "-X github.com/margince/margince/backend/internal/shared/buildinfo.Revision=$(REVISION)" \
	      -o bin/$$role github.com/margince/margince/backend/cmd/$$role || exit 1; done; \
	fi
	@go version -m bin/api | grep -E 'composition|my-extension' | sed 's/^/  /'

db: ## Create this instance's database once
	@echo "SELECT 'CREATE DATABASE $(DB)' WHERE NOT EXISTS (SELECT FROM pg_database WHERE datname = '$(DB)')\gexec" | $(PG) -d postgres
	@$(PG) -d $(DB) -v ON_ERROR_STOP=1 < $(CORE)/scripts/db-init.sql

migrate: ## Apply core and extension migrations
	@MARGINCE_OWNER_DSN=$(OWNER_DSN) ./bin/migrate up

run: ## Start the api and the worker
	@for p in api worker; do [ ! -f $(RUN)/$$p.pid ] || { kill $$(cat $(RUN)/$$p.pid) 2>/dev/null || true; rm -f $(RUN)/$$p.pid; }; done
	@mkdir -p $(RUN)
	@[ -f config/margince-admin-password ] || { printf 'instance-first-password' > config/margince-admin-password; chmod 600 config/margince-admin-password; }
	@MARGINCE_SCHEMA_DSN=$(OWNER_DSN) nohup ./bin/api --addr :$(API_PORT) --dsn $(APP_DSN) --config config/margince.yaml \
	  --redis $(REDIS) --public-base-url http://localhost:$(WEB_PORT) --ai-fake > $(RUN)/api.log 2>&1 & echo $$! > $(RUN)/api.pid
	@for i in $$(seq 60); do curl -sf localhost:$(API_PORT)/readyz >/dev/null && break; sleep 1; done; \
	  curl -sf localhost:$(API_PORT)/readyz >/dev/null || { echo "api not ready — see $(RUN)/api.log"; exit 1; }
	@nohup ./bin/worker --dsn $(APP_DSN) --redis $(REDIS) --config config/margince.yaml \
	  --public-base-url http://localhost:$(WEB_PORT) --ai-fake > $(RUN)/worker.log 2>&1 & echo $$! > $(RUN)/worker.pid
	@echo "api :$(API_PORT) and worker running"

web: ## Serve the SPA with this instance's screens (core's Vite, a stand-in)
	@[ ! -f $(RUN)/web.pid ] || { kill $$(cat $(RUN)/web.pid) 2>/dev/null || true; rm -f $(RUN)/web.pid; }
	@mkdir -p $(RUN)
	@BACKEND_PORT=$(API_PORT) MARGINCE_BUILD_REVISION=$(REVISION) \
	  MARGINCE_COMPOSITION_FRONTEND=$(abspath $(CORE))/build/composition/frontend \
	  nohup pnpm --dir $(CORE)/frontend exec vite --port $(WEB_PORT) --strictPort > $(RUN)/web.log 2>&1 & echo $$! > $(RUN)/web.pid
	@for i in $$(seq 30); do curl -sf -o /dev/null localhost:$(WEB_PORT) && break; sleep 1; done; echo "web :$(WEB_PORT) running"

stop: ## Stop everything this instance started
	@for p in api worker web; do [ ! -f $(RUN)/$$p.pid ] || { kill $$(cat $(RUN)/$$p.pid) 2>/dev/null || true; rm -f $(RUN)/$$p.pid; }; done

logs: ## Tail the api and worker logs
	@tail -n 40 $(RUN)/api.log $(RUN)/worker.log

clean: stop ## Remove build output and binaries (keeps the database)
	rm -rf bin $(RUN)
	rm -rf .margince
