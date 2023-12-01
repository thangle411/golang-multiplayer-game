const Application = PIXI.Application;
const Graphics = PIXI.Graphics;
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

let globalPlayers = [];

window.EventBus.subscribe(window.EventBus.events['sync'], data => {
  //this is syncing with states from server
  const { state } = data;
  globalPlayers = [];
  app.stage.removeChildren();
  state.forEach(p => {
    const { state } = p;
    const player = new Graphics();
    player
      .beginFill(getRandomColor())
      .drawRect(center.x + state.Point.x * 1, center.y + state.Point.y * 1, 10, 10)
      .endFill();
    app.stage.addChild(player);
    globalPlayers.push({ player, id: p.id });
  });
});

function gameLoop() {
  globalPlayers.forEach(p => {
    if (p.id !== window.Store.playerid) return;
    switch (window.Store.currentInput) {
      case 'ArrowUp':
        p.player.y -= 1;
        lobbySocket.send(JSON.stringify({ key: 'arrow-up' }));
        break;
      case 'ArrowDown':
        p.player.y += 1;
        lobbySocket.send(JSON.stringify({ key: 'arrow-down' }));
        break;
      case 'ArrowLeft':
        p.player.x -= 1;
        lobbySocket.send(JSON.stringify({ key: 'arrow-left' }));
        break;
      case 'ArrowRight':
        p.player.x += 1;
        lobbySocket.send(JSON.stringify({ key: 'arrow-right' }));
        break;
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
