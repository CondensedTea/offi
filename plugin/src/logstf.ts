import { Match, Player } from "./api/types";
import { getSettingValue } from "./web-extension/settings";
import { getMatch } from "./api/get_match";
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
    match = await getMatch(apiBaseUrl, matchId);
  } catch (e) {
    console.warn("offi: could not get match: " + e.toString());
    return;
  }

  const competitionBlock = document.createElement("h3");
  competitionBlock.innerHTML =
    `<a href="https://etf2l.org/matches/${match.match_id}">${match.competition}</a>`;

  const matchBlock = document.createElement("h3");

  if (match.tier) {
    matchBlock.innerText = `${match.tier} | ${match.stage}`;
  } else {
    matchBlock.innerText = match.stage;
  }

  const logDateElem = document.getElementById("log-date");

  logDateElem.after(matchBlock);
  logDateElem.after(competitionBlock);
}

export async function replacePlayerNames() {
  const apiBaseUrl = await getSettingValue("apiBaseURL");

  const playerNodes = document.querySelectorAll("[id^=player_]");

  const steamPlayerNames: Map<string, string> = new Map();

  const playerIDs = Array.from(playerNodes).map((node) => {
    const steamId = node.id.replace("player_", "");
    const oldName = node.querySelector(".log-player-name a").textContent;

    steamPlayerNames.set(steamId, oldName);

    return steamId;
  });

  const players = await getPlayers(apiBaseUrl, playerIDs, false);

  const steamIDToETF2LID = new Map<string, string>();
  players.map((player) => {
    steamIDToETF2LID.set(player.steam_id, player.id)
  });

  const selectors = [
    "#class_k .log-player-name",         // Player names in main table
    ".log-player-name .dropdown-toggle", // Player names in "kills" table
    ".healtable h6",                     // Medic names in heals table
    "td.log-player-name",                // Player names in heal target table
    ".chat-name"                         // Player names in chat
  ]

  replacePlayerNamesInNodes(players, steamPlayerNames, selectors.join(", "));

  // On toggle of death/kills tables trigger name replacement
  document.querySelectorAll("#classtab [data-toggle]").forEach((node) => {
    node.addEventListener("click", () => {
      replacePlayerNamesInNodes(players, steamPlayerNames, node.attributes["href"].value);
    });
  })

  document.querySelectorAll("a[href^='http://etf2l.org/search']").forEach((node) => {
    const match = node.getAttribute("href").match(/\/search\/(\d+)/);
    if (match.length < 2) return;

    const etf2l_id = steamIDToETF2LID.get(match[1])
    if (!etf2l_id) return;

    node.setAttribute("href", `https://etf2l.org/forum/user/${etf2l_id}/`);
  })
}

function replacePlayerNamesInNodes(players: Player[], steamPlayerNames: Map<string, string>, selector: string) {
  players.forEach((player) => {
    const oldName = steamPlayerNames.get(player.steam_id);

    const playerNameNodes = document.querySelectorAll(selector);
    playerNameNodes.forEach((node) => {
      replaceInText(node, oldName, player.name);
    });
  });
}

export function replaceInText(element: ChildNode, pattern: string, replacement: string) {
  for (const node of element.childNodes) {
    switch (node.nodeType) {
      case Node.ELEMENT_NODE || Node.DOCUMENT_NODE:
        replaceInText(node, pattern, replacement);
        break;
      case Node.TEXT_NODE:
        if (node.textContent === pattern) {
          node.textContent = replacement;
        }
        break;
    }
  }
}
