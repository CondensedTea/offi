import { getSettingValue } from "@kocal/web-extension-library";
import { getLogs, NoLogsError } from "./api/get_logs";
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
  const LogHeader = document.createElement("div");
  LogHeader.className = "offi";
  let text = "";

  if (!isPrimary) {
    text = "Other ";
  }
  if (logList.childNodes.length === 1) {
    text += "1 log";
  } else {
    text += `${logList.childNodes.length} logs`;
  }
  LogHeader.innerHTML = `<h2>${text}</h2>`;

  LogHeader.append(logList);

  const playersSection = document.getElementsByClassName("fix match-players");
  if (playersSection === null || playersSection.length < 1) {
    return;
  }
  playersSection[0].after(LogHeader);
}

export async function addLogLinks(): Promise<void> {
  const matchId = getMatchID();
  let logs: Log[];

  const apiBaseUrl = await getSettingValue("apiBaseURL");

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
      secondaryLogs.unshift(log);
    } else {
      primaryLogs.unshift(log);
    }
  });

  if (secondaryLogs.length > 0) {
    const OtherLogList = document.createElement("ul");
    secondaryLogs.forEach((log) => {
      const logItem = document.createElement("li");
      logItem.innerHTML = `<a href="https://logs.tf/${log.id}">#${log.id}</a> | ${log.map} | ${log.played_at.toLocaleString()}`;
      OtherLogList.appendChild(logItem);
    });
    createLogHeader(OtherLogList, false);
  }

  if (primaryLogs.length > 0) {
    const PrimaryLogList = document.createElement("ul");
    primaryLogs.forEach((log) => {
      const logItem = document.createElement("li");
      logItem.innerHTML = `<a href="https://logs.tf/${log.id}">#${log.id}</a> | ${log.map} | ${log.played_at.toLocaleString()}`;
      PrimaryLogList.appendChild(logItem);
    });
    createLogHeader(PrimaryLogList, true);
  }
}
