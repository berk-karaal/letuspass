/**
 * @fileOverview Provides functions to generate Elliptic Curve key pairs
 * and get shared key between two key pairs.
 */

import { ec as EC } from "elliptic";
import { stringToArrayBuffer } from "./bufferConverter";
import { generateRandomString } from "./randomKey";

export const EC_CURVE = "curve25519";

export function generateRandomSalt() {
  return generateRandomString(32);
}

/**
 * Internal function to derive a strengthend key from a password with a salt
 * using PBKDF2.
 * @param password Password to derive key from
 * @param salt Salt to use in the derivation process.
 */
async function _deriveKeyFromPassword(
  password: string,
  salt: string
): Promise<CryptoKey> {
  const keyMaterial = await crypto.subtle.importKey(
    "raw",
    stringToArrayBuffer(password),
    { name: "PBKDF2" },
    false,
    ["deriveKey"]
  );

  return crypto.subtle.deriveKey(
    {
      name: "PBKDF2",
      salt: stringToArrayBuffer(salt),
      iterations: 100000,
      hash: "SHA-256",
    },
    keyMaterial,
    { name: "HMAC", hash: "SHA-256", length: 256 },
    true,
    ["sign", "verify"]
  );
}

/**
 * Returns a deterministically generated Elliptic Curve key pair from a password and salt.
 */
export async function generateECKeyPairFromPassword(
  password: string,
  salt: string
): Promise<{
  privateKey: string;
  publicKey: string;
}> {
  const derivedKey = await _deriveKeyFromPassword(password, salt);
  const seed = await crypto.subtle.exportKey("raw", derivedKey);

  const ECKeyPair = new EC(EC_CURVE).keyFromPrivate(new Uint8Array(seed));
  return {
    privateKey: ECKeyPair.getPrivate("hex"),
    publicKey: ECKeyPair.getPublic("hex"),
  };
}

/**
 * Returns the shared key between two Elliptic Curve key pairs.
 */
export function getSharedKey(privateKey: string, publicKey: string) {
  const ec = new EC(EC_CURVE);
  const userKeyPair = ec.keyFromPrivate(privateKey, "hex");
  const otherUserKeyPair = ec.keyFromPublic(publicKey, "hex");
  const sharedSecret = userKeyPair
    .derive(otherUserKeyPair.getPublic())
    .toString("hex");
  return sharedSecret;
}
