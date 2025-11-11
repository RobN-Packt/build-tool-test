'use client';

import { FormEvent, useState } from 'react';
import { useRouter } from 'next/navigation';
import { createBook, type BookCreate } from '@/lib/api';

const defaultForm = {
  title: '',
  author: '',
  price: '',
  currency: 'USD',
  stock: ''
};

export function NewBookForm() {
  const router = useRouter();
  const [form, setForm] = useState(defaultForm);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const handleChange = (event: React.ChangeEvent<HTMLInputElement>) => {
    const { name, value } = event.target;
    setForm((prev) => ({ ...prev, [name]: value }));
  };

  const handleSubmit = async (event: FormEvent<HTMLFormElement>) => {
    event.preventDefault();
    setError(null);
    setLoading(true);
    try {
      const payload: BookCreate = {
        title: form.title.trim(),
        author: form.author.trim(),
        price: Number(form.price),
        currency: form.currency.trim().toUpperCase(),
        stock: Number(form.stock)
      };
      await createBook(payload);
      setForm(defaultForm);
      router.push('/');
    } catch (err) {
      const message = err instanceof Error ? err.message : 'Failed to create book.';
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
          {loading ? 'Saving…' : 'Create book'}
        </button>
      </div>
    </form>
  );
}
