import 'regenerator-runtime/runtime';

const matchRe = RegExp('https://logs.tf/(\\d+)');
const apiUrl = 'https://offi.lemontea.dev/log/';

class Match {
  id: number;
  competition: string;
  stage: string;
  constructor(data: Object) {
    this.id = data['match_id'];
    this.competition = data['competition'];
    this.stage = data['stage'];
  }
}

function getLogID(): number {
  const match = document.URL.match(matchRe);

  if (match === null || match.length < 1) {
    throw new Error('could not find log ID');
  }
  return parseInt(match[1]);
}

async function getMatchFromAPI(matchId: number): Promise<Match> {
  const res = await fetch(apiUrl + matchId.toString());

  if (!res.ok) {
    throw new Error('offi api returned error: ' + res.statusText);
  }

  const apiResponse = await res.json();

  return new Match(apiResponse['match']);
}

async function addMatchLink(): Promise<void> {
  const matchId = getLogID();
  let match: Match;

  try {
    match = await getMatchFromAPI(matchId);
  } catch (e) {
    console.log('could not get match: ' + e.toString());
    return;
  }

  const competitionBlock = document.createElement('h3');
  competitionBlock.innerHTML =
    `<a href="https://etf2l.org/matches/${match.id}">${match.competition}</a>`;

  const matchBlock = document.createElement('h3');
  matchBlock.innerText = match.stage;

  const logDateElem = document.getElementById('log-date');

  logDateElem.after(matchBlock);
  logDateElem.after(competitionBlock);
}

addMatchLink();
