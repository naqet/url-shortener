run:
	@npx tailwindcss -o ./static/css/styles.css --minify
	@templ generate
	@go build -o ./tmp/main cmd/main.go
