import '@testing-library/jest-dom/vitest';

// Vitest's jsdom environment does not ship the WHATWG fetch classes in Node.
// Polyfill them using undici so libraries like openapi-fetch can construct Requests.
import { fetch as undiciFetch, Headers, Request, Response } from 'undici';

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