import { TextDecoder, TextEncoder } from 'node:util';
import * as undici from 'undici';

if (!globalThis.TextEncoder) {
  globalThis.TextEncoder = TextEncoder;
}

if (!globalThis.TextDecoder) {
  globalThis.TextDecoder = TextDecoder;
}

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
