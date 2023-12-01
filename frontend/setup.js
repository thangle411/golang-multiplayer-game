class EventBus {
  events = {
    sync: 'sync',
  };
  constructor() {
    this.events = {};
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
    this.currentInput = null;
    this.playerid = 0;
  }
}

const store = new Store();
window.Store = store;

//input handlers
addEventListener('keydown', e => {
  window.Store.currentInput = e.key;
});

addEventListener('keyup', e => {
  window.Store.currentInput = null;
});
