import '@testing-library/jest-dom';

// Vitest's jsdom environment does not ship the WHATWG fetch classes in Node.
// Polyfill them using undici so libraries like openapi-fetch can construct Requests.
import { fetch as undiciFetch, Headers, Request, Response } from 'undici';
import { TextDecoder, TextEncoder } from 'node:util';

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

if (!globalThis.TextEncoder) {
  globalThis.TextEncoder = TextEncoder;
}

if (!globalThis.TextDecoder) {
  globalThis.TextDecoder = TextDecoder;
}

