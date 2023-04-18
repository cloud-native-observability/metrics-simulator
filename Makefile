

build: pkg/metrics/metrics.go pkg/config/config.go cmd/main.go
	mkdir -pv build/
	go build -o build/metrics-simulator ./cmd/
