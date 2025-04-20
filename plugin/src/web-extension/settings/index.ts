import { readFromSyncStorage, writeToSyncStorage } from './sync-storage';

type ConfigValue = string | boolean

/* tslint:disable variable-name */
let _settings: InterfaceSettings = {};

export interface InterfaceSettings {
  [k: string]: InterfaceSettingsItem;
}

export interface InterfaceSettingsItem {
  type: string;
  label: string;
  help?: string;
  children?: InterfaceSettings;
  // for any type
  value?: ConfigValue;
  defaultValue?: ConfigValue;
  // when type is `range'
  min?: number;
  max?: number;
}

export const getSettings = (): InterfaceSettings => _settings;

export const loadSettings = async (): Promise<void> => {
  const items = await readFromSyncStorage('settings');
  const settings: InterfaceSettings | undefined = items.settings;

  Object.entries(settings || getFlattenSettings()).forEach(([dottedName, value]) => {
    setSettingValue(dottedName, value, false);
  });

  return writeToSyncStorage({ settings: getFlattenSettings() });
};

/**
 * Should be called only one time in your whole app
 */
export const registerSettings = (settings: InterfaceSettings): Promise<void> => {
  _settings = settings;
  return loadSettings();
};

export function getSettingValue<T = ConfigValue>(key: string): T | undefined {
  const setting = getSetting(key);

  if (setting === undefined) return undefined;
  if (setting.value === undefined) return setting.defaultValue as T;

  return setting.value as T;
}

export async function setSettingValue<T = ConfigValue>(key: string, value: T, synchronize = true): Promise<void> {
  const setting = getSetting(key);
  const previousValue = getSettingValue(key);

  if (setting === undefined) return Promise.resolve();

  Object.assign(setting, { value });

  if (!synchronize) {
    return Promise.resolve();
  }

  try {
    await writeToSyncStorage({ settings: getFlattenSettings() });
    return Promise.resolve();
  } catch {
    Object.assign(setting, { value: previousValue });
    return Promise.reject();
  }
}

function getSetting(key: string): InterfaceSettingsItem | undefined {
  const keyParts = key.split('.');
  let setting = getSettings()[keyParts[0]];

  if (keyParts.length === 1) return setting;

  while (true) {
    keyParts.shift();

    if (keyParts.length === 0) return setting;
    if (setting === undefined || setting.children === undefined) return undefined;

    setting = setting.children[keyParts[0]];
  }
}

type Callback = (key: string, setting: InterfaceSettingsItem) => void;

function recursivelyIterate(settings: InterfaceSettings, cb: Callback, prefix: string = ''): void {
  Object.entries(settings).forEach(([name, setting]) => {
    if (setting.children) {
      recursivelyIterate(setting.children, cb, `${name}.`);
    }

    cb(`${prefix}${name}`, setting);
  });
}

function getFlattenSettings(): { [k: string]: ConfigValue } {
  const flattenSettings: { [k: string]: ConfigValue } = {};

  recursivelyIterate(getSettings(), dottedName => {
    flattenSettings[dottedName] = getSettingValue(dottedName);
  });

  return flattenSettings;
}
