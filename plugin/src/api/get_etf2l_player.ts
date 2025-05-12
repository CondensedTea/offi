import { Player, PlayersResponse } from "./types";
import { requestHeaders } from "./api";
import { APIError } from "./error";

export async function getEtf2lPlayer(apiBaseUrl: string, steamID: string, withRecruitmentStatus: boolean = false): Promise<Player> {
  const playersURL = new URL(apiBaseUrl + "/api/v1/etf2l/players");

  playersURL.searchParams.set("id", steamID.toString());

  if (withRecruitmentStatus) playersURL.searchParams.append("with_recruitment_status", "true");

  const res = await fetch(playersURL.toString(), {
    headers: requestHeaders,
  });
  if (res.status !== 200) {
    throw await APIError.fromResponse(res);
  }

  const response = await res.json() as PlayersResponse;
  return response.players[0];
}
