build:
	docker build -t log735:latest .
	cd topology; python3 topology-creation.py --miners 10 --clients 20 --malicious-miners 3

run:
	cd topology; docker-compose up

destroy:
	cd topology; docker-compose down
