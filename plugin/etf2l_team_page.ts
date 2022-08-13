import {api, apiUrl, type} from "./utils";
import {TeamResponse, Team} from "./types";

const playerRe = RegExp("https://etf2l.org/teams/(\\d+)/");
const NoRecruitmentInfo = new Error("this team doesn't have recruitment post");

function getTeamID(): number {
  const match = document.URL.match(playerRe);

  if (match === null || match.length < 1) {
    throw new Error("could not find team ID");
  }
  return parseInt(match[1]);
}

async function getTeamStatusFromAPI(teamId: number): Promise<Team> {
  const getTeamURL = new URL(apiUrl + "team/" + teamId.toString());
  getTeamURL.searchParams.append("version", api.runtime.getManifest().version);
  getTeamURL.searchParams.append("browser", type);

  const res = await fetch(getTeamURL.toString());
  if (!res.ok) {
    throw new Error("offi api returned error: " + res.statusText);
  }

  const response = (await res.json()) as TeamResponse;
  if (response.team === null || response.team.recruitment.empty) {
    throw NoRecruitmentInfo;
  }

  return response.team;
}

async function addTeamStatus() {
  const playerId = getTeamID();

  let teamStatus: Team;

  try {
    teamStatus = await getTeamStatusFromAPI(playerId);
  } catch (e) {
    if (e === NoRecruitmentInfo) {
      return;
    } else {
      console.error("failed to get team status: ", e.toString());
    }
  }

  let classesString: string;

  if (teamStatus.recruitment.classes.length > 3) {
    classesString = "3+ classes";
  } else {
    classesString = teamStatus.recruitment.classes.join(", ");
  }

  const node = document.createElement("tr");
  node.innerHTML = `
        <td>LFP</td>
        <td><a href=${teamStatus.recruitment.url}>${classesString}</a></td>`;
  document.querySelector(".teaminfo > tbody").appendChild(node);
}

api.storage.sync.get((fields) => {
  if (fields.etf2l_show_lfp === true) {
    addTeamStatus();
  }
});
