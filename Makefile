tw:
	./tailwindcss -i ./static/css/input.css -o ./static/css/style.css --watch

dev: 
	air

templ-proxy:
	templ generate -watch -proxy=http://localhost:4000

.PHONY: tailwind-build
tailwind-build:
	./tailwindcss -i ./static/css/input.css -o ./static/css/style.css --minify

.PHONY: templ-generate
templ-generate:
	templ generate

.PHONY: build
build:
	make tailwind-build && make templ-generate && go build -o ./bin/$(APP_NAME) ./$(APP_NAME)
