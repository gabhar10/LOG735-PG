const { spawn } = require('child_process');
const { exec } = require('child_process');
const dc = spawn('tail', ['-f', '-n', '+1', 'docker-compose.logs']);

var WebSocketServer = require('ws').Server,
  wss = new WebSocketServer({port: 40510})

wss.on('connection', function (ws) {
  dc.stdout.on('data', (data) => {
    exec('cat docker-compose.logs | aha > frontend/dc-logs.htm');
    ws.send("update");
  });

  ws.on('message', function (message) {
    console.log('received: %s', message)
  });
})
