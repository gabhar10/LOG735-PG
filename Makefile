clean:
	-cd topology; docker-compose down
	-docker rm -f log735-webapp
	-docker rm -f log735-chat
	-rm -f webapp/docker-compose.logs
build:
	cd topology; python3 topology-creation.py --miners $(MINERS) --clients $(CLIENTS) --malicious-miners $(MMINERS)
	docker build -t log735:latest .
	docker build -f Dockerfile_chat -t log735-chat:latest .
	cd webapp; docker build -t log735-webapp:latest .
run:
	cd topology; docker-compose up -d
	docker-compose -f topology/docker-compose.yaml logs -f > webapp/docker-compose.logs &
	docker run -dt -v $(shell pwd)/webapp/docker-compose.logs:/root/webapp/docker-compose.logs -p 3000:3000 -p 40510:40510 --privileged --name log735-webapp log735-webapp:latest
	python3 -m webbrowser -n "http://localhost:3000"
	python3 -m webbrowser -n "http://localhost:8000"
