dev-recreate:
	@docker-compose --project-name=go-go --env-file deploy/dev/.env -f deploy/dev/docker-compose.yaml up -d --build --force-recreate

build-and-push-image: build-image push-image

build-image:
	@docker build -f release.Dockerfile . -t aqaurius6666/mainserver:pre-release

push-image:
	@docker tag supermedicalchain/main-service:pre-release supermedicalchain/main-service${TAG}
	@docker push supermedicalchain/main-service${TAG}

build:
	@go build -o ./dist/server ./src/internal

debug:
	@/bin/sh -c "dlv --listen=0.0.0.0:2345 --headless=true --api-version=2 exec ./dist/server -- --disable-profiler --allow-kill serve"

serve:
	@./dist/server serve

dev:
	@./dist/server --log-format plain --log-level debug --disable-profiler --allow-kill serve

cleanDB:
	@./dist/server clean

seed:
	@./dist/server seed --clean
	@echo Hello
docker-seed:
	@docker-compose --project-name=go-go --env-file deploy/dev/.env -f deploy/dev/docker-compose.yaml exec apiservice make seed  

dlv:
	@./scripts/dlv.sh
	
kill:
	@(echo '{}' | grpc-client-cli -service Common -method Kill localhost:${CONFIG_GRPC_PORT}) > /nil 2> /nil || return 0
	@echo "Enemy has been slain"

logs:
	@docker-compose --project-name=go-go -f deploy/dev/docker-compose.yaml logs -f apiservice

proto:
	@./scripts/gen-proto.sh
	@$(MAKE) fix-import
	@$(MAKE) swagger

monitor:
	@docker-compose --project-name=go-go -f deploy/dev/docker-monitor.yaml up -d --build --force-recreate

lint:
	@golangci-lint run

wire:
	@./scripts/wire.sh

swagger:
	@go generate src/services/swagger github.com/aqaurius6666/apiservice/src/services/swagger


sql:
	@go generate src/internal/db/seed github.com/aqaurius6666/apiservice/src/internal/db/seed

# "github.com/aqaurius6666/apiservice/src/internal/var/c"

fix-import:
	@find . -type f -name "*.pb.*" -print0 | xargs -0 sed -i 's/\.\/c/github.com\/aqaurius6666\/apiservice\/src\/internal\/var\/c/' 
fix-omit-empty:
	@find . -type f -name "*.pb.*.go" -print0 | xargs -0 sed -i 's/omitempty//'
ensure:
	@./scripts/ensure.sh

test:
	@go test github.com/aqaurius6666/apiservice/src/internal/api -count=1