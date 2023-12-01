// const enterBtn = document.getElementById('enter-button');
// enterBtn.addEventListener('click', async () => {
//   const name = document.getElementById('nickname');
//   if (!name.value) {
//     return alert('Please provide a name');
//   }
//   const resp = await fetch('http://0.0.0.0:8080/createGame', {
//     method: 'POST',
//     headers: {
//       'Content-Type': 'application/json',
//     },
//     body: JSON.stringify({ name: name.value }),
//   });
//   if (!resp.ok) {
//     alert('Maximum games allowed');
//     name.value = '';
//     return;
//   }

//   getGames();
// });

// //Fetching data
// const getGames = async () => {
//   const resp = await fetch('http://0.0.0.0:8080/getGames');
//   const data = await resp.json();
//   const roomsDiv = document.getElementById('rooms');
//   roomsDiv.innerHTML = '';
//   data.forEach(game => {
//     console.log(game);
//     const div = document.createElement('div');
//     div.className = 'game-container';
//     const innerDiv = document.createElement('div');
//     innerDiv.innerText = 'Name: ' + game.name;
//     const playerNum = document.createElement('div');
//     playerNum.innerText = '# of Players: ' + game.numberOfPlayers;

//     div.append(innerDiv, playerNum);
//     roomsDiv.append(div);
//   });
// };

// getGames();
// const gamesInterval = setInterval(async () => {
//   getGames();
// }, 60 * 1000);
