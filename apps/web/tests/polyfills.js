const { TextDecoder, TextEncoder } = require('node:util');

if (!globalThis.TextEncoder) {
  globalThis.TextEncoder = TextEncoder;
}

if (!globalThis.TextDecoder) {
  globalThis.TextDecoder = TextDecoder;
}

const undici = require('undici');

if (!globalThis.fetch) {
  globalThis.fetch = undici.fetch;
}

if (!globalThis.Headers) {
  globalThis.Headers = undici.Headers;
}

if (!globalThis.Request) {
  globalThis.Request = undici.Request;
}

if (!globalThis.Response) {
  globalThis.Response = undici.Response;
}
