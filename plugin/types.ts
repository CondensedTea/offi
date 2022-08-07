export type Recruitment = {
  id: number;
  skill: string;
  url: string;
  game_mode: string;
  classes: string[];
  empty: boolean;
}

export type Ban = {
  start: number;
  end: number;
  reason: string;
}

export type Player = {
  id: number;
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
  id: number;
  competition: string;
  stage: string;
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
