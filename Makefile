run:
	go run ./main.go

docker-build:
	docker build -t github.com/tomkaith13/decision-engine-poc .

docker-run:
	docker run -p 8080:8080 github.com/tomkaith13/decision-engine-poc

docker-stop:
	docker stop github.com/tomkaith13/decision-engine-poc

docker-clean:
	docker rmi -f github.com/tomkaith13/decision-engine-poc