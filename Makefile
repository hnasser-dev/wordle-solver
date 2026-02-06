newdist:
	rm -rf dist
	mkdir -p dist

data:
	go generate ./...

wasm:
	GOOS=js GOARCH=wasm go build -o dist/main.wasm ./cmd/wasm
	cp web/static/wasm_exec.js dist/

static:
	npx @tailwindcss/cli -i ./web/static/input.css -o ./dist/main.css
	cp web/static/index.html dist/
	cp web/static/main.js dist/
	cp web/static/favicon.ico dist/

build: newdist data wasm static
	@echo "Build complete - populated into dist/"
 