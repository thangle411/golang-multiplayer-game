class EventBus {
  constructor() {
    this.events = {};
    this.eventNames = {
      sync: 'sync',
      removePlayer: 'removePlayer',
    };
  }

  dispatch(eventName, data) {
    const cbs = this.events[eventName];
    if (!cbs) return;
    for (let cb of cbs) {
      cb(data);
    }
  }

  subscribe(eventName, cb) {
    if (!this.events[eventName]) {
      this.events[eventName] = [];
    }
    this.events[eventName].push(cb);
  }
}

const eventBus = new EventBus();
window.EventBus = eventBus;

class Store {
  constructor() {
    this.world = {
      players: [],
    };
    this.input = {
      ArrowUp: false,
      ArrowDown: false,
      ArrowLeft: false,
      ArrowRight: false,
    };
    this.playerid = 0;
  }
}

const store = new Store();
window.Store = store;

//input handlers
addEventListener('keydown', e => {
  window.Store.input[e.key] = true;
});

addEventListener('keyup', e => {
  window.Store.input[e.key] = false;
});
