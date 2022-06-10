import {apiUrl} from "./utils";

const playerRe = RegExp("https://etf2l.org/teams/(\\d+)/");
const NoRecruitmentInfo = new Error("this team doesn't have recruitment post");

class TeamStatus {
  id: number;
  skill: string;
  url: string;
  game_mode: string;
  classes: string[];
  empty: boolean;
}

type ApiResponse = {
  status: TeamStatus;
};

function getTeamID(): number {
  const match = document.URL.match(playerRe);

  if (match === null || match.length < 1) {
    throw new Error("could not find team ID");
  }
  return parseInt(match[1]);
}

async function getTeamStatusFromAPI(teamId: number): Promise<TeamStatus> {
  const getTeamURL = new URL(apiUrl + "team/" + teamId.toString());
  getTeamURL.searchParams.append("version", chrome.runtime.getManifest().version);

  const res = await fetch(getTeamURL.toString());

  if (!res.ok) {
    throw new Error("offi api returned error: " + res.statusText);
  }

  const teamResponse = (await res.json()) as ApiResponse;
  if (teamResponse.status === null || teamResponse.status.empty) {
    throw NoRecruitmentInfo;
  }

  return teamResponse.status;
}

async function addTeamStatus() {
  const playerId = getTeamID();

  let teamStatus: TeamStatus;

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

  if (teamStatus.classes.length > 3) {
    classesString = "3+ classes";
  } else {
    classesString = teamStatus.classes.join(", ");
  }

  const node = document.createElement("tr");
  node.innerHTML = `
        <td>LFP</td>
        <td><a href=${teamStatus.url}>${classesString}</a></td>`;
  document.querySelector(".teaminfo > tbody").appendChild(node);
}

addTeamStatus();
