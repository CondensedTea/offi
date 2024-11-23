import { Player, PlayersResponse } from "./types";

export async function getPlayers(apiBaseUrl: string, ids: string[]): Promise<Player[]> {
  const playersURL = new URL(apiBaseUrl + "/players");

  const idsString = ids.join(",");

  playersURL.searchParams.append("id", idsString);

  const res = await fetch(playersURL.toString());
  if (!res.ok) {
    throw new Error("offi api returned error: " + res.statusText);
  }

  const response = await res.json() as PlayersResponse;
  return response.players;
}
