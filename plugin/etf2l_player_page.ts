import {apiUrl} from "./utils";
import {PlayersResponse, Recruitment, Ban} from "./types";

const playerRe = RegExp("https://etf2l.org/forum/user/(\\d+)/");

function getPlayerID(): number {
  const match = document.URL.match(playerRe);

  if (match === null || match.length < 1) {
    throw new Error("could not find match ID");
  }
  return parseInt(match[1]);
}

async function getPlayerStatusFromAPI(playerId: number): Promise<PlayersResponse> {
  const getPlayerURL = new URL(apiUrl + "players");
  getPlayerURL.searchParams.append("id", playerId.toString());
  getPlayerURL.searchParams.append("version", chrome.runtime.getManifest().version);

  const res = await fetch(getPlayerURL.toString());

  if (!res.ok) {
    throw new Error("offi api returned error: " + res.statusText);
  }

  return (await res.json()) as PlayersResponse;
}

async function addPlayerStatus(recruitment: Recruitment) {
  if (recruitment.empty) {
    return;
  }

  const node = document.createElement("a");
  node.setAttribute("href", recruitment.url);
  node.className = "recruitment-status";
  node.innerText = `LFT ${recruitment.skill} ${recruitment.game_mode}`;

  document
      .querySelector("#rs-discuss")
      .querySelector("h2")
      .appendChild(node);

  // I love WordPress
  document
      .querySelector(".playerinfo")
      .querySelector("tbody")
      .querySelectorAll("tr")[1]
      .querySelectorAll("td")[1]
      .querySelectorAll("img")
      .forEach((imgNode) => {
        if (recruitment.classes.includes(imgNode.title)) {
          imgNode.className = "invert-img";
        }
      });
}

async function addPlayersBans(bans: Ban[]) {
  const container = document.createElement("div");
  container.className = "player-bans";

  const header = document.createElement("h2");
  header.innerText = "Bans";

  const banList = document.createElement("ul");
  bans.forEach((ban) => {
    const banStart = new Date(ban.start * 1000);
    const banEnd = new Date(ban.end * 1000);
    banList.appendChild(document.createElement("li")).innerHTML = `<b>${ban.reason}</b>: ${banStart.toLocaleDateString()} to ${banEnd.toLocaleDateString()}`;
  });
  container.appendChild(header);
  container.appendChild(banList);

  document.getElementById("rs-discuss").appendChild(container);
}

async function updatePlayerPage() {
  const playerId = getPlayerID();

  let response: PlayersResponse;

  try {
    response = await getPlayerStatusFromAPI(playerId);
  } catch (e) {
    console.error("failed to get player status: ", e.toString());
  }

  if (response.players.length != 1) {
    console.error("api returned more than 1 player");
    return;
  }

  const player = response.players[0];

  if (player.recruitment != null) {
    await addPlayerStatus(player.recruitment);
  }

  if (player.bans.length > 1) {
    await addPlayersBans(player.bans);
  }
}

updatePlayerPage();
