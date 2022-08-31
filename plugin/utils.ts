export const apiUrl = "https://offi.lemontea.dev/";

class Options {
  logstf_link_matchpage = true;
  logstf_replace_names = true;

  etf2l_show_bans = true;
  etf2l_show_lft = true;
  etf2l_show_lfp = true;
  etf2l_show_logs = true;
}

export let api;
export let type: string;

if (typeof chrome !== "undefined" && typeof chrome.runtime !== "undefined") {
  api = chrome;
  type = "chrome";
} else {
  api = browser;
  type = "firefox";
}

api.storage.sync.get((fields: Object) => {
  if (fields["logstf_link_matchpage"] === undefined) {
    const fields = Object(new Options());
    api.storage.sync.set(fields);
  }
});

export function replaceInText(element, pattern, replacement) {
  for (const node of element.childNodes) {
    switch (node.nodeType) {
      case Node.ELEMENT_NODE || Node.DOCUMENT_NODE:
        replaceInText(node, pattern, replacement);
        break;
      case Node.TEXT_NODE:
        node.textContent = node.textContent.replace(pattern, replacement);
        break;
    }
  }
}
