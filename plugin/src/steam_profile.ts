import { getSettingValue } from "./web-extension/settings";
import { getPlayers } from "./api/get_players";

export async function addPlayerLinks() {
  const apiBaseUrl = await getSettingValue("apiBaseURL");
  const steamID = document.querySelector("[id^=commentthread_Profile_]").id.split("_")[2];

  const players = await getPlayers(apiBaseUrl, [steamID]);
  if (players === null) {
    console.warn("offi: player does not have an etf2l account");
    return;
  }
  const player = players[0];

  const etf2lProfileNode = createItemListElement(
      "ETF2L",
      `https://etf2l.org/forum/user/${player.id}/`,
  );

  const logsTfProfileNode = createItemListElement(
      "Logs",
      `https://logs.tf/profile/${player.steam_id}`,
  );

  const itemListNode = document.querySelector(".profile_item_links");
  itemListNode.insertBefore(logsTfProfileNode, itemListNode.firstChild);
  itemListNode.insertBefore(etf2lProfileNode, itemListNode.firstChild);
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
