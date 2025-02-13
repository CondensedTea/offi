import { Player, Recruitment, Ban } from "./api/types";
import { getSettingValue } from "./web-extension/settings";
import { getPlayers } from "./api/get_players";

const playerRe = RegExp("https://etf2l.org/forum/user/(\\d+)/");

function getPlayerID(): string {
  const match = document.URL.match(playerRe);

  if (match === null || match.length < 1) {
    throw new Error("could not find match ID");
  }
  return match[1];
}

async function addPlayerStatus(recruitment: Recruitment) {
  const node = document.createElement("a");
  node.href = recruitment.url;
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

    banList.appendChild(createBanEntryNode(ban));
  });
  container.appendChild(header);
  container.appendChild(banList);

  document.getElementById("rs-discuss").appendChild(container);
}

function createBanEntryNode(ban: Ban): HTMLLIElement {
  const node = document.createElement("li");

  let banCommentHTML = `${ban.startDate.toLocaleDateString()} to ${ban.endDate.toLocaleDateString()}`;
  if (ban.end - ban.start < 0) {
    banCommentHTML = `<span style="text-decoration: line-through">${banCommentHTML}</span> reverted`;
  }

  node.innerHTML = `<b>${ban.reason}</b>: ${banCommentHTML}`;

  return node;
}

export async function updatePlayerPage(showLft: boolean, showBans: boolean) {
  const apiBaseUrl = await getSettingValue("apiBaseURL");
  const playerId = getPlayerID();

  let players: Player[];

  try {
    players = await getPlayers(apiBaseUrl, [playerId], true);
  } catch (e) {
    console.error("failed to get player status: ", e.toString());
  }

  if (players.length != 1) {
    console.error("api returned more than 1 player");
    return;
  }

  const player = players[0];

  if (showLft && player.recruitment != null) {
    await addPlayerStatus(player.recruitment);
  }

  if (showBans && player.bans.length > 0) {
    await addPlayersBans(player.bans);
  }
}
