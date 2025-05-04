import { Player, PlayersResponse } from "./types";
import { APIError } from "./error";
import { requestHeaders } from "./api";

export async function getPlayers(apiBaseUrl: string, ids: string[], withRecruitmentStatus: boolean = false): Promise<Player[]> {
  const playersURL = new URL(apiBaseUrl + "/players");

  const idsString = ids.join(",");

  playersURL.searchParams.append("id", idsString);

  if (withRecruitmentStatus) playersURL.searchParams.append("with_recruitment_status", "true");

  const res = await fetch(playersURL.toString(), {
    headers: requestHeaders,
  });
  if (res.status !== 200) {
    throw await APIError.fromResponse(res);
  }

  const response = await res.json() as PlayersResponse;
  return response.players;
}
