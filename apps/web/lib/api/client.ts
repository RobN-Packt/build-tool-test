import createClient from 'openapi-fetch';
import type { paths, components } from './types';

const publicBaseUrl = process.env.NEXT_PUBLIC_API_BASE_URL?.replace(/\/$/, '');
const serverBaseUrl =
  process.env.API_SERVER_BASE_URL?.replace(/\/$/, '') ??
  publicBaseUrl ??
  'http://localhost:8080';

function resolveBrowserBaseUrl() {
  if (!publicBaseUrl) {
    return `${window.location.origin}/api`;
  }

  const isHttpsPage = window.location.protocol === 'https:';
  const isPublicBaseHttp = publicBaseUrl.startsWith('http://');

  if (isHttpsPage && isPublicBaseHttp) {
    // Avoid mixed-content rejections by falling back to the same-origin proxy.
    return `${window.location.origin}/api`;
  }

  return publicBaseUrl;
}

const baseUrl = typeof window === 'undefined' ? serverBaseUrl : resolveBrowserBaseUrl();

export const client = createClient<paths>({
  baseUrl,
  fetch: (...args) => fetch(...args)
});

export type Book = components['schemas']['Book'];
export type BookCreate = components['schemas']['BookCreate'];
export type BookUpdate = components['schemas']['BookUpdate'];
export type ErrorResponse = components['schemas']['Error'];

function extractErrorMessage(error: unknown): string {
  if (error && typeof error === 'object' && 'message' in error && typeof error.message === 'string') {
    return error.message;
  }
  if (error && typeof error === 'object' && 'detail' in error && typeof error.detail === 'string') {
    return error.detail;
  }
  return 'Unknown error';
}

export async function listBooks() {
  const { data, error } = await client.GET('/books', {});
  if (error) {
    throw new Error(extractErrorMessage(error));
  }
  return data?.books ?? [];
}

export async function getBook(id: string) {
  const { data, error } = await client.GET('/books/{id}', {
    params: { path: { id } }
  });
  if (error) {
    throw new Error(extractErrorMessage(error));
  }
  return data ?? null;
}

export async function createBook(body: BookCreate) {
  const { data, error } = await client.POST('/books', {
    body
  });
  if (error) {
    throw new Error(extractErrorMessage(error));
  }
  if (!data) {
    throw new Error('Unexpected empty response');
  }
  return data;
}

export async function updateBook(id: string, body: BookUpdate) {
  const { data, error } = await client.PUT('/books/{id}', {
    params: { path: { id } },
    body
  });
  if (error) {
    throw new Error(extractErrorMessage(error));
  }
  if (!data) {
    throw new Error('Unexpected empty response');
  }
  return data;
}

export async function deleteBook(id: string) {
  const { error } = await client.DELETE('/books/{id}', {
    params: { path: { id } }
  });
  if (error) {
    throw new Error(extractErrorMessage(error));
  }
}
