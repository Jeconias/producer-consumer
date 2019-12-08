$(document).ready(() => {
  const connection = new WebSocket("ws://localhost:8000");

  // Log errors
  connection.onerror = error => {
    console.log("WebSocket Error " + error);
  };

  // Log messages from the server
  connection.onmessage = e => {
    const v = e.data;
    $("#list").append(`<p>${v}</p>`);
    console.log("Server: " + v);
  };

  $("#formNickname").submit(e => {
    e.preventDefault();
    const toSend = $("#toSend").val();

    connection.send(`{"to":"all","data":"${toSend}"}`);
    console.log(`{"to":"all","data":"${toSend}"}`);
  });
});
