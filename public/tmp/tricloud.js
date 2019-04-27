class TriCloud {
  constructor() {
    this.commandmaps = [];
  }

  init() {
    this.ws = new WebSocket("ws://localhost:8080/websocket");
    this.ws.onopen = function(evt) {
      Print("OPEN");
    };

    this.ws.onclose = function(evt) {
      Print("CLOSE");
      ws = null;
    };
    this.ws.onerror = function(evt) {
      Print("ERROR: " + evt.data);
    };

    this.ws.onmessage = this.receive;
  }

  ParseHeader(data) {}

  receive(msg) {
    header = ParseHeader(msg);
    callback = this.commandmaps[str(header.msgtype)];
    callback(msg);
  }

  RegisterReceiver(key, callback) {
    this.commandmaps[key] = callback;
  }

  Send(msgtype, msg) {}

  SendAgent(msgtype, agentid, msg) {}

  newHead(connid, msgtype) {}
}

/* Utilities funcs */

function Print(arg) {
  console.log(arg);
}

function str(msgtype) {
  /*
    "terminal"
    "systemstat"
    */
  const mapping = {
    1: "serverhello",
    2: "terminal",
    3: "services"
  };
  return mapping[msgtype];
}
