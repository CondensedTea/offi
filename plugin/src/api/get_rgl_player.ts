import { Player, PlayersResponse } from "./types";
import { requestHeaders } from "./api";
import { APIError } from "./error";

export async function getRGLPlayer(apiBaseUrl: string, steamID: string): Promise<Player> {
  const playersURL = new URL(apiBaseUrl + "/api/v1/rgl/players");

  playersURL.searchParams.set("id", steamID.toString());

  const res = await fetch(playersURL.toString(), {
    headers: requestHeaders,
  });
  if (res.status !== 200) {
    throw await APIError.fromResponse(res);
  }

  const response = await res.json() as PlayersResponse;
  return response.players[0];
}
