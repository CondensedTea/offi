import { Player, PlayersResponse } from "./types";

export async function getPlayers(apiBaseUrl: string, ids: string[], withRecruitmentStatus: boolean): Promise<Player[]> {
  const playersURL = new URL(apiBaseUrl + "/players");

  const idsString = ids.join(",");

  playersURL.searchParams.append("id", idsString);

  if (withRecruitmentStatus) playersURL.searchParams.append("with_recruitment_status", "true");

  const res = await fetch(playersURL.toString());
  if (!res.ok) {
    throw new Error("offi api returned error: " + res.statusText);
  }

  const response = await res.json() as PlayersResponse;
  return response.players;
}
