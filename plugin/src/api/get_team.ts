import { Team, TeamResponse } from "./types";
import { NoRecruitmentInfo, APIError } from "./error";
import {requestHeaders} from "./api";

export async function getTeam(apiBaseUrl: string, teamId: string): Promise<Team> {
  const getTeamURL = new URL(apiBaseUrl + `/team/${teamId}`);

  const res = await fetch(getTeamURL, {
    headers: requestHeaders,
  });
  if (res.status === 404) {
    throw NoRecruitmentInfo;
  } else if (res.status !== 200) {
    throw await APIError.fromResponse(res);
  }

  const response = (await res.json()) as TeamResponse;

  return response.team;
}
