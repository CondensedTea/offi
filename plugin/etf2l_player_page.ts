import {apiUrl, api, type} from "./utils";
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
  getPlayerURL.searchParams.append("version", api.runtime.getManifest().version);
  getPlayerURL.searchParams.append("browser", type);

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
    ban.startDate = new Date(ban.start * 1000);
    ban.endDate = new Date(ban.end * 1000);

    banList.appendChild(createBanEntryNode(ban))
  });
  container.appendChild(header);
  container.appendChild(banList);

  document.getElementById("rs-discuss").appendChild(container);
}

function createBanEntryNode(ban: Ban): HTMLLIElement {
  let node = document.createElement("li")

  let banReasonHTML = `<b>${ban.reason}</b>`
  let banCommentHTML = `${ban.startDate.toLocaleDateString()} to ${ban.endDate.toLocaleDateString()}`
  if (ban.end - ban.start < 0) {
    banCommentHTML = `<span style="text-decoration: line-through">${banCommentHTML}</span> reverted`
  }

  node.innerHTML = `${banReasonHTML}: ${banCommentHTML}`

  return node
}

async function updatePlayerPage(options: Options) {
  if (options.etf2l_show_bans === false && options.etf2l_show_lft === false) {
    return;
  }

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

  if (options.etf2l_show_lft && player.recruitment != null) {
    await addPlayerStatus(player.recruitment);
  }

  if (options.etf2l_show_bans && player.bans.length > 1) {
    await addPlayersBans(player.bans);
  }
}

api.storage.sync.get((fields) => {
  updatePlayerPage(fields);
});

