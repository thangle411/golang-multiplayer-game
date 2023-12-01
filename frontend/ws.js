//Websockets
const lobbySocket = new WebSocket('ws://0.0.0.0:8080/ws/joinLobby');
window.lobbySocket = lobbySocket;
lobbySocket.addEventListener('message', event => {
  const data = JSON.parse(event.data);
  if (data.content) {
    if (data.content === 'Welcome') {
      window.Store.playerid = data.playerid;
    } else {
      appendToLobbyChat(data.content, data.type, data.playerid);
    }
  } else if (data.state) {
    handleWorldState(data);
  }
});

const handleWorldState = data => {
  window.Store.world.players = data.state;
  window.EventBus.dispatch(window.EventBus.events['sync'], data);
};

const chat = event => {
  if (event.code === 'Enter') {
    const input = document.getElementById('lobby-input');
    lobbySocket.send(input.value);
    input.value = '';
  }
};

const appendToLobbyChat = (text, type, id) => {
  const container = document.getElementById('lobby-chat-messages');
  const div = document.createElement('div');
  div.innerText = id ? `${id}: ${text}` : text;
  if (type === 1) {
    div.style.color = 'grey';
  }
  container.appendChild(div);
};
