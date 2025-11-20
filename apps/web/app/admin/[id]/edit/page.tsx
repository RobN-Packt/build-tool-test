import { notFound } from 'next/navigation';
import Link from 'next/link';
import { getBook } from '@/lib/api/server';
import { BookForm } from '@/components/BookForm';

interface EditBookPageProps {
  params: {
    id: string;
  };
}

export const metadata = {
  title: 'Edit Book — Packt Library'
};

export default async function EditBookPage({ params }: EditBookPageProps) {
  const book = await getBook(params.id).catch(() => null);
  if (!book) {
    notFound();
  }

  return (
    <div className="page">
      <Link href="/" className="back-link">
        Back to inventory
      </Link>
      <section className="card form-card">
        <div>
          <h1 className="page-title">Edit book</h1>
          <p className="text-muted">
            Update the details for <strong>{book.title}</strong>. Your changes will be saved instantly.
          </p>
        </div>
        <BookForm mode="edit" book={book} />
      </section>
    </div>
  );
}

