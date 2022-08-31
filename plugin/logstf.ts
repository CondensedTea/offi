import "regenerator-runtime/runtime";
import {api, apiUrl, type, replaceInText, getPlayers} from "./utils";
import {MatchResponse, Match} from "./types";

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
  logURL.searchParams.append("version", api.runtime.getManifest().version);
  logURL.searchParams.append("browser", type);

  const res = await fetch(logURL.toString());

  if (!res.ok) {
    throw new Error("api returned error: " + res.statusText);
  }

  const apiResponse = await res.json() as MatchResponse;

  return apiResponse.match;
}

async function addMatchLink(): Promise<void> {
  let matchId: number;

  try {
    matchId = getLogID();
  } catch (e) {
    return;
  }

  let match: Match;

  try {
    match = await getMatchFromAPI(matchId);
  } catch (e) {
    console.log("off: could not get match: " + e.toString());
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

async function replacePlayerNames() {
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

  const players = await getPlayers(playerIDs);

  players.forEach((player) => {
    const oldName = steamPlayerNames.get(player.steam_id);

    const logSelectionNode = document.querySelector("div#log-section-players");
    replaceInText(logSelectionNode, oldName, player.name);

    const healSpreadNode = document.querySelector("div.healspread");
    replaceInText(healSpreadNode, oldName, player.name);

    const tabContentNode = document.querySelector("div.tab-content");
    replaceInText(tabContentNode, oldName, player.name);

    const showstreaksNode = document.querySelector("#showstreaks");
    replaceInText(showstreaksNode, oldName, player.name);
  });
}

api.storage.sync.get((fields: Options) => {
  if (fields.logstf_link_matchpage === true) {
    addMatchLink();
  }
  if (fields.logstf_replace_names === true) {
    replacePlayerNames();
  }
});

