all:
	go build -o main && ./main
dev:
	ENV=development \
	go run main.go