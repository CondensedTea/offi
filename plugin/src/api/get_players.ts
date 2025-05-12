import { Player, PlayersResponse } from "./types";
import { APIError } from "./error";
import { requestHeaders } from "./api";

export async function getPlayers(apiBaseUrl: string, league: string, steamIDs: string[]): Promise<Player[]> {
  if (!["etf2l", "rgl"].includes(league)) throw new Error("Invalid league: " + league);

  const playersURL = new URL(apiBaseUrl + `/api/v1/${league}/players`);

  playersURL.searchParams.append("id", steamIDs.join(","));

  const res = await fetch(playersURL.toString(), {
    headers: requestHeaders,
  });
  if (res.status !== 200) {
    throw await APIError.fromResponse(res);
  }

  const response = await res.json() as PlayersResponse;
  return response.players;
}
