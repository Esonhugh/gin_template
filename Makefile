

.PHONY: app
app:
	go build -o app .

run: app
	./app serve