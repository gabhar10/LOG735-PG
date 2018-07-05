var express = require('express')
var ws = require('./scripts/ws')
var app = express()

app.use("/scripts", express.static(__dirname + "/frontend"));

app.get('/', function (req, res) {
    res.sendFile(__dirname + '/frontend/index.html');
});


app.listen(3000, function () {
  console.log('Webapp listening on port 3000')
});
