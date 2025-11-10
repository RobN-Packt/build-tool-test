import Link from 'next/link';
import { BooksView } from './components/BooksView';
import { apiClient } from '@/lib/api/client';
import type { Book } from './components/BookTable';

export default async function HomePage() {
  const { data, error } = await apiClient.GET('/books');
  const books = (data ?? []) as Book[];

  return (
    <main>
      <header>
        <h1>Book Shop Inventory</h1>
        <p>Track catalog titles, prices, and stock for the shop.</p>
        <Link href="/admin/new">Go to admin create screen</Link>
      </header>
      {error ? (
        <p className="error">Failed to load books. Try again soon.</p>
      ) : (
        <BooksView initialBooks={books} />
      )}
    </main>
  );
}
