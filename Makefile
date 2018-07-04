destroy:
	-pkill -f topology.creation.py
	-cd topology; docker-compose down

build:
	docker build -t log735:latest .

run:	destroy
	cd topology; python3 topology-creation.py --miners $(MINERS) --clients $(CLIENTS) --malicious-miners $(MMINERS) &
	cd topology; python3 -m webbrowser -n "http://127.0.0.1:8000"
	cd topology; docker-compose up
