newdist:
	rm -rf dist
	mkdir -p dist

wasm:
	GOOS=js GOARCH=wasm go build -o dist/main.wasm ./cmd/wasm
	cp web/static/wasm_exec.js dist/

static:
	npx @tailwindcss/cli -i ./web/static/input.css -o ./dist/main.css
	cp web/static/index.html dist/
	cp web/static/main.js dist/

build: newdist wasm static
	@echo "Build complete - populated into dist/"
 