import { MatchResponse, Player } from "./api/types";
import { getSettingValue } from "./web-extension/settings";
import { getMatch } from "./api/get_match";
import { getPlayers } from "./api/get_players";
import { NoLogsError } from "./api/error";

const matchRe = RegExp("https://logs.tf/(\\d+)");

function getLogID(): number {
  const match = document.URL.match(matchRe);

  if (match === null || match.length < 1) {
    throw new Error("could not find log ID");
  }
  return parseInt(match[1]);
}

export async function addMatchLink() {
  const apiBaseUrl = getSettingValue<string>("apiBaseURL") as string;

  let matchId: number;

  try {
    matchId = getLogID();
  } catch {
    return;
  }

  let res: MatchResponse;

  try {
    res = await getMatch(apiBaseUrl, matchId);
  } catch (e) {
    if (e === NoLogsError) {
      return;
    }

    console.error("offi: could not get match: " + e.toString());
    return;
  }

  const competitionBlock = document.createElement("h3");

  const competitionLink = document.createElement("a")
  competitionLink.href = `https://etf2l.org/matches/${res.match.match_id}`
  competitionLink.innerText = res.match.competition;

  competitionBlock.appendChild(competitionLink)

  const matchBlock = document.createElement("h3");

  if (res.match.tier) {
    matchBlock.innerText = `${res.match.tier} | ${res.match.stage}`;
  } else {
    matchBlock.innerText = res.match.stage;
  }

  document.getElementById("log-date").after(competitionBlock, matchBlock);

  if (res.log.demo_id) {
    const demoLogo = document.createElement("img");
    demoLogo.className = "demostf-logo medium";

    const demoLink = document.createElement("a");
    demoLink.href = `https://demos.tf/${res.log.demo_id}`;
    demoLink.innerText = "demos.tf";

    const demoContainer = document.createElement("h3");
    demoContainer.append(demoLogo, demoLink);

    document.getElementById("log-map").after(demoContainer);
  }
}

export async function replacePlayerNames() {
  const apiBaseUrl = getSettingValue<string>("apiBaseURL") as string;

  const playerNodes = document.querySelectorAll("[id^=player_]");

  const steamPlayerNames: Map<string, string> = new Map();

  const playerIDs = Array.from(playerNodes).map((node) => {
    const steamId = node.id.replace("player_", "");
    const oldName = node.querySelector(".log-player-name a").textContent;

    steamPlayerNames.set(steamId, oldName);

    return steamId;
  });

  const players = await getPlayers(apiBaseUrl, playerIDs);

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
