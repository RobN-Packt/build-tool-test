'use client';

import { useState, useTransition } from 'react';
import { BookForm, type Book } from './BookForm';
import { BookTable } from './BookTable';
import { apiClient } from '@/lib/api/client';

interface BooksViewProps {
  initialBooks: Book[];
}

export function BooksView({ initialBooks }: BooksViewProps) {
  const [books, setBooks] = useState<Book[]>(initialBooks);
  const [feedback, setFeedback] = useState<string | null>(null);
  const [refreshing, startRefresh] = useTransition();

  const reloadBooks = async () => {
    const { data, error } = await apiClient.GET('/books');
    if (error) {
      throw error;
    }
    setBooks(data ?? []);
  };

  const handleCreated = (book: Book) => {
    setFeedback(`Added “${book.title}” to inventory.`);
    startRefresh(() => {
      reloadBooks().catch((err) => {
        console.error('Failed to refresh inventory', err);
      });
    });
  };

  return (
    <div>
      {feedback ? <p className="success" role="status">{feedback}</p> : null}
      {refreshing ? <p>Refreshing inventory…</p> : null}
      <BookForm onCreated={handleCreated} submitLabel="Create" showSuccessMessage={false} />
      <section>
        <h2>Inventory</h2>
        <BookTable books={books} />
      </section>
    </div>
  );
}
