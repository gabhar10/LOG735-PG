build:
	docker build -t log735:latest .

run:
	cd topology; python3 topology-creation.py --miners $(MINERS) --clients $(CLIENTS) --malicious-miners $(MMINERS)
	cd topology; docker-compose up

destroy:
	cd topology; docker-compose down
