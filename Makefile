run:
	@npx tailwindcss -i ./css/root.css -o ./static/css/styles.css --minify
	@templ generate
	@go build -o ./tmp/main cmd/main.go
