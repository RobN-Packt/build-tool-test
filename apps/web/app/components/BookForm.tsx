'use client';

import { FormEvent, useState } from 'react';
import type { components } from '@/lib/api/types';
import type { Book } from './BookTable';
import { apiClient } from '@/lib/api/client';

export type { Book } from './BookTable';
type BookCreate = components['schemas']['BookCreate'];

type BookFormProps = {
  onCreated?: (book: Book) => void;
  submitLabel?: string;
  showSuccessMessage?: boolean;
};

type FormState = {
  title: string;
  author: string;
  price: string;
  currency: string;
  stock: string;
};

const defaultState: FormState = {
  title: '',
  author: '',
  price: '',
  currency: 'USD',
  stock: '0'
};

export function BookForm({ onCreated, submitLabel = 'Save Book', showSuccessMessage = false }: BookFormProps) {
  const [form, setForm] = useState<FormState>(defaultState);
  const [submitting, setSubmitting] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [success, setSuccess] = useState(false);

  const handleChange = (field: keyof FormState) => (event: FormEvent<HTMLInputElement>) => {
    const target = event.target as HTMLInputElement;
    setForm((prev) => ({ ...prev, [field]: target.value }));
  };

  const handleSubmit = async (event: FormEvent<HTMLFormElement>) => {
    event.preventDefault();
    setSubmitting(true);
    setError(null);
    setSuccess(false);

    const parsedPrice = parseFloat(form.price || '0');
    const parsedStock = parseInt(form.stock || '0', 10);

    const payload: BookCreate = {
      title: form.title.trim(),
      author: form.author.trim(),
      price: Number.isNaN(parsedPrice) ? 0 : parsedPrice,
      currency: form.currency.trim().toUpperCase() || 'USD',
      stock: Number.isNaN(parsedStock) ? 0 : parsedStock
    };

    try {
      const { data, error: clientError } = await apiClient.POST('/books', {
        body: payload
      });

      if (clientError || !data) {
        throw clientError ?? new Error('Unable to create book');
      }

      setForm(defaultState);
      setSuccess(true);
      onCreated?.(data);
    } catch (err) {
      if (err instanceof Error) {
        setError(err.message);
      } else {
        setError('Unexpected error occurred');
      }
    } finally {
      setSubmitting(false);
    }
  };

  return (
    <form onSubmit={handleSubmit} aria-label="Create book form">
      <h2>Add a Book</h2>
      {error ? <div role="alert" className="error">{error}</div> : null}
      {showSuccessMessage && success ? (
        <div role="status" className="success">Book saved</div>
      ) : null}
      <label htmlFor="title">
        Title
        <input id="title" name="title" required value={form.title} onInput={handleChange('title')} />
      </label>
      <label htmlFor="author">
        Author
        <input id="author" name="author" required value={form.author} onInput={handleChange('author')} />
      </label>
      <label htmlFor="price">
        Price
        <input
          id="price"
          name="price"
          type="number"
          min="0"
          step="0.01"
          required
          value={form.price}
          onInput={handleChange('price')}
        />
      </label>
      <label htmlFor="currency">
        Currency
        <input
          id="currency"
          name="currency"
          maxLength={3}
          value={form.currency}
          onInput={handleChange('currency')}
        />
      </label>
      <label htmlFor="stock">
        Stock
        <input
          id="stock"
          name="stock"
          type="number"
          min="0"
          step="1"
          value={form.stock}
          onInput={handleChange('stock')}
        />
      </label>
      <div className="actions">
        <button type="submit" disabled={submitting} aria-busy={submitting}>
          {submitting ? 'Saving…' : submitLabel}
        </button>
      </div>
    </form>
  );
}
