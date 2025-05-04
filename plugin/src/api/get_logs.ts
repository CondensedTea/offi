import { Log, LogResponse } from "./types";
import { APIError, NoLogsError } from "./error";
import { requestHeaders } from "./api";

export async function getLogs(apiBaseUrl: string, matchId: number): Promise<Log[]> {
  const getMatchURL = new URL(apiBaseUrl + `/match/${matchId}`);

  const res = await fetch(getMatchURL, {
    headers: requestHeaders,
  });
  if (res.status === 404 || res.status === 425) {
    throw NoLogsError;
  } else if (res.status !== 200) {
    throw await APIError.fromResponse(res);
  }

  const { logs } = (await res.json()) as LogResponse;

  return logs.map((log => new Log(log)));
}
