let api;

if (typeof chrome !== "undefined" && typeof chrome.runtime !== "undefined") {
  api = chrome;
} else {
  api = browser;
}

class Options {
  logstf_link_matchpage = true;
  logstf_replace_names = true;

  etf2l_show_bans = true;
  etf2l_show_lft = true;
  etf2l_show_lfp = true;
  etf2l_show_logs = true;
}

function saveOptions() {
  const opts = new Options();
  const fields = Object.keys(opts);

  fields.forEach((optionName) => {
    opts[optionName] = (document.getElementById(optionName) as HTMLInputElement).checked;
  });

  const {...object} = opts;
  chrome.storage.sync.set(object, () => {
    const status = document.getElementById("status");
    status.textContent = "Options saved.";
    setTimeout(() => {
      status.textContent = "";
    }, 750);
  });
}

function restoreOptions() {
  api.storage.sync.get((fields) => {
    console.log("restoreOptions ", fields);
    Object.entries(fields).forEach((value: [string, boolean]) => {
      if (value[0] !== "defaults_loaded") {
        const node = document.getElementById(value[0]) as HTMLInputElement;
        node.checked = value[1];
      }
    });
  });
}

function restoreDefaults() {
  const fields = Object(new Options());
  console.log("new fields", fields);
  api.storage.sync.set(fields, () => {
    const status = document.getElementById("status");
    status.textContent = "Options reset.";
    setTimeout(() => {
      status.textContent = "";
    }, 750);
  });
}

document.addEventListener("DOMContentLoaded", restoreOptions);
document.getElementById("save").addEventListener("click", saveOptions);
document.getElementById("reset").addEventListener("click", restoreDefaults);
