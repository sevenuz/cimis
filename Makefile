test:
	go test ./... -v --cover

test-report:
	go test ./... -v --cover -coverprofile=coverage.out
	go tool cover -func=coverage.out
	go tool cover -html=coverage.out

todo:
# "|| true" is necesarry because grep returns 1 on empty results
	@grep -rn TODO ui/src/ || true
	@grep -rn TODO routes/ || true
	@grep -rn TODO hooks/ || true
	@grep -rn TODO models/ || true
	@grep -rn TODO migrations/ || true

build-ui:
	cd ui; npm install; npm run build

build: build-ui
	go build

dev:
	@echo "This command does not run the backend"
	cd ui; npm run dev;

run:
	@echo "This command does not build the frontend!"
	go run main.go serve

build-windows: build-ui
	@echo "Build is for Arch amd64."
	GOOS=windows GOARCH=amd64 go build

build-darwin: build-ui
	@echo "Build is for Arch amd64."
	GOOS=darwin GOARCH=amd64 go build

build-linux: build-ui
	@echo "Build is for Arch amd64."
	GOOS=linux GOARCH=amd64 go build
