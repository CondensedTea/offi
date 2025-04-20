import { MatchResponse } from "./types";
import { MatchNotFound, APIError } from "./error";
import { requestHeaders } from "./api";

export async function getMatch(apiBaseUrl: string, logId: number): Promise<MatchResponse> {
  const logURL = new URL(apiBaseUrl + `/log/${logId}`);

  const res = await fetch(logURL, {
    headers: requestHeaders,
  });
  if (res.status === 404) {
    throw MatchNotFound
  }
  if (res.status !== 200) {
    throw await APIError.fromResponse(res);
  }

  return await res.json() as MatchResponse;
}
