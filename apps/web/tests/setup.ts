import '@testing-library/jest-dom';
import { TextDecoder, TextEncoder } from 'node:util';

// Ensure WHATWG text encoders exist before loading undici, which expects them.
if (!globalThis.TextEncoder) {
  globalThis.TextEncoder = TextEncoder;
}

if (!globalThis.TextDecoder) {
  globalThis.TextDecoder = TextDecoder;
}

// Vitest's jsdom environment does not ship the WHATWG fetch classes in Node.
// Load undici *after* TextEncoder/TextDecoder are available so it can polyfill.
const { fetch: undiciFetch, Headers, Request, Response } = await import('undici');

if (!globalThis.fetch) {
  globalThis.fetch = undiciFetch;
}

if (!globalThis.Headers) {
  globalThis.Headers = Headers;
}

if (!globalThis.Request) {
  globalThis.Request = Request;
}

if (!globalThis.Response) {
  globalThis.Response = Response;
}


