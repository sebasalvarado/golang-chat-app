$(function() {
  var socket = null;
  var msgBox = $("#chatbox textarea");
  var messages = $("#messages");

  // Handle function to submit message
  $("#chatbox").submit(function() {
    if (!msgBox.val()) return false;
    if (!socket) {
      alert("Error: There is no socket connection");
      return false
    }

    socket.send(msgBox.val());
    msgBox.val("")
    return false
  });

  // Check WebSocket support
  if (!window["WebSocket"]) {
    alert("Error: Browser does not support websocket")
  } else {
    socket = new WebSocket("ws://localhost:8080/room");
    socket.onclose = function() {
      alert("Connection closed");
    }
    socket.onmessage = function(e) {
      messages.append($("<li>").text(e.data));
    }
  }
});