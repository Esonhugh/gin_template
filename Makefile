app = "./app"

build: 
	go build -o app .

run:
	$(app) serve