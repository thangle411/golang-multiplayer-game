//Buttons
const startBtn = document.getElementById('start-btn');
const endBtn = document.getElementById('end-btn');
startBtn.addEventListener('click', () => {
  fetch('http://0.0.0.0:8080/startGame', {
    method: 'POST',
  });
});
endBtn.addEventListener('click', () => {
  fetch('http://0.0.0.0:8080/endGame', {
    method: 'POST',
  });
});

//Websockets
const lobbySocket = new WebSocket('ws://0.0.0.0:8080/ws/joinLobby');
window.lobbySocket = lobbySocket;
lobbySocket.addEventListener('message', event => {
  const data = JSON.parse(event.data);
  // console.log(data);
  if (data.content) {
    handleNonStateUpdate(data);
  } else if (data.worldState && data.gameState) {
    handleWorldState(data);
  }
});

const handleNonStateUpdate = data => {
  if (data.content === 'Welcome') {
    window.Store.playerid = data.playerid;
  } else if (data.content === 'Disconnect') {
    window.EventBus.dispatch(window.EventBus.eventNames['removePlayer'], data);
  } else {
    appendToLobbyChat(data.content, data.type, data.playerid);
  }
};

const handleWorldState = data => {
  window.Store.world.players = data.worldState;
  window.Store.gameState = data.gameState;
  window.EventBus.dispatch(window.EventBus.eventNames['sync'], data);
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
