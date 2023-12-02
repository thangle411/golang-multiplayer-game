const Application = PIXI.Application;
const Container = PIXI.Container;
const Graphics = PIXI.Graphics;
const app = new Application({
  width: window.innerWidth,
  height: window.innerHeight,
  transparent: false,
});
const container = new Container({
  width: app.view.width,
  height: app.view.height,
});
const center = {
  x: app.view.width / 2,
  y: app.view.height / 2,
};
document.body.appendChild(app.view);
app.ticker.add(gameLoop);

let globalPlayers = {};

window.EventBus.subscribe(window.EventBus.eventNames['sync'], data => {
  // this is syncing with states from server
  const { state } = data;
  if (!state) return;
  state.forEach(p => {
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
      globalPlayers[id].position.x = p.state.Point.x;
      globalPlayers[id].position.y = p.state.Point.y;
    }
  });
});

window.EventBus.subscribe(window.EventBus.eventNames['removePlayer'], data => {
  app.stage.removeChild(globalPlayers[data.playerid]);
  delete globalPlayers[data.playerid];
});

function gameLoop() {
  Object.keys(globalPlayers).forEach(key => {
    const p = globalPlayers[key];
    if (Number(key) !== Number(window.Store.playerid)) return;
    if (window.Store.input['ArrowUp']) {
      p.y -= 1;
      lobbySocket.send(JSON.stringify({ key: 'arrow-up' }));
    }
    if (window.Store.input['ArrowDown']) {
      p.y += 1;
      lobbySocket.send(JSON.stringify({ key: 'arrow-down' }));
    }
    if (window.Store.input['ArrowLeft']) {
      p.x -= 1;
      lobbySocket.send(JSON.stringify({ key: 'arrow-left' }));
    }
    if (window.Store.input['ArrowRight']) {
      p.x += 1;
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
