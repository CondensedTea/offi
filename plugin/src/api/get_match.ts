import {Match, MatchResponse} from "./types";

export async function getMatch(apiBaseUrl: string, matchId: number): Promise<Match> {
  const logURL = new URL(apiBaseUrl + "/log/" + matchId.toString());

  const res = await fetch(logURL.toString());

  if (!res.ok) {
    throw new Error("api returned error: " + res.statusText);
  }

  const apiResponse = await res.json() as MatchResponse;

  return apiResponse.match;
}
