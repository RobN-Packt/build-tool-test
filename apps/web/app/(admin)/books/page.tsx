import Link from "next/link";

import { apiClient } from "@/lib/api/client";

async function getBooks() {
  const response = await apiClient.GET("/books");
  return response.data?.data ?? [];
}

export default async function AdminBooksPage() {
  const books = await getBooks();

  return (
    <section className="space-y-6">
      <header className="flex items-center justify-between">
        <div>
          <h2 className="text-2xl font-semibold">Admin: Books</h2>
          <p className="text-sm text-slate-600">Create, update, or delete books.</p>
        </div>
        <Link className="rounded bg-blue-600 px-3 py-2 text-sm font-medium text-white" href="/admin/books/new">
          New book
        </Link>
      </header>
      <ul className="space-y-3">
        {books.map((book) => (
          <li key={book.id} className="flex items-center justify-between rounded border border-slate-200 bg-white px-4 py-3 text-sm">
            <div>
              <p className="font-medium">{book.title}</p>
              <p className="text-slate-500">{book.author}</p>
            </div>
            <div className="flex items-center gap-2">
              <Link className="rounded border border-slate-300 px-2 py-1" href={`/admin/books/${book.id}`}>
                Edit
              </Link>
            </div>
          </li>
        ))}
        {books.length === 0 && <li className="text-sm text-slate-500">No books yet.</li>}
      </ul>
    </section>
  );
}
