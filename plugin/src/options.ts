import { InterfaceSettings } from "@kocal/web-extension-library";

export default {
  linkLogsOnMatchpages: {
    type: "boolean",
    label: "Show logs on ETF2L match pages",
    defaultValue: "true",
  },
  linkMatchepagesOnLogs: {
    type: "boolean",
    label: "Link etf2l matchpages on logs.tf",
    defaultValue: "true",
  },
  replaceNamesInLogs: {
    type: "boolean",
    label: "Replace player names on logs.tf with etf2l names",
    defaultValue: "true",
  },
  showBansForPlayers: {
    type: "boolean",
    label: "Show bans on player pages",
    defaultValue: "true",
  },
  showLftForPlayer: {
    type: "boolean",
    label: "Link recruitments on player pages",
    defaultValue: "true",
  },
  addLinksOnSteamProfiles: {
    type: "boolean",
    label: "Add links to etf2l and logs.tf to Steam profiles",
    defaultValue: "true",
  },
  apiBaseURL: {
    type: "string",
    label: "Offi API URL",
    defaultValue: "https://offi.lemontea.dev",
  },
} as InterfaceSettings;

