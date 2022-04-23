import "regenerator-runtime/runtime";

const matchRe = RegExp("https://etf2l.org/matches/(\\d+)/");
const apiUrl = "https://offi.lemontea.dev/match/";

class Log {
  id: number;
  title: string;
  map: string;
  played_at: Date;
  is_secondary: boolean;

  constructor(data: Log) {
    this.id = data.id;
    this.title = data.title;
    this.map = data.map;
    this.played_at = new Date(data.played_at);
    this.is_secondary = data.is_secondary;
  }
}

type ApiLogResponse = {
  logs: Log[];
};

function getMatchID(): number {
  const match = document.URL.match(matchRe);

  if (match === null || match.length < 1) {
    throw new Error("could not find match ID");
  }
  return parseInt(match[1]);
}

async function getLogsFromAPI(matchId: number): Promise<Log[]> {
  const res = await fetch(apiUrl + matchId.toString());

  if (!res.ok) {
    throw new Error("offi api returned error: " + res.statusText);
  }

  const { logs } = (await res.json()) as ApiLogResponse;
  const parsedLogs: Log[] = [];

  for (const rawLog of logs) {
    parsedLogs.push(new Log(rawLog));
  }

  return parsedLogs;
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

async function addLogLinks(): Promise<void> {
  const matchId = getMatchID();
  let logs: Log[];

  try {
    logs = await getLogsFromAPI(matchId);
  } catch (e) {
    console.error("could not get logs: " + e.toString());
    return;
  }

  const PrimaryLogList = document.createElement("ul");
  const OtherLogList = document.createElement("ul");

  for (const log of logs) {
    const logItem = document.createElement("li");
    if (log.is_secondary) {
      logItem.innerHTML = `
        <a href="https://logs.tf/${log.id}">#${log.id}</a> | ${
        log.title !== "" ? log.title : log.map
      } | ${log.played_at.toLocaleString()}`;
      OtherLogList.appendChild(logItem);
    } else {
      logItem.innerHTML = `
        <a href="https://logs.tf/${log.id}">#${log.id}</a> | ${
        log.map
      } | ${log.played_at.toLocaleString()}`;
      PrimaryLogList.appendChild(logItem);
    }
  }

  createLogHeader(OtherLogList, false);
  createLogHeader(PrimaryLogList, true);
}

addLogLinks();
