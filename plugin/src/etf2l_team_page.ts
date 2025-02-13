import { Team } from "./api/types";
import { getSettingValue } from "./web-extension/settings";
import { getTeam, NoRecruitmentInfo } from "./api/get_team";

const playerRe = RegExp("https://etf2l.org/teams/(\\d+)/");

function getTeamID(): string {
  const match = document.URL.match(playerRe);

  if (match === null || match.length < 1) {
    throw new Error("could not find team ID");
  }
  return match[1];
}

export async function addTeamStatus() {
  const apiBaseUrl = await getSettingValue("apiBaseURL");
  const playerId = getTeamID();

  let teamStatus: Team;

  try {
    teamStatus = await getTeam(apiBaseUrl, playerId);
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
