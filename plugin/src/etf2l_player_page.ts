import { Player, Recruitment, Ban } from "./api/types";
import { getSettingValue } from "./web-extension/settings";
import { getEtf2lPlayer } from "./api/get_etf2l_player";

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
          imgNode.alt = "This player is looking for a team";
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

  const banHeader = document.createElement("b");
  banHeader.innerText = ban.reason;

  node.append(banHeader, ": ");

  if (ban.end - ban.start > 0) {
    node.append(ban.startDate.toLocaleDateString(), " to ", ban.endDate.toLocaleDateString());
  } else {
    const revertedBan = document.createElement("span");
    revertedBan.innerText = ban.startDate.toLocaleDateString() + " to " + ban.endDate.toLocaleDateString();
    revertedBan.setAttribute("style", "text-decoration: line-through");
    node.append(revertedBan, " reverted");
  }

  return node;
}

export async function updatePlayerPage(showLft: boolean, showBans: boolean) {
  const apiBaseUrl = getSettingValue<string>("apiBaseURL") as string;
  const playerID = getPlayerID();

  let player: Player;
  try {
    player = await getEtf2lPlayer(apiBaseUrl, playerID, true);
  } catch (e) {
    console.error("failed to get player status: ", e.toString());
    return
  }

  if (showLft && player.recruitment != null) {
    await addPlayerStatus(player.recruitment);
  }

  if (showBans && player.bans.length > 0) {
    await addPlayersBans(player.bans);
  }
}
