"use client";

import { useRouter } from "next/navigation";
import { useState } from "react";

import { apiClient } from "@/lib/api/client";

export default function NewBookPage() {
  const router = useRouter();
  const [form, setForm] = useState({ title: "", author: "", price: "", currency: "USD", stock: "0" });
  const [error, setError] = useState<string | null>(null);
  const [submitting, setSubmitting] = useState(false);

  async function handleSubmit(event: React.FormEvent<HTMLFormElement>) {
    event.preventDefault();
    setSubmitting(true);
    setError(null);
    const result = await apiClient.POST("/books", {
      body: {
        title: form.title,
        author: form.author,
        price: Number(form.price),
        currency: form.currency,
        stock: Number(form.stock),
      },
    });

    if (result.error) {
      setError(result.error.message ?? "Failed to create book");
      setSubmitting(false);
      return;
    }

    router.push("/admin/books");
  }

  return (
    <section className="mx-auto max-w-lg space-y-6">
      <header className="space-y-2">
        <h2 className="text-2xl font-semibold">Create new book</h2>
        <p className="text-sm text-slate-600">Fill in the details below.</p>
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
            onChange={(e) => setForm((prev) => ({ ...prev, title: e.target.value }))}
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
            onChange={(e) => setForm((prev) => ({ ...prev, author: e.target.value }))}
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
              onChange={(e) => setForm((prev) => ({ ...prev, price: e.target.value }))}
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
              onChange={(e) => setForm((prev) => ({ ...prev, stock: e.target.value }))}
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
            onChange={(e) => setForm((prev) => ({ ...prev, currency: e.target.value.toUpperCase() }))}
          />
        </div>
        {error && <p className="text-sm text-red-600">{error}</p>}
        <button
          className="w-full rounded bg-blue-600 px-3 py-2 text-sm font-medium text-white disabled:bg-blue-300"
          disabled={submitting}
          type="submit"
        >
          {submitting ? "Saving..." : "Save"}
        </button>
      </form>
    </section>
  );
}
