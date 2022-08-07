import "regenerator-runtime/runtime";
import {api, apiUrl} from "./utils";
import {MatchResponse, Match, Player, PlayersResponse} from "./types";

const matchRe = RegExp("https://logs.tf/(\\d+)");

function getLogID(): number {
  const match = document.URL.match(matchRe);

  if (match === null || match.length < 1) {
    throw new Error("could not find log ID");
  }
  return parseInt(match[1]);
}

async function getMatchFromAPI(matchId: number): Promise<Match> {
  const logURL = new URL(apiUrl + "log/" + matchId.toString());
  logURL.searchParams.append("version", api().runtime.getManifest().version);

  const res = await fetch(logURL.toString());

  if (!res.ok) {
    throw new Error("offi api returned error: " + res.statusText);
  }

  const apiResponse = await res.json() as MatchResponse;

  return apiResponse.match;
}

async function getPlayers(ids: string[]): Promise<Player[]> {
  const playersURL = new URL(apiUrl + "players");

  const idsString = ids.join(",");

  playersURL.searchParams.append("id", idsString);
  playersURL.searchParams.append("version", browser.runtime.getManifest().version);

  const res = await fetch(playersURL.toString());
  if (!res.ok) {
    throw new Error("offi api returned error: " + res.statusText);
  }

  const response = await res.json() as PlayersResponse;
  return response.players;
}

async function addMatchLink(): Promise<void> {
  const matchId = getLogID();
  let match: Match;

  try {
    match = await getMatchFromAPI(matchId);
  } catch (e) {
    console.error("could not get match: " + e.toString());
    return;
  }

  const competitionBlock = document.createElement("h3");
  competitionBlock.innerHTML =
    `<a href="https://etf2l.org/matches/${match.id}">${match.competition}</a>`;

  const matchBlock = document.createElement("h3");
  matchBlock.innerText = match.stage;

  const logDateElem = document.getElementById("log-date");

  logDateElem.after(matchBlock);
  logDateElem.after(competitionBlock);
}

async function replacePlayerNames() {
  const playerNodes = document.querySelectorAll("[id^=player_]");

  const playerIDs = Array.from(playerNodes).map((node) => {
    return node.id.replace("player_", "");
  });

  const players = await getPlayers(playerIDs);

  players.forEach((player) => {
    const query = `#player_${player.steam_id} td.log-player-name div.dropdown a.dropdown-toggle`;

    const playerNode = document.querySelector(query);
    playerNode.innerHTML = player.name;
  });
}

addMatchLink();
replacePlayerNames();
