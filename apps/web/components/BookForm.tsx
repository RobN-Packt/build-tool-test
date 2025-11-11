'use client';

import { FormEvent, useMemo, useState } from 'react';
import { useRouter } from 'next/navigation';
import { createBook, updateBook, type Book, type BookCreate, type BookUpdate } from '@/lib/api';

type Mode = 'create' | 'edit';

interface BookFormProps {
  mode: Mode;
  book?: Book;
}

interface FormState {
  title: string;
  author: string;
  price: string;
  currency: string;
  stock: string;
}

function initFormState(book?: Book): FormState {
  if (!book) {
    return {
      title: '',
      author: '',
      price: '',
      currency: 'USD',
      stock: ''
    };
  }
  return {
    title: book.title,
    author: book.author,
    price: String(book.price ?? ''),
    currency: book.currency,
    stock: String(book.stock ?? '')
  };
}

export function BookForm({ mode, book }: BookFormProps) {
  const router = useRouter();
  const [form, setForm] = useState<FormState>(() => initFormState(book));
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const title = useMemo(() => (mode === 'edit' ? 'Update book' : 'Create book'), [mode]);

  const handleChange = (event: React.ChangeEvent<HTMLInputElement>) => {
    const { name, value } = event.target;
    setForm((prev) => ({ ...prev, [name]: value }));
  };

  const handleSubmit = async (event: FormEvent<HTMLFormElement>) => {
    event.preventDefault();
    if (loading) return;

    setError(null);
    setLoading(true);

    try {
      const payloadBase: BookCreate = {
        title: form.title.trim(),
        author: form.author.trim(),
        price: Number(form.price),
        currency: form.currency.trim().toUpperCase(),
        stock: Number(form.stock)
      };

      if (mode === 'create') {
        await createBook(payloadBase);
      } else if (book) {
        const updatePayload: BookUpdate = {
          title: payloadBase.title,
          author: payloadBase.author,
          price: payloadBase.price,
          currency: payloadBase.currency,
          stock: payloadBase.stock
        };
        await updateBook(book.id, updatePayload);
      }

      setForm(initFormState(mode === 'edit' ? book : undefined));
      router.push('/');
      router.refresh();
    } catch (err) {
      const message = err instanceof Error ? err.message : 'Failed to submit the form.';
      setError(message);
    } finally {
      setLoading(false);
    }
  };

  return (
    <form className="form-card" onSubmit={handleSubmit}>
      <div className="form-grid">
        <div className="form-field">
          <label htmlFor="title">Title</label>
          <input
            id="title"
            name="title"
            type="text"
            required
            value={form.title}
            onChange={handleChange}
            placeholder="Hands-On Machine Learning with Go"
            className="input-field"
          />
        </div>

        <div className="form-field">
          <label htmlFor="author">Author</label>
          <input
            id="author"
            name="author"
            type="text"
            required
            value={form.author}
            onChange={handleChange}
            placeholder="Jane Doe"
            className="input-field"
          />
        </div>

        <div className="form-field">
          <label htmlFor="price">Price</label>
          <input
            id="price"
            name="price"
            type="number"
            min="0"
            step="0.01"
            required
            value={form.price}
            onChange={handleChange}
            placeholder="34.99"
            className="input-field"
          />
        </div>

        <div className="form-field">
          <label htmlFor="currency">Currency</label>
          <input
            id="currency"
            name="currency"
            type="text"
            required
            maxLength={3}
            value={form.currency}
            onChange={handleChange}
            placeholder="USD"
            className="input-field"
          />
        </div>

        <div className="form-field">
          <label htmlFor="stock">Stock</label>
          <input
            id="stock"
            name="stock"
            type="number"
            min="0"
            required
            value={form.stock}
            onChange={handleChange}
            placeholder="25"
            className="input-field"
          />
        </div>
      </div>

      {error ? (
        <p role="alert" className="form-error">
          {error}
        </p>
      ) : null}

      <div className="form-actions">
        <button type="submit" className="button" disabled={loading}>
          {loading ? 'Saving…' : title}
        </button>
      </div>
    </form>
  );
}

