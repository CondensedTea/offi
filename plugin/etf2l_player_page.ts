import {apiUrl} from "./utils";

const playerRe = RegExp("https://etf2l.org/forum/user/(\\d+)/");

class PlayerStatus {
  id: number;
  skill: string;
  url: string;
  game_mode: string;
  classes: string[];
  empty: boolean;
}

class PlayerInfo {
  id: number;
  bans: {
      start: number,
      end: number,
      reason: string
  }[];
}

type ApiPlayerResponse = {
  status: PlayerStatus;
  player: PlayerInfo;
};

function getPlayerID(): number {
  const match = document.URL.match(playerRe);

  if (match === null || match.length < 1) {
    throw new Error("could not find match ID");
  }
  return parseInt(match[1]);
}

async function getPlayerStatusFromAPI(playerId: number): Promise<ApiPlayerResponse> {
  const getPlayerURL = new URL(apiUrl + "player/" + playerId.toString());
  getPlayerURL.searchParams.append("version", chrome.runtime.getManifest().version);

  const res = await fetch(getPlayerURL.toString());

  if (!res.ok) {
    throw new Error("offi api returned error: " + res.statusText);
  }

  return (await res.json()) as ApiPlayerResponse;
}

async function addPlayerStatus(status: PlayerStatus) {
  const node = document.createElement("a");
  node.setAttribute("href", status.url);
  node.className = "recruitment-status";
  node.innerText = `LFT ${status.skill} ${status.game_mode}`;

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
        if (status.classes.includes(imgNode.title)) {
          imgNode.className = "invert-img";
        }
      });
}

async function addPlayersBans(playerInfo: PlayerInfo) {
  const container = document.createElement("div");
  container.className = "player-bans";

  const header = document.createElement("h2");
  header.innerText = "Bans";

  const banList = document.createElement("ul");
  playerInfo.bans.forEach((ban) => {
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

  let player: ApiPlayerResponse;

  try {
    player = await getPlayerStatusFromAPI(playerId);
  } catch (e) {
    console.error("failed to get player status: ", e.toString());
  }

  if (player.status != null && !player.status.empty) {
    await addPlayerStatus(player.status);
  }

  if (player.player != null) {
    await addPlayersBans(player.player);
  }
}

updatePlayerPage();
