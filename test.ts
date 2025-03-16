#!/usr/bin/env -S deno run
import {Buffer} from "node:buffer"

function decodeString(hexString) {
  return Buffer.from(hexString, 'hex');
}

// Encode a byte array to a hex string
function encodeToString(byteArray) {
  return Buffer.from(byteArray).toString('hex');
}

const secret: string = "425111bc21a56c0e8e55ceb856685a512017c510a769146d98692149456b2dd3"

const messageBuffer = decodeString(secret)

const hashBuffer = await crypto.subtle.digest("SHA-256", messageBuffer);

const hash = encodeToString(hashBuffer)

console.log(hash)

