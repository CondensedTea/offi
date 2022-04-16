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
  let match;

  try {
    match = await getMatchFromAPI(matchId);
  } catch (e) {
    console.log('could not get match: ' + e.toString());
    return;
  }

  const matchBlock = document.createElement('h3');

  matchBlock.innerHTML =
    `<a href="https://etf2l.org/matches/${match.id}"> ${match.competition} | ${match.stage} </a>`;

  const logDateElem = document.getElementById('log-date');

  logDateElem.after(matchBlock);
}

addMatchLink();
