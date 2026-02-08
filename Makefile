newdist:
	rm -rf dist
	mkdir -p dist/assets

data:
	go generate ./...

wasm:
	GOOS=js GOARCH=wasm go build -o dist/main.wasm ./cmd/wasm

styling:
	npx @tailwindcss/cli -i ./web/static/input.css -o ./dist/main.css

static:
	cp web/static/index.html dist/
	cp web/static/main.js dist/
	cp web/static/wasm_exec.js dist/
	cp web/static/assets/favicon.ico dist/assets/
	cp web/static/assets/help-icon.svg dist/assets/

build: newdist data wasm styling static
	@echo "Build complete - populated into dist/"
 