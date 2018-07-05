# LOG735-PG

## Build docker image

```bash
make build CLIENTS=<amount-clients> MINERS=<amount-miners> MMINERS=<amount-malicious-miners>
```

Variable Name | Value Type | Description
--- | --- | ---
amount-miners | Integer | Total number of miners in the network. A subset are anchor miners
amount-clients | Integer | Total number of clients in the network
amount-malicious-miners | Integer | Total number of malicious miners in the network. Is a subnet of *honest* miners.

## Create and run topolgy

```bash
make run
```

## Destroy topology
```bash
make clean
```

## Inspect logs
The logs are visible at http:localhost:3000

You can also consult them with
```bash
docker-compose -f topology/docker-compose.yaml logs -f
```

## Example
```bash
make clean; make build CLIENTS=20 MINERS=5 MMINERS=2; make run
```
