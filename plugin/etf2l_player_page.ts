import {apiUrl} from "./utils";

const playerRe = RegExp("https://etf2l.org/forum/user/(\\d+)/");
const NoRecruitmentInfo = new Error("this user doesn't have recruitment post");

class PlayerStatus {
  id: number;
  skill: string;
  url: string;
  game_mode: string;
  classes: string[];
  empty: boolean;
}

type ApiPlayerResponse = {
  status: PlayerStatus;
};

function getPlayerID(): number {
  const match = document.URL.match(playerRe);

  if (match === null || match.length < 1) {
    throw new Error("could not find match ID");
  }
  return parseInt(match[1]);
}

async function getPlayerStatusFromAPI(playerId: number): Promise<PlayerStatus> {
  const getPlayerURL = new URL(apiUrl + "player/" + playerId.toString());
  getPlayerURL.searchParams.append("version", chrome.runtime.getManifest().version);

  const res = await fetch(getPlayerURL.toString());

  if (!res.ok) {
    throw new Error("offi api returned error: " + res.statusText);
  }

  const playerStatus = (await res.json()) as ApiPlayerResponse;
  if (playerStatus.status === null || playerStatus.status.empty) {
    throw NoRecruitmentInfo;
  }

  return playerStatus.status;
}

async function addPlayerStatus() {
  const playerId = getPlayerID();

  let playerStatus: PlayerStatus;

  try {
    playerStatus = await getPlayerStatusFromAPI(playerId);
  } catch (e) {
    if (e === NoRecruitmentInfo) {
      return;
    } else {
      console.error("failed to get player status: ", e.toString());
    }
  }

  const node = document.createElement("a");
  node.setAttribute("href", playerStatus.url);
  node.className = "recruitment-status";
  node.innerText = `LFT ${playerStatus.skill} ${playerStatus.game_mode}`;

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
        if (playerStatus.classes.includes(imgNode.title)) {
          imgNode.className = "invert-img";
        }
      });
}

addPlayerStatus();
