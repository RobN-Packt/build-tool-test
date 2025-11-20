import Link from 'next/link';
import { listBooks } from '@/lib/api/server';
import { BookTable } from '@/components/BookTable';

export const revalidate = 0;

export default async function Page() {
  const books = await listBooks().catch(() => []);

  return (
    <div className="page">
      <section className="hero">
        <span className="hero-eyebrow">Packt Publishing</span>
        <h1 className="hero-title">Curate your modern technical library.</h1>
        <p className="hero-subtitle">
          Manage a collection of Packt titles, keep inventory healthy, and empower your teams with the
          latest insights from experts around the world.
        </p>
        <div className="hero-actions">
          <Link href="/admin/new" className="button">
            Add a book
          </Link>
          <a
            href="https://www.packtpub.com"
            target="_blank"
            rel="noreferrer"
            className="button secondary"
          >
            Explore Packt.com
          </a>
        </div>
      </section>

      <section className="card">
        <header className="section-header">
          <div>
            <h2 className="page-title">Inventory</h2>
            <p className="text-muted">
              {books.length === 0
                ? 'No books yet — add your first Packt title to get started.'
                : `You have ${books.length} curated ${books.length === 1 ? 'title' : 'titles'} ready for your readers.`}
            </p>
          </div>
          <Link href="/admin/new" className="button secondary">
            New entry
          </Link>
        </header>
        <BookTable books={books} />
      </section>
    </div>
  );
}
