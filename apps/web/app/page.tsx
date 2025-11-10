import { listBooks } from '@/lib/api';
import { BookTable } from '@/components/BookTable';

export const revalidate = 0;

export default async function Page() {
  const books = await listBooks().catch(() => []);

  return (
    <section style={{ display: 'grid', gap: '1.5rem' }}>
      <header>
        <h1 style={{ marginBottom: '0.25rem' }}>Books</h1>
        <p style={{ color: '#4b5563' }}>
          Browse the catalog. Use the New Book link to add items to the store.
        </p>
      </header>
      <BookTable books={books} />
    </section>
  );
}
