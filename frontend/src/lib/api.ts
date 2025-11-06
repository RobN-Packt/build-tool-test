import { Book, BookFormValues } from "@/types/book";

const API_BASE_URL = process.env.NEXT_PUBLIC_API_URL ?? "http://localhost:8080";

async function handleResponse<T>(response: Response): Promise<T> {
  if (!response.ok) {
    let message = response.statusText;

    try {
      const errorBody = await response.json();
      message = errorBody?.error ?? JSON.stringify(errorBody);
    } catch {
      // ignore parsing errors
    }

    throw new Error(message || "Request failed");
  }

  return response.json() as Promise<T>;
}

export async function fetchBooks(): Promise<Book[]> {
  const response = await fetch(`${API_BASE_URL}/books`, { cache: "no-store" });
  const payload = await handleResponse<{ data: Book[] }>(response);
  return payload.data;
}

export async function createBook(values: BookFormValues): Promise<Book> {
  const response = await fetch(`${API_BASE_URL}/books`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(toRequestPayload(values)),
  });

  return handleResponse<Book>(response);
}

export async function updateBook(
  id: number,
  values: BookFormValues,
): Promise<Book> {
  const response = await fetch(`${API_BASE_URL}/books/${id}`, {
    method: "PUT",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(toRequestPayload(values)),
  });

  return handleResponse<Book>(response);
}

export async function deleteBook(id: number): Promise<void> {
  const response = await fetch(`${API_BASE_URL}/books/${id}`, {
    method: "DELETE",
  });

  await handleResponse(response);
}

function toRequestPayload(values: BookFormValues) {
  return {
    title: values.title,
    author: values.author,
    isbn: values.isbn,
    price: Number(values.price),
    stock: Number(values.stock),
    description: values.description,
    publishedDate: values.publishedDate,
  };
}

