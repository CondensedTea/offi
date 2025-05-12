import { registerSettings, setSettingValue, getSettings } from "./web-extension/settings";
import settings from "./options";

document.addEventListener("DOMContentLoaded", async () => {
  const optionsContainer = document.getElementById("options");
  if (!optionsContainer) {
    return;
  }

  await registerSettings(settings);

  const settingsValues = getSettings();

  for (const [key, item] of Object.entries(settingsValues)) {
    const labelEl = document.createElement("label");
    const inputEl = document.createElement("input");
    const selectEl = document.createElement("select");
    const textEl = document.createElement("div");
    textEl.className = "item-label";

    switch (item.type) {
      case "boolean":
        textEl.textContent = item.label;

        inputEl.type = "checkbox";
        inputEl.checked = item.value as boolean;
        inputEl.addEventListener("change", async () => {
          await setSettingValue(key, inputEl.checked);
        });
        break;
      case "string":
        textEl.textContent = item.label;

        inputEl.type = "url";
        inputEl.value = item.value as string;
        inputEl.addEventListener("change", async () => {
          await setSettingValue(key, inputEl.value);
        });
        break;
      case "enum":
        textEl.textContent = item.label;

        item.enumValues?.forEach((value) => {
          const optionEl = document.createElement("option");
          optionEl.value = value;
          optionEl.textContent = value;
          selectEl.appendChild(optionEl);
        });
        selectEl.value = item.value as string;
        selectEl.addEventListener("change", async () => {
          await setSettingValue(key, selectEl.value);
        });
        break;
    }

    if (item.type === "enum") {
      labelEl.appendChild(selectEl);
    } else {
      labelEl.appendChild(inputEl);
    }

    labelEl.appendChild(textEl);
    optionsContainer.appendChild(labelEl);
  }
});
