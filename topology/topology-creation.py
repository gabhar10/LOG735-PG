import argparse, random, sys

# Define topology variables
MAX_MINER_PEERS=3
ANCHOR_PEERS=3
MIN_PORT = 7000
MAX_PORT = 7499
USED_PORTS = []
ANCHOR_MINERS = []
MINERS = []
CLIENTS = []

# Parse command-line arguments
parser = argparse.ArgumentParser(description='Generate blockchain demo topology.')
parser.add_argument('--miners', dest='num_miners', help='total number of miners in the network',
                    type=int, required=True)
parser.add_argument('--clients', dest='num_clients', help='total number of clients in the network',
                    type=int, required=True)
parser.add_argument('--malicious-miners', dest='num_mal_miners', help='total number of malicious miners in the network',
                    type=int, required=True)
args = parser.parse_args()

# Sanity check for minimal amount of anchor-miners and number of malicious miners
if args.num_miners < ANCHOR_PEERS:
    print('Number of miners must be greater or equal than %d' % ANCHOR_PEERS)
    sys.exit(1)

if args.num_mal_miners >= args.num_miners:
    print('Number of malicious miners must be smaller than miners')
    sys.exit(1)

# Define anchor-miners
while len(USED_PORTS) < ANCHOR_PEERS:
    anchor_port=random.randrange(MIN_PORT, MAX_PORT)
    if anchor_port not in USED_PORTS:
        node = {'ID': 'AM%d' % anchor_port, 'port': anchor_port, 'role': 'anchor-miner', 'peers': []}
        ANCHOR_MINERS.append(node)
        USED_PORTS.append(anchor_port)

# Define miners
while len(USED_PORTS) < args.num_miners:
    port=random.randrange(MIN_PORT, MAX_PORT)
    if port not in USED_PORTS:
        node = {'ID': 'M%d' % port, 'port': port, 'role': 'miner', 'malicious': False, 'peers': []}
        MINERS.append(node)
        USED_PORTS.append(port)

# Define clientd
while len(USED_PORTS) < args.num_miners + args.num_clients:
    port=random.randrange(MIN_PORT, MAX_PORT)
    if port not in USED_PORTS:
        node = {'ID': 'C%d' % port, 'port': port, 'role': 'client', 'peers': []}
        CLIENTS.append(node)
        USED_PORTS.append(port)
		
# Fully connect anchor-miners together
for i in ANCHOR_MINERS:
    for j in ANCHOR_MINERS:
        if i != j and j['port'] not in i['peers']:
            i['peers'].append(j['port'])

# Connect miners to at least 1 anchor-miner and other miners
for i in MINERS:
    i['peers'].append(random.choice(ANCHOR_MINERS)['port'])
    # Find remaining peers from concatenation of anchors-miners and miners MINUS the current miner
    sample = random.sample(ANCHOR_MINERS + list(filter(lambda x: x['port'] != i['port'], MINERS)) , MAX_MINER_PEERS-1)
    for j in sample:
        i['peers'].append(j['port'])

# Connect clients to a random anchor
for i in CLIENTS:
    i['peers'].append(random.choice(ANCHOR_MINERS)['port'])
	#for anchor in ANCHOR_MINERS:
	#	i['peers'].append(anchor['port'])

# Randomly select malicious miners
for i in random.sample(MINERS, args.num_mal_miners):
    i['malicious'] = True

# Create all docker-compose client services and vis.js content
services = ''
visjs_vertices = ''
visjs_edges = ''
for i in CLIENTS:
    services += '%s:\n  image: log735:latest\n  container_name: node-%s\n  environment:\n\
    - PEERS=%s\n    - ROLE=client\n    - PORT=%s\n  networks:\n    - blockchain\n' \
                % (i['ID'], i['port'], " ".join(str(x) for x in i['peers']), i['port'])
    visjs_vertices += '        {id: %s, label: \'%s\'},\n' % (i['port'], i['ID'])

    for x in i['peers']:
        visjs_edges += '        {from: %s, to: %s},\n' % (i['port'], x)

for i in ANCHOR_MINERS:
    services += '%s:\n  image: log735:latest\n  container_name: node-%s\n  environment:\n\
    - PEERS=%s\n    - ROLE=miner\n    - PORT=%s\n  networks:\n    - blockchain\n' \
                % (i['ID'], i['port'], " ".join(str(x) for x in i['peers']), i['port'])
    visjs_vertices += '        {id: %s, label: \'%s\'},\n' % (i['port'], i['ID'])
    for x in i['peers']:
        visjs_edges += '        {from: %s, to: %s},\n' % (i['port'], x)

for i in MINERS:
    services += '%s:\n  image: log735:latest\n  container_name: node-%s\n  environment:\n\
    - PEERS=%s\n    - ROLE=miner\n    - PORT=%s\n  networks:\n    - blockchain\n' \
                % (i['ID'], i['port'], " ".join(str(x) for x in i['peers']), i['port'])
    visjs_vertices += '        {id: %s, label: \'%s\'},\n' % (i['port'], i['ID'])
    for x in i['peers']:
        visjs_edges += '        {from: %s, to: %s},\n' % (i['port'], x)

# Add chat container
services += 'node-8001:\n  image: log735-chat:latest\n  container_name: chat_app\n  environment:\n\
    - PEERS=%s\n    - ROLE=chat_interface\n    - PORT=8001\n  networks:\n    - blockchain\n' \
% (CLIENTS[1]['port'])


# Write docker-compose.yaml
with open('docker-compose.template', 'r') as f:
    content = f.read()

content = content.replace('%SERVICES%', ''.join('  '+line for line in services.splitlines(True)))

with open('docker-compose.yaml', 'w+') as f:
    f.write(content)

# Write index.html
with open('../webapp/scripts/index.html.template', 'r') as f:
    content = f.read()
	
# Write chat application environment variables
with open('../chat/environment_variable.list', 'w') as f:
	f.write('PEERS=%s' % CLIENTS[1]['port'])
	

content = content.replace('%VERTICES%', ''.join(''+line for line in visjs_vertices[:-2].splitlines(True)))
content = content.replace('%EDGES%', ''.join(''+line for line in visjs_edges[:-2].splitlines(True)))

with open('../webapp/frontend/index.html', 'w+') as f:
    f.write(content)
