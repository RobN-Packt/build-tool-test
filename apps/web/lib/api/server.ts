import type { components } from './types';

type Book = components['schemas']['Book'];

interface ListBooksResponse {
  books: Book[];
}

interface BooksHealthResponse {
  status: string;
  checkedAt: string;
}

interface BackendDiagnostics {
  url: string;
  status: number;
  ok: boolean;
  body: string;
  fetchedAt: string;
  error?: string;
}

const ERROR_PREVIEW_LIMIT = 512;

function getBackendBaseUrl() {
  const base = process.env.BACKEND_INTERNAL_URL;
  if (!base) {
    throw new Error('BACKEND_INTERNAL_URL is not configured');
  }
  return base.endsWith('/') ? base.slice(0, -1) : base;
}

function previewBody(body: string) {
  if (body.length <= ERROR_PREVIEW_LIMIT) {
    return body;
  }
  return `${body.slice(0, ERROR_PREVIEW_LIMIT)}…`;
}

function logBackendFailure(message: string, metadata: Record<string, unknown>) {
  if (process.env.NODE_ENV === 'test') {
    return;
  }
  console.error(`[backend] ${message}`, metadata);
}

async function requestJson<T>(path: string): Promise<T> {
  const url = `${getBackendBaseUrl()}${path}`;
  let response: Response;

  try {
    response = await fetch(url, {
      method: 'GET',
      headers: {
        Accept: 'application/json'
      },
      cache: 'no-store'
    });
  } catch (error) {
    const message = error instanceof Error ? error.message : 'Unknown error';
    logBackendFailure('Failed to reach backend', { path, url, message });
    throw new Error(`Failed to reach backend for ${path}: ${message}`);
  }

  const rawBody = await response.text();
  if (!response.ok) {
    const preview = previewBody(rawBody);
    logBackendFailure('Backend returned non-2xx response', {
      path,
      url,
      status: response.status,
      body: preview
    });
    throw new Error(`Backend request to ${path} failed (${response.status}): ${preview}`);
  }

  try {
    return JSON.parse(rawBody) as T;
  } catch (error) {
    logBackendFailure('Failed to parse backend JSON', {
      path,
      url,
      status: response.status,
      body: previewBody(rawBody),
      message: error instanceof Error ? error.message : 'Unknown error'
    });
    throw new Error(`Backend response for ${path} was not valid JSON`);
  }
}

export async function listBooks() {
  const payload = await requestJson<ListBooksResponse>('/books');
  return payload.books ?? [];
}

export async function getBook(id: string) {
  return requestJson<Book>(`/books/${encodeURIComponent(id)}`);
}

export async function fetchBooksHealth() {
  return requestJson<BooksHealthResponse>('/health/books');
}

export async function fetchBackendHealthDiagnostics(): Promise<BackendDiagnostics> {
  const url = `${getBackendBaseUrl()}/health/books`;
  try {
    const response = await fetch(url, {
      method: 'GET',
      headers: {
        Accept: 'application/json'
      },
      cache: 'no-store'
    });
    const body = await response.text();
    return {
      url,
      status: response.status,
      ok: response.ok,
      body,
      fetchedAt: new Date().toISOString()
    };
  } catch (error) {
    const message = error instanceof Error ? error.message : 'Unknown error';
    logBackendFailure('Backend diagnostics request failed', { url, message });
    return {
      url,
      status: 0,
      ok: false,
      body: '',
      fetchedAt: new Date().toISOString(),
      error: message
    };
  }
}
