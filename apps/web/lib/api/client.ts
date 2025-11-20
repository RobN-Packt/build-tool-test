import createClient from 'openapi-fetch';
import type { paths, components } from './types';

function resolveOrigin() {
  if (typeof window !== 'undefined' && window.location?.origin) {
    return window.location.origin;
  }

  const explicit =
    process.env.NEXT_PUBLIC_SITE_URL ?? process.env.NEXT_PUBLIC_VERCEL_URL;

  if (explicit) {
    const normalized = explicit.startsWith('http') ? explicit : `https://${explicit}`;
    return normalized.replace(/\/$/, '');
  }

  return 'http://localhost:3000';
}

function resolveBaseUrl() {
  return `${resolveOrigin()}/api`;
}

export const client = createClient<paths>({
  baseUrl: resolveBaseUrl(),
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
