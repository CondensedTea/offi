import {Match, MatchResponse} from "./types";

export async function getMatch(apiBaseUrl: string, logId: number): Promise<Match> {
  const logURL = new URL(apiBaseUrl + `/log/${logId}`);

  const res = await fetch(logURL);

  if (!res.ok) {
    throw new Error("api returned error: " + res.statusText);
  }

  const apiResponse = await res.json() as MatchResponse;

  return apiResponse.match;
}
