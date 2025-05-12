import { InterfaceSettings } from "./web-extension/settings";

export default {
  linkLogsOnMatchpages: {
    type: "boolean",
    label: "Show logs on ETF2L match pages",
    defaultValue: "true",
  },
  linkMatchepagesOnLogs: {
    type: "boolean",
    label: "Link ETF2L matches on logs.tf",
    defaultValue: "true",
  },
  replaceNamesInLogs: {
    type: "boolean",
    label: "Replace player names on logs.tf with league names",
    defaultValue: "true",
  },
  showBansForPlayers: {
    type: "boolean",
    label: "Show bans on ETF2L player pages",
    defaultValue: "true",
  },
  showLftForPlayer: {
    type: "boolean",
    label: "Link recruitments on ETF2L player pages",
    defaultValue: "true",
  },
  showLfpForTeam: {
    type: "boolean",
    label: "Link recruitments on ETF2L team pages",
    defaultValue: "true",
  },
  addLinksOnSteamProfiles: {
    type: "boolean",
    label: "Add links to league profiles and logs.tf to Steam profiles",
    defaultValue: "true",
  },
  league: {
    type: "enum",
    enumValues: ["ETF2L", "RGL"],
    label: "Player names' source",
    defaultValue: "ETF2L",
  },
  apiBaseURL: {
    type: "string",
    label: "Offi API URL",
    defaultValue: "https://offi.lemontea.dev",
  },
} as InterfaceSettings;
