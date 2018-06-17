# LOG735-PG

## Topology Creation
To create a valid _docker-compose.yaml_ file, use the _topology-creation.py_ script as follows:

```bash
# This script expects docker-compose.template to be in the present working directory
python3 topology-creation.py --miners <amount-miners> --clients <amount-clients> --malicious-miners <amount-malicious-miners>
```

Variable Name | Value Type | Description
--- | --- | ---
amount-miners | Int | Total number of miners in the network. A subset are anchor miners
amount-clients | Int | Total number of clients in the network
amount-malicious-miners | Int | Total number of malicious miners in the network. Is a subnet of *honest* miners.

This script will take as input _docker-compose.template_ and generate a valid _docker-compose.yaml_ file in the present working directory, **overwriting** any current copies. 