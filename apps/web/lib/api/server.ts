import type { components } from './types';

type Book = components['schemas']['Book'];

interface ListBooksResponse {
  books: Book[];
}

function getBackendBaseUrl() {
  const base = process.env.BACKEND_INTERNAL_URL;
  if (!base) {
    throw new Error('BACKEND_INTERNAL_URL is not configured');
  }
  return base.endsWith('/') ? base.slice(0, -1) : base;
}

async function requestJson<T>(path: string): Promise<T> {
  const url = `${getBackendBaseUrl()}${path}`;
  const response = await fetch(url, {
    method: 'GET',
    headers: {
      Accept: 'application/json'
    },
    cache: 'no-store'
  });

  if (!response.ok) {
    const detail = await response.text();
    throw new Error(detail || `Backend request to ${path} failed with status ${response.status}`);
  }

  return (await response.json()) as T;
}

export async function listBooks() {
  const payload = await requestJson<ListBooksResponse>('/books');
  return payload.books ?? [];
}

export async function getBook(id: string) {
  return requestJson<Book>(`/books/${encodeURIComponent(id)}`);
}
