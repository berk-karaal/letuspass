import { arrayBufferToHex } from "./bufferConverter";

/**
 * Generate a random string of the specified length.
 *
 * String is generated using crypto API and encoded to hex. Returned string
 * can be used as IV or salt.
 */
export function generateRandomString(length: number): string {
  const randomValues = window.crypto.getRandomValues(new Uint8Array(length));
  return arrayBufferToHex(randomValues.buffer);
}
