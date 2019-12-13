$(document).ready(() => {
  const connection = new WebSocket("ws://localhost:8000");

  // Log errors
  connection.onerror = error => {
    console.log("WebSocket Error " + error);
  };

  // Log messages from the server
  connection.onmessage = e => {
    const date = new Date();
    const v = e.data;
    $(".chat").append(`<li class="other">
    <div class="avatar">
      <img src="https://i.imgur.com/DY6gND0.png" draggable="false" />
    </div>
    <div class="msg">
      <p>${v}</p>
      <time>${date.getHours()}:${date.getMinutes()}</time>
    </div>
  </li>`);
  };

  $("#textareaInput").keypress(function(e) {
    const key = e.charCode;
    if (key !== 13) return;

    const toSend = $(this).val();
    const date = new Date();
    $(this).val("");

    connection.send(`{"to":"all","data":"${toSend}"}`);

    $(".chat").append(`
    <li class="self">
        <div class="avatar">
          <img src="https://i.imgur.com/HYcn9xO.png" draggable="false" />
        </div>
        <div class="msg">
          <p>${toSend}</p>
          <time>${date.getHours()}:${date.getMinutes()}</time>
        </div>
      </li>`);
  });
});
