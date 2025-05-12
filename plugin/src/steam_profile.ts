import { getSettingValue } from "./web-extension/settings";
import { getRGLPlayer } from "./api/get_rgl_player";
import { getEtf2lPlayer } from "./api/get_etf2l_player";

export async function addPlayerLinks() {
  const apiBaseUrl = getSettingValue<string>("apiBaseURL") as string;
  const league = (getSettingValue<string>("league") as string).toLowerCase();

  const steamID = document.querySelector("[id^=commentthread_Profile_]").id?.split("_")[2];
  if (!steamID) {
    console.warn("offi: could not find steam ID");
    return;
  }

  let playerURL: string
  let playerSteamID: string

  if (league === "rgl") {
    const player = await getRGLPlayer(apiBaseUrl, steamID);
    playerURL = `https://rgl.gg/Public/PlayerProfile.aspx?p=${player.steam_id}`
    playerSteamID = player.steam_id
  } else {
    const player = await getEtf2lPlayer(apiBaseUrl, steamID);
    playerURL = `https://etf2l.org/forum/user/${player.id}/`
    playerSteamID = player.steam_id
  }

  if (playerURL === "") {
    console.warn("offi: player does not have a league account");
    return;
  }

  const leagueLink = createItemListElement(
    league.toUpperCase(),
    playerURL,
  );

  const logsTfProfileNode = createItemListElement(
      "Logs",
      `https://logs.tf/profile/${playerSteamID}`,
  );

  const itemListNode = document.querySelector(".profile_item_links");
  itemListNode.insertBefore(logsTfProfileNode, itemListNode.firstChild);
  itemListNode.insertBefore(leagueLink, itemListNode.firstChild);
}

function createItemListElement(label: string, href: string): Element {
  const labelNode = document.createElement("span");
  labelNode.className = "count_link_label";
  labelNode.textContent = label;

  const linkNode = document.createElement("a");
  linkNode.href = href;

  const containerNode = document.createElement("div");
  containerNode.setAttribute("data-panel", "{\"focusable\":true,\"clickOnActivate\":true}");
  containerNode.className = "profile_count_link ellipsis";

  linkNode.appendChild(labelNode);
  containerNode.appendChild(linkNode);

  return containerNode;
}
