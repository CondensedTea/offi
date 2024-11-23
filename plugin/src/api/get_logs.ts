import { Log, LogResponse } from "./types";

export const NoLogsError = new Error("api didnt found logs for this match");

export async function getLogs(apiBaseUrl: string, matchId: number): Promise<Log[]> {
  const getMatchURL = new URL(apiBaseUrl + "/match/" + matchId.toString());

  const res = await fetch(getMatchURL.toString());

  if (!res.ok) {
    throw new Error("offi api returned error: " + res.statusText);
  }

  const { logs } = (await res.json()) as LogResponse;
  if (logs === null) {
    throw NoLogsError;
  }

  const parsedLogs: Log[] = [];

  for (const rawLog of logs) {
    parsedLogs.push(new Log(rawLog));
  }

  return parsedLogs;
}
