APP := tax-calculator
HTTPPORT := 8080
ENV ?= local

# Run the app
run:
	go run main.go

# Run tests
test:
	go test ./...

# Format code
fmt:
	go fmt ./...

# Lint (if using golangci-lint or similar)
lint:
	golangci-lint run

# Clean build artifacts
clean:
	rm -f $(APP)

# Build the binary
build-binary:
	CGO_ENABLED=0 go build .

# Pre-commit workflow
pre-commit: 
	$(MAKE) test 
	$(MAKE) docker-build 
	$(MAKE) docker-launch 
	$(MAKE) docker-smoke-test

# Docker utilities
docker-clean:
	-docker rm -f $(APP)
	-docker network rm net-$(APP)
	-docker rm -f smoke-$(APP)

docker-network:
	-docker network create net-$(APP)

docker-build:
	docker build -t $(APP):$(ENV) --build-arg ENV=$(ENV) .

docker-launch:
	-docker network create net-$(APP)
	docker run -d --name $(APP) \
		--network net-$(APP) \
		-p ${HTTPPORT}:${HTTPPORT} \
		$(APP):$(ENV)
		
docker-smoke-test:
	# Remove previous smoke test container
	-docker rm -f smoke-$(APP)
	# Timeout to wait for the container to initialize (adjust as needed)
	sleep 5
	# Run the smoke test
	docker run -v $$PWD/.tests/smoke_test.yml:/etc/smoke/conf.d/smoke_test.yml \
		--network net-$(APP) \
		--name smoke-$(APP) \
		bluehoodie/smoke \
		-f /etc/smoke/conf.d/smoke_test.yml \
		-u http://$(APP):${HTTPPORT} -v
