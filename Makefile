dotenv:
	@echo "Generating .env file..."
	@cp .env.empty cmd/.env

build:
	@echo "Building the project..."
	@CGO_ENABLED=0 go build -o bin/cli cmd/main.go
	@echo "Moving .env to bin..."
	@cp cmd/.env bin/.env
	@echo ".env copy to bin directory."
	@echo "Build complete."