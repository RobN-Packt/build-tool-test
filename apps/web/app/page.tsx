import Link from "next/link";

import { BookTable } from "./components/book-table";
import { apiClient } from "../lib/api/client";

async function getBooks() {
  const response = await apiClient.GET("/books");
  if (response.error) {
    console.error("Failed to fetch books", response.error);
    return [];
  }
  return response.data?.data ?? [];
}

export default async function BooksPage() {
  const books = await getBooks();

  return (
    <section className="space-y-4">
      <header className="flex items-center justify-between">
        <div>
          <h2 className="text-2xl font-semibold">Books</h2>
          <p className="text-sm text-slate-600">Browse available titles and check inventory.</p>
        </div>
        <Link className="rounded bg-blue-600 px-3 py-2 text-sm font-medium text-white" href="/admin/books/new">
          Add book
        </Link>
      </header>
      <BookTable books={books} />
    </section>
  );
}
