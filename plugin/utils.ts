export const apiUrl = "http://offi.lemontea.dev/";

export function api() {
  if (typeof chrome !== "undefined" && typeof chrome.runtime !== "undefined") {
    return chrome;
  }
  return browser;
}
