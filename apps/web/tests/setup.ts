import '@testing-library/jest-dom/vitest';
import { fetch, Headers, Request, Response } from 'undici';

globalThis.fetch = fetch;
globalThis.Headers = Headers;
globalThis.Request = Request;
globalThis.Response = Response;