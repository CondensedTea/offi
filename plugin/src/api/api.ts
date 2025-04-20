const VERSION =
  (process.env.OFFI_PLATFORM || "unknown") + "/" + (process.env.OFFI_VERSION || "dev");

export const requestHeaders = {
    "X-Offi-Version": VERSION,
    "Content-Type": "application/json",
}

