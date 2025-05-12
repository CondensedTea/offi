import { getSettingValue } from "./web-extension/settings";
import { getLogs } from "./api/get_logs";
import { NoLogsError } from "./api/error"
import { Log } from "./api/types";

const matchRe = RegExp("https://etf2l.org/matches/(\\d+)/");

function getMatchID(): number {
  const match = document.URL.match(matchRe);

  if (match === null || match.length < 1) {
    throw new Error("could not find match ID");
  }
  return parseInt(match[1]);
}

function createLogHeader(logList: Node, isPrimary: boolean) {
  const headerContainer = document.createElement("div");

  let text = "";

  if (!isPrimary) {
    text = "Other ";
  }
  if (logList.childNodes.length === 1) {
    text += "1 log";
  } else {
    text += `${logList.childNodes.length} logs`;
  }

  const header = document.createElement("h2");
  header.innerText = text;

  headerContainer.append(header, logList);

  document.querySelector(".match-players")?.after(headerContainer);
}

export async function addLogLinks(): Promise<void> {
  const matchId = getMatchID();
  let logs: Log[];

  const apiBaseUrl = getSettingValue<string>("apiBaseURL");

  try {
    logs = await getLogs(apiBaseUrl, matchId);
  } catch (e) {
    if (e === NoLogsError) {
      return;
    }
    console.error("could not get logs: " + e.toString());
    return;
  }

  const primaryLogs: Log[] = [];
  const secondaryLogs: Log[] = [];

  logs.forEach((log) => {
    if (log.is_secondary) {
      secondaryLogs.push(log);
    } else {
      primaryLogs.push(log);
    }
  });

  if (secondaryLogs.length > 0) {
    const secondaryLogListNode = document.createElement("ul");

    secondaryLogListNode.append(...secondaryLogs.map(buildLogList))

    createLogHeader(secondaryLogListNode, false);
  }

  if (primaryLogs.length > 0) {
    const primaryLogListNode = document.createElement("ul");

    primaryLogListNode.append(...primaryLogs.map(buildLogList))

    createLogHeader(primaryLogListNode, true);
  }
}

function buildLogList(log: Log): Node {
  const logItem = document.createElement("li");

  if (log.demo_id) {
    const demosLink = document.createElement("a");
    const demosLogo = document.createElement("img");
    demosLogo.className = "demostf-logo logo-small";
    demosLink.append(demosLogo);
    demosLink.href = "https://demos.tf/" + log.demo_id;
    logItem.append(demosLink);
  }

  const logLink = document.createElement("a");
  logLink.href = `https://logs.tf/${log.id}`
  logLink.innerText = `#${log.id}`
  logItem.append(logLink)

  if (log.map) {
    logItem.append(document.createTextNode(` | ${log.map}`))
  }

  logItem.append(document.createTextNode(` | ${log.played_at.toLocaleString()}`))

  return logItem
}
