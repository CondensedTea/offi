export type Recruitment = {
  skill: string;
  url: string;
  game_mode: string;
  classes: string[];
}

export type Ban = {
  start: number
  startDate: Date;
  end: number;
  endDate: Date;
  reason: string;
}

export type Player = {
  id: string;
  steam_id: string;
  name: string;
  bans: Ban[];
  recruitment: Recruitment;
}

export type Team = {
  ID: number;
  recruitment: Recruitment;
}



export class Log {
  id: number;
  title: string;
  map: string;
  played_at: Date;
  is_secondary: boolean;

  constructor(data: Log) {
    this.id = data.id;
    this.title = data.title;
    this.map = data.map;
    this.played_at = new Date(data.played_at);
    this.is_secondary = data.is_secondary;
  }
}

export type Match = {
  match_id: number;
  competition: string;
  stage: string;
  tier: string;
}

export type MatchResponse = {
  match: Match;
}

export type LogResponse = {
  logs: Log[];
};

export type PlayersResponse = {
  players: Player[];
}

export type TeamResponse = {
  team: Team;
}
