const urlParams = new URLSearchParams(window.location.search);
const room = urlParams.get("room") || "general";
const username =
  urlParams.get("username") || "anonymous_" + Math.floor(Math.random() * 1000);

document.title = `${room} — Nebula`;
const wsProtocol = location.protocol === "https:" ? "wss" : "ws";
const ws = new WebSocket(
  `${wsProtocol}://${location.host}/ws?room=${room}&username=${encodeURIComponent(
    username
  )}`
);

const messages = document.getElementById("messages");
const input = document.getElementById("input");
const status = document.getElementById("status");
const form = document.getElementById("form");

function addMessage(text, type = "message") {
  const div = document.createElement("div");
  if (type === "system") div.className = "system";
  div.textContent = text;
  messages.appendChild(div);
  messages.scrollTop = messages.scrollHeight;
}

ws.onopen = () => {
  status.textContent = `Connected as ${username} in room « ${room} »`;
  addMessage(`You joined the room « ${room} »`, "system");
};

ws.onmessage = (event) => {
  addMessage(event.data);
};

ws.onclose = () => {
  status.textContent = "Disconnected";
  addMessage("You have been disconnected", "system");
};

form.onsubmit = (e) => {
  e.preventDefault();
  if (input.value.trim() === "") return;
  ws.send(input.value);
  input.value = "";
};
