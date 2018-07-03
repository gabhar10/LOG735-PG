# LOG735-PG

## Build docker image

```bash
make build
```

## Create and run topolgy

```bash
make run CLIENTS=<amount-clients> MINERS=<amount-miners> MMINERS=<amount-malicious-miners>
```

Variable Name | Value Type | Description
--- | --- | ---
amount-miners | Integer | Total number of miners in the network. A subset are anchor miners
amount-clients | Integer | Total number of clients in the network
amount-malicious-miners | Integer | Total number of malicious miners in the network. Is a subnet of *honest* miners.


## Destroy topology
```bash
make destroy
```
