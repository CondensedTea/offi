import {Team, TeamResponse} from "./types";

export const NoRecruitmentInfo = new Error("this team doesn't have recruitment post");

export async function getTeam(apiBaseUrl: string, teamId: string): Promise<Team> {
  const getTeamURL = new URL(apiBaseUrl + `/team/${teamId}`);

  const res = await fetch(getTeamURL);
  if (!res.ok) {
    throw new Error("offi api returned error: " + res.statusText);
  }

  const response = (await res.json()) as TeamResponse;
  if (response.team === null || !response.team.recruitment) {
    throw NoRecruitmentInfo;
  }

  return response.team;
}
