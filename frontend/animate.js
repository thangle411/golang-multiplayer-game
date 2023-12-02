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

//start button
const startButton = new PIXI.Text('Start!', {
  fontFamily: 'Arial',
  fontSize: 12,
  fill: 0xffffff,
});
startButton.anchor.set(0.5);
startButton.position.set(center.x, 10);
app.stage.addChild(startButton);

let globalPlayers = {};
let globalSquares = {};

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
      player.x = state.Point.x;
      player.y = state.Point.y;
      app.stage.addChild(player);
      globalPlayers[id] = player;
    } else {
      globalPlayers[id].position.x = (center.x / 400) * p.state.Point.x;
      globalPlayers[id].position.y = (center.y / 400) * p.state.Point.y;
    }
  });
});

window.EventBus.subscribe(window.EventBus.eventNames['removePlayer'], data => {
  app.stage.removeChild(globalPlayers[data.playerid]);
  delete globalPlayers[data.playerid];
});

function renderSquares(squares) {
  squares?.forEach(s => {
    const { point, id } = s;
    if (!globalSquares[id]) {
      globalSquares[id] = s;
      const square = new Graphics();
      square.beginFill('#FFBF00').drawRect(center.x, center.y, 20, 20).endFill();
      square.x = (center.x / 400) * point.x;
      square.y = (center.y / 400) * point.y;
      app.stage.addChildAt(square, 1);
    }
  });
}

function gameLoop() {
  if (window.Store.gameState.level != 0) {
    app.stage.removeChild(startButton);
    renderSquares(window.Store.gameState.squares);
  } else {
    app.stage.addChild(startButton);
  }
  Object.keys(globalPlayers).forEach(key => {
    const p = globalPlayers[key];
    const ratioX = center.x / 400;
    const ratioY = center.y / 400;
    if (Number(key) !== Number(window.Store.playerid)) return;
    if (window.Store.input['ArrowUp']) {
      p.y -= ratioY;
      lobbySocket.send(JSON.stringify({ key: 'arrow-up' }));
    }
    if (window.Store.input['ArrowDown']) {
      p.y += ratioY;
      lobbySocket.send(JSON.stringify({ key: 'arrow-down' }));
    }
    if (window.Store.input['ArrowLeft']) {
      p.x -= ratioX;
      lobbySocket.send(JSON.stringify({ key: 'arrow-left' }));
    }
    if (window.Store.input['ArrowRight']) {
      p.x += ratioX;
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
