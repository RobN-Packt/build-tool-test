"use client";

import { useRouter } from "next/navigation";
import { useEffect, useState } from "react";

import { apiClient } from "@/lib/api/client";

type FormState = {
  title: string;
  author: string;
  price: string;
  currency: string;
  stock: string;
};

interface Props {
  params: { id: string };
}

export default function EditBookPage({ params }: Props) {
  const router = useRouter();
  const [form, setForm] = useState<FormState | null>(null);
  const [error, setError] = useState<string | null>(null);
  const [submitting, setSubmitting] = useState(false);
  const [loading, setLoading] = useState(true);
  const id = params.id;

  useEffect(() => {
    async function load() {
      const response = await apiClient.GET("/books/{id}", { params: { path: { id } } });
      if (response.error || !response.data) {
        setError(response.error?.message ?? "Book not found");
        setLoading(false);
        return;
      }
      const { data } = response.data;
      setForm({
        title: data.title,
        author: data.author,
        price: String(data.price),
        currency: data.currency,
        stock: String(data.stock),
      });
      setLoading(false);
    }
    load();
  }, [id]);

  async function handleSubmit(event: React.FormEvent<HTMLFormElement>) {
    event.preventDefault();
    if (!form) return;
    setSubmitting(true);
    setError(null);
    const result = await apiClient.PUT("/books/{id}", {
      params: { path: { id } },
      body: {
        title: form.title,
        author: form.author,
        price: Number(form.price),
        currency: form.currency,
        stock: Number(form.stock),
      },
    });

    if (result.error) {
      setError(result.error.message ?? "Failed to update book");
      setSubmitting(false);
      return;
    }

    router.push("/admin/books");
  }

  async function handleDelete() {
    setSubmitting(true);
    const result = await apiClient.DELETE("/books/{id}", { params: { path: { id } } });
    if (result.error) {
      setError(result.error.message ?? "Failed to delete book");
      setSubmitting(false);
      return;
    }
    router.push("/admin/books");
  }

  if (loading || !form) {
    return <p className="text-sm text-slate-500">Loading...</p>;
  }

  return (
    <section className="mx-auto max-w-lg space-y-6">
      <header className="space-y-2">
        <h2 className="text-2xl font-semibold">Edit book</h2>
        <p className="text-sm text-slate-600">Update details or remove the book.</p>
      </header>
      <form className="space-y-4" onSubmit={handleSubmit}>
        <div className="space-y-1">
          <label className="text-sm font-medium" htmlFor="title">
            Title
          </label>
          <input
            className="w-full rounded border border-slate-300 px-3 py-2"
            id="title"
            required
            value={form.title}
            onChange={(e) => setForm((prev) => prev && { ...prev, title: e.target.value })}
          />
        </div>
        <div className="space-y-1">
          <label className="text-sm font-medium" htmlFor="author">
            Author
          </label>
          <input
            className="w-full rounded border border-slate-300 px-3 py-2"
            id="author"
            required
            value={form.author}
            onChange={(e) => setForm((prev) => prev && { ...prev, author: e.target.value })}
          />
        </div>
        <div className="grid grid-cols-2 gap-4">
          <div className="space-y-1">
            <label className="text-sm font-medium" htmlFor="price">
              Price
            </label>
            <input
              className="w-full rounded border border-slate-300 px-3 py-2"
              id="price"
              required
              type="number"
              step="0.01"
              value={form.price}
              onChange={(e) => setForm((prev) => prev && { ...prev, price: e.target.value })}
            />
          </div>
          <div className="space-y-1">
            <label className="text-sm font-medium" htmlFor="stock">
              Stock
            </label>
            <input
              className="w-full rounded border border-slate-300 px-3 py-2"
              id="stock"
              required
              type="number"
              min="0"
              value={form.stock}
              onChange={(e) => setForm((prev) => prev && { ...prev, stock: e.target.value })}
            />
          </div>
        </div>
        <div className="space-y-1">
          <label className="text-sm font-medium" htmlFor="currency">
            Currency
          </label>
          <input
            className="w-full rounded border border-slate-300 px-3 py-2"
            id="currency"
            required
            maxLength={3}
            value={form.currency}
            onChange={(e) => setForm((prev) => prev && { ...prev, currency: e.target.value.toUpperCase() })}
          />
        </div>
        {error && <p className="text-sm text-red-600">{error}</p>}
        <div className="flex items-center gap-3">
          <button
            className="flex-1 rounded bg-blue-600 px-3 py-2 text-sm font-medium text-white disabled:bg-blue-300"
            disabled={submitting}
            type="submit"
          >
            {submitting ? "Saving..." : "Save"}
          </button>
          <button
            className="rounded border border-red-200 px-3 py-2 text-sm text-red-600 disabled:text-red-300"
            disabled={submitting}
            onClick={(event) => {
              event.preventDefault();
              void handleDelete();
            }}
          >
            Delete
          </button>
        </div>
      </form>
    </section>
  );
}
