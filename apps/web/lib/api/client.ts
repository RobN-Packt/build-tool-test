import createClient from 'openapi-fetch';
import type { paths, components } from './types';

function resolveBrowserBaseUrl(envBaseUrl?: string) {
  if (envBaseUrl?.startsWith('/')) {
    return envBaseUrl;
  }
  if (envBaseUrl) {
    return envBaseUrl;
  }
  return '/api';
}

function resolveServerBaseUrl(envBaseUrl?: string) {
  const vercelUrl = process.env.VERCEL_URL?.replace(/\/$/, '');
  if (vercelUrl) {
    return `https://${vercelUrl}/api`;
  }
  return envBaseUrl ?? 'http://localhost:8080';
}

const normalizedEnvBaseUrl = process.env.NEXT_PUBLIC_API_BASE_URL?.replace(/\/$/, '');
const baseUrl =
  typeof window === 'undefined'
    ? resolveServerBaseUrl(normalizedEnvBaseUrl)
    : resolveBrowserBaseUrl(normalizedEnvBaseUrl);

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
