const Application = PIXI.Application;
const Container = PIXI.Container;
const Graphics = PIXI.Graphics;
const Text = PIXI.Text;
const app = new Application({
  width: window.innerWidth,
  height: window.innerHeight,
  transparent: false,
});
const center = {
  x: app.view.width / 2,
  y: app.view.height / 2,
};
document.body.appendChild(app.view);
app.ticker.add(gameLoop);

let globalPlayers = {};
let globalSquares = {};
let globalBots = {};

window.EventBus.subscribe(window.EventBus.eventNames['sync'], data => {
  // this is syncing with states from server
  const { worldState } = data;
  if (!worldState) return;
  worldState.forEach(p => {
    const { state, id } = p;
    //If player doesn't exist, then we add to the array
    if (!globalPlayers[id]) {
      const player = new Graphics();
      player.beginFill(getRandomColor()).drawRect(center.x, center.y, 10, 10).endFill();
      player.x = state.center.x;
      player.y = state.center.y;
      app.stage.addChild(player);
      globalPlayers[id] = player;
    } else {
      globalPlayers[id].position.x = (center.x / 400) * p.state.center.x;
      globalPlayers[id].position.y = (center.y / 400) * p.state.center.y;
    }
  });
});

window.EventBus.subscribe(window.EventBus.eventNames['removePlayer'], data => {
  app.stage.removeChild(globalPlayers[data.playerid]);
  delete globalPlayers[data.playerid];
});

function renderSquares(squares) {
  squares?.forEach(s => {
    const { state, id } = s;
    const currentSquare = globalSquares[id];
    if (!currentSquare || currentSquare?.x != state.center.x || currentSquare?.y != state.center.y) {
      app.stage.removeChild(currentSquare);
      const square = new Graphics();
      square.beginFill('#FFBF00').drawRect(center.x, center.y, state.width, state.height).endFill();
      square.x = (center.x / 400) * state.center.x;
      square.y = (center.y / 400) * state.center.y;
      globalSquares[id] = square;
      app.stage.addChildAt(square, 1);
    }
  });
}

function removeSquares() {
  Object.keys(globalSquares).forEach(key => {
    app.stage.removeChild(globalSquares[key]);
  });
  globalSquares = {};
}

function renderBots(bots) {
  console.log(bots);
  bots?.forEach(b => {
    const { state, id } = b;
    const currentBot = globalBots[id];
    if (!currentBot || currentBot?.x != state.center.x || currentBot?.y != state.center.y) {
      app.stage.removeChild(currentBot);
      const square = new Graphics();
      square.beginFill('#808080').drawRect(center.x, center.y, state.width, state.height).endFill();
      square.x = (center.x / 400) * state.center.x;
      square.y = (center.y / 400) * state.center.y;
      globalBots[id] = square;
      app.stage.addChildAt(square, 1);
    }
  });
}

function removeBots() {}

function gameLoop() {
  if (window.Store.gameState.squares?.length > 0) {
    const { Store } = window;
    renderSquares(Store.gameState.squares);
    renderBots(Store.gameState.bots);
  } else {
    removeSquares();
    removeBots();
  }

  Object.keys(globalPlayers).forEach(key => {
    const p = globalPlayers[key];
    const ratioX = center.x / 400;
    const ratioY = center.y / 400;
    if (Number(key) !== Number(window.Store.playerid)) return;
    if (window.Store.input['ArrowUp']) {
      p.y -= ratioY * 3;
      lobbySocket.send(JSON.stringify({ key: 'arrow-up' }));
    }
    if (window.Store.input['ArrowDown']) {
      p.y += ratioY * 3;
      lobbySocket.send(JSON.stringify({ key: 'arrow-down' }));
    }
    if (window.Store.input['ArrowLeft']) {
      p.x -= ratioX * 3;
      lobbySocket.send(JSON.stringify({ key: 'arrow-left' }));
    }
    if (window.Store.input['ArrowRight']) {
      p.x += ratioX * 3;
      lobbySocket.send(JSON.stringify({ key: 'arrow-right' }));
    }
  });
}

function getRandomColor() {
  var letters = '0123456789ABCDEF';
  var color = '#';
  for (var i = 0; i < 6; i++) {
    color += letters[Math.floor(Math.random() * 16)];
  }
  return color;
}
