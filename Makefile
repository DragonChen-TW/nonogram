default: run
run:
	go run main.go
build:
	GOOS=js GOARCH=wasm go build -o build_wasm/main.wasm main.go