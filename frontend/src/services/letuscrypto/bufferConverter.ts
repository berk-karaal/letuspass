/**
 * @fileOverview Provides functions to convert strings and hex strings to ArrayBuffers
 * and vice versa. These functions are useful for serializing and deserializing
 * data for cryptographic operations.
 */

export function stringToArrayBuffer(str: string): ArrayBuffer {
  const encoder = new TextEncoder();
  return encoder.encode(str).buffer;
}

export function arrayBufferToString(buffer: ArrayBuffer): string {
  const decoder = new TextDecoder();
  return decoder.decode(buffer);
}

export function hexToArrayBuffer(hex: string): ArrayBuffer {
  if (hex.length % 2 !== 0) {
    throw new Error("Invalid hex string length. hex: " + hex);
  }

  // Create a Uint8Array to hold the bytes
  const byteArray = new Uint8Array(hex.length / 2);

  // Convert each pair of hex characters to a byte
  for (let i = 0; i < hex.length; i += 2) {
    byteArray[i / 2] = parseInt(hex.substring(i, i + 2), 16);
  }

  return byteArray.buffer;
}

export function arrayBufferToHex(buffer: ArrayBuffer): string {
  return Array.from(new Uint8Array(buffer))
    .map((byte) => byte.toString(16).padStart(2, "0"))
    .join("");
}
