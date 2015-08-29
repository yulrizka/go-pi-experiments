var WebSocket = require('ws') , 
ws = new WebSocket('ws://192.168.1.43:8080/button');
  ws.on('open', function() {
    //noop
  });
  ws.on('message', function(message) {
    console.log('received: %s', message);
  });
