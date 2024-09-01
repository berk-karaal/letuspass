/**
 * @fileOverview Provides functions to encrypt and decrypt data using AES-GCM.
 * Note: The words "aes" and "aes-gcm" are used interchangeably in this file.
 * Since this project only uses aes-gcm encryption, the term "aes" is used to
 * refer to aes-gcm encryption.
 */

import {
  arrayBufferToHex,
  arrayBufferToString,
  hexToArrayBuffer,
  stringToArrayBuffer,
} from "./bufferConverter";
import { generateRandomString } from "./randomKey";

export function generateRandomIV(): string {
  return generateRandomString(12);
}

/**
 * Internal function to convert hex encoded aes key to CryptoKey
 * @param hex hex encoded AES-GCM key
 */
function _hexKeyToCryptoKey(hex: string): Promise<CryptoKey> {
  return window.crypto.subtle.importKey(
    "raw",
    hexToArrayBuffer(hex),
    "AES-GCM",
    true,
    ["encrypt", "decrypt"]
  );
}

/**
 * Decrypt data using AES-GCM
 * @param key hex encoded aes key
 * @param iv hex encoded Initialization Vector
 * @param data hex encoded data to be decrypted
 * @returns utf-16 (default string encoding) encoded decrypted data
 */
export async function decrypt(
  key: string,
  iv: string,
  data: string
): Promise<string> {
  const cryptoKey = await _hexKeyToCryptoKey(key);
  let decrypted: ArrayBuffer;
  try {
    decrypted = await window.crypto.subtle.decrypt(
      { name: "AES-GCM", iv: hexToArrayBuffer(iv) },
      cryptoKey,
      hexToArrayBuffer(data)
    );
  } catch (e) {
    console.error("AES-GCM decrypt error:", e);
    throw e;
  }
  return arrayBufferToString(decrypted);
}

/**
 * Encrypt data using AES-GCM
 * @param key hex encoded aes key
 * @param iv hex encoded Initialization Vector
 * @param data utf-16 (default string encoding) encoded data to be encrypted
 * @returns hex encoded encrypted data
 */
export async function encrypt(
  key: string,
  iv: string,
  data: string
): Promise<string> {
  const cryptoKey = await _hexKeyToCryptoKey(key);
  const encrypted = await window.crypto.subtle.encrypt(
    { name: "AES-GCM", iv: hexToArrayBuffer(iv) },
    cryptoKey,
    stringToArrayBuffer(data)
  );
  return arrayBufferToHex(new Uint8Array(encrypted));
}
