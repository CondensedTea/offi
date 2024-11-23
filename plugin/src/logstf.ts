import "regenerator-runtime/runtime";
import { Match } from "./api/types";
import { getSettingValue } from "@kocal/web-extension-library";
import { getMatchFromAPI } from "./api/get_match";
import { getPlayers } from "./api/get_players";

const matchRe = RegExp("https://logs.tf/(\\d+)");

function getLogID(): number {
  const match = document.URL.match(matchRe);

  if (match === null || match.length < 1) {
    throw new Error("could not find log ID");
  }
  return parseInt(match[1]);
}

export async function addMatchLink() {
  const apiBaseUrl = await getSettingValue("apiBaseURL");

  let matchId: number;

  try {
    matchId = getLogID();
  } catch {
    return;
  }

  let match: Match;

  try {
    match = await getMatchFromAPI(apiBaseUrl, matchId);
  } catch (e) {
    console.warn("off: could not get match: " + e.toString());
    return;
  }

  const competitionBlock = document.createElement("h3");
  competitionBlock.innerHTML =
    `<a href="https://etf2l.org/matches/${match.match_id}">${match.competition}</a>`;

  const matchBlock = document.createElement("h3");
  matchBlock.innerText = match.stage;

  const logDateElem = document.getElementById("log-date");

  logDateElem.after(matchBlock);
  logDateElem.after(competitionBlock);
}

export async function replacePlayerNames() {
  const apiBaseUrl = await getSettingValue("apiBaseURL");

  const match = document.URL.match(matchRe);
  if (match === null) {
    return;
  }

  const playerNodes = document.querySelectorAll("[id^=player_]");

  const steamPlayerNames: Map<string, string> = new Map();

  const playerIDs = Array.from(playerNodes).map((node) => {
    const steamId = node.id.replace("player_", "");
    const oldName = node.querySelector(".log-player-name a").textContent;

    steamPlayerNames.set(steamId, oldName);

    return steamId;
  });

  const players = await getPlayers(apiBaseUrl, playerIDs);

  players.forEach((player) => {
    const oldName = steamPlayerNames.get(player.steam_id);

    const playerNameNodes = document.querySelectorAll(
        // select player name cells or heal table headers or player names in chat
        ".log-player-name, .healtable h6, .chat-name",
    );
    playerNameNodes.forEach((node) => {
      replaceInText(node, oldName, player.name);
    });
  });
}

export function replaceInText(element, pattern, replacement) {
  for (const node of element.childNodes) {
    switch (node.nodeType) {
      case Node.ELEMENT_NODE || Node.DOCUMENT_NODE:
        replaceInText(node, pattern, replacement);
        break;
      case Node.TEXT_NODE:
        node.textContent = node.textContent.replace(pattern, replacement);
        break;
    }
  }
}
