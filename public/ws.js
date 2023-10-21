window.onload = main

function writeOutput(msg) {
  const outputElement = document.getElementById("output");
  outputElement.innerHTML += `<p>${msg}</p><br>`;
}

function main() {
  const inputElement = document.getElementById("text");
  const sendElement = document.getElementById("send");

  inputElement.addEventListener("keypress", (evt) => {
    if (evt.key === "Enter") {
      evt.preventDefault();
      sendElement.click();
    }
  });

  const socket = new WebSocket("ws:\/\/127.0.0.1:3000/ws/board");

  socket.addEventListener("open", (evt) => {
    console.log(evt)
    writeOutput("socket opened");
  });

  socket.addEventListener("message", (evt) => {
    console.log(evt)
    console.log("message from server: ", evt.data)
    writeOutput(evt.data)
  });

  socket.addEventListener("close", () => {
    writeOutput("socket closed");
  });

  sendElement.addEventListener("click", () => {
    const message = inputElement.value
    inputElement.value = "";
    if (socket.readyState == socket.OPEN) socket.send(message);
  })
}

