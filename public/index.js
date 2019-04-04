window.addEventListener("load", function(evt) {
  var output = document.getElementById("output");
  var input = document.getElementById("input");
  var ws;

  const CMD_SERVER_HELLO = 0;
  const CMD_LIST_AGENTS = 1;
  const CMD_SERVER_MAX = 2;
  const CMD_SYSTEM_INFO = 3;

  var commandmap = {
    helloserver: CMD_SERVER_HELLO,
    listagents: CMD_LIST_AGENTS,
    systeminfo: CMD_SYSTEM_INFO
  };

  var print = function(message) {
    output.innerText = message;
  };

  document.getElementById("connect").onclick = function(evt) {
    if (ws) {
      return false;
    }
    ws = new WebSocket("ws://localhost:8080/websocket");
    ws.onopen = function(evt) {
      print("OPEN");
    };
    ws.onclose = function(evt) {
      print("CLOSE");
      ws = null;
    };
    ws.onmessage = function(evt) {
      print("RESPONSE: " + evt.data);
    };
    ws.onerror = function(evt) {
      print("ERROR: " + evt.data);
    };
    return false;
  };

  document.getElementById("send").onclick = function(evt) {
    if (!ws) {
      return false;
    }

    var allcmd = document.getElementById("input").value;
    var cmd = allcmd.split(" ");
    //helloserver;
    if (!(cmd[0] in commandmap)) {
      print("invalidcommand");
      return;
    }

    // if (cmd.lenght < 3) {
    //   return
    // }

    var response = {
      rident: "456456",
      cmdtype: commandmap[cmd[0]]
    };

    ws.send(JSON.stringify(response));
    return false;
  };
});
