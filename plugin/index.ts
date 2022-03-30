import 'regenerator-runtime/runtime';

const matchRe = RegExp('https://etf2l.org/matches/(\\d+)/');
const apiUrl = 'https://stg.lemontea.dev/offi/match/';

class Log {
  id: number;
  map: string;
  played_at: Date;
  constructor(data: Object) {
    this.id = data['id'];
    this.map = data['map'];
    this.played_at = new Date(data['played_at']);
  }
}

function getMatchID(): number {
  const match = document.URL.match(matchRe);

  if (match === null || match.length < 1) {
    throw new Error('could not find match ID');
  }
  return parseInt(match[1]);
}

async function getLogsFromAPI(matchId: number) {
  const res = await fetch(apiUrl + matchId.toString());
  const apiResponse = await res.json();
  const logs: Log[] = [];

  for (const logData of apiResponse['logs']) {
    const l = new Log(logData);
    logs.push(l);
  }
  return logs;
}

async function addLogLinks(): Promise<void> {
  const matchId = getMatchID();
  const logs = await getLogsFromAPI(matchId);

  const LogList = document.createElement('ul');

  for (const log of logs) {
    const logItem = document.createElement('li');
    logItem.innerHTML =`<a href="https://logs.tf/${log.id}"> #${log.id} </a> | ${log.map} | ${log.played_at.toLocaleString()}`;
    LogList.appendChild(logItem);
  }
  const LogHeader = document.createElement('div');
  LogHeader.className = 'offi';

  if (logs.length === 1) {
    LogHeader.innerHTML = `<h2>1 Log</h2>`;
  } else {
    LogHeader.innerHTML = `<h2>${logs.length} Logs</h2>`;
  }
  LogHeader.append(LogList);

  const playersSection = document.getElementsByClassName('fix match-players');
  if (playersSection === null || playersSection.length < 1) {
    return;
  }
  playersSection[0].after(LogHeader);
  return;
}

addLogLinks();
