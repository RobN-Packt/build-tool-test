'use client';

import Link from 'next/link';
import { useRouter } from 'next/navigation';
import { BookForm, type Book } from '../../components/BookForm';

export default function NewBookPage() {
  const router = useRouter();

  const handleCreated = (_book: Book) => {
    router.push('/');
  };

  return (
    <main>
      <header>
        <h1>Add a New Book</h1>
        <p>Provide complete details for the book. Submit to add to inventory.</p>
        <Link href="/">Back to inventory</Link>
      </header>
      <BookForm onCreated={handleCreated} submitLabel="Create Book" showSuccessMessage />
    </main>
  );
}
