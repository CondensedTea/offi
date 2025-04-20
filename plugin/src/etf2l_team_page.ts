import { Team } from "./api/types";
import { getSettingValue } from "./web-extension/settings";
import { getTeam } from "./api/get_team";
import { NoRecruitmentInfo } from "./api/error"

const playerRe = RegExp("https://etf2l.org/teams/(\\d+)/");

function getTeamID(): string {
  const match = document.URL.match(playerRe);

  if (match === null || match.length < 1) {
    throw new Error("could not find team ID");
  }
  return match[1];
}

export async function addTeamStatus() {
  const apiBaseUrl = getSettingValue<string>("apiBaseURL") as string;
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

  const row = document.createElement("tr");

  const headerCell = document.createElement("td");
  headerCell.innerText = "LFP";

  const linkCell = document.createElement("td");
  const link = document.createElement("a");

  link.href = teamStatus.recruitment.url;
  link.innerText = classesString;
  linkCell.appendChild(link);

  row.append(headerCell, linkCell);

  document.querySelector(".teaminfo > tbody").appendChild(row);
}
