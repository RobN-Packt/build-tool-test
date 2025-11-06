import Link from "next/link";

import { apiClient } from "@/lib/api/client";

interface Props {
  params: { id: string };
}

export default async function BookDetailPage({ params }: Props) {
  const { id } = params;
  const response = await apiClient.GET("/books/{id}", { params: { path: { id } } });

  if (response.error || !response.data) {
    return (
      <section className="space-y-4">
        <h2 className="text-2xl font-semibold">Book not found</h2>
        <Link className="text-blue-600" href="/">
          Go back
        </Link>
      </section>
    );
  }

  const { data: book } = response.data;

  return (
    <section className="space-y-6">
      <header className="space-y-1">
        <h2 className="text-3xl font-semibold">{book.title}</h2>
        <p className="text-slate-600">by {book.author}</p>
      </header>
      <div className="rounded border border-slate-200 bg-white p-6 text-sm">
        <p>
          <span className="font-medium">Price:</span> {new Intl.NumberFormat("en-US", { style: "currency", currency: book.currency }).format(book.price)}
        </p>
        <p>
          <span className="font-medium">Stock:</span> {book.stock}
        </p>
        <p>
          <span className="font-medium">Created:</span> {new Date(book.created_at).toLocaleString()}
        </p>
        <p>
          <span className="font-medium">Updated:</span> {new Date(book.updated_at).toLocaleString()}
        </p>
      </div>
      <div className="flex items-center gap-4">
        <Link className="rounded bg-slate-200 px-3 py-2 text-sm" href="/">
          Back to list
        </Link>
        <Link className="rounded bg-blue-600 px-3 py-2 text-sm text-white" href={`/admin/books/${book.id}`}>
          Edit
        </Link>
      </div>
    </section>
  );
}
