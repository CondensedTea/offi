import { addLogLinks } from "./etf2l_match_page";
import { updatePlayerPage } from "./etf2l_player_page";
import { addPlayerLinks } from "./steam_profile";
import { registerSettings, loadSettings, getSettingValue } from "./web-extension/settings";
import settings from "./options";
import { addMatchLink, replacePlayerNames } from "./logstf";
import { addTeamStatus } from "./etf2l_team_page";

async function main() {
  await registerSettings(settings);
  await loadSettings();

  const url = document.URL;

  if (url.match("https://etf2l.org/matches/\\d+/") && getSettingValue("linkLogsOnMatchpages") as boolean) {
    return await addLogLinks();
  }

  if (url.match(`https://logs.tf/\\d+`)) {
    const replaceNames = getSettingValue("replaceNamesInLogs") as boolean;
    const linkMatchpages = getSettingValue("linkMatchepagesOnLogs") as boolean;

    if (replaceNames) {
      await replacePlayerNames();
    }

    if (linkMatchpages) {
      await addMatchLink();
    }

    return;
  }

  if (url.startsWith("https://etf2l.org/forum/user/")) {
    const showBans = getSettingValue<boolean>("showBansForPlayers");
    const showLft = getSettingValue("showLftForPlayer") as boolean;

    if (!showBans && !showLft) return;

    return await updatePlayerPage(showLft, showBans);
  }

  if (url.match("https://etf2l.org/teams/\\d+/")) {
    const showTeamStatus = getSettingValue("showLfpForTeam") as boolean;

    if (!showTeamStatus) return;

    return await addTeamStatus();
  }

  if (url.startsWith("https://steamcommunity.com/profiles/") || url.startsWith("https://steamcommunity.com/id/")) {
    const addLinks = getSettingValue("addLinksOnSteamProfiles") as boolean;

    if (!addLinks) return;

    return await addPlayerLinks();
  }
}

main();
