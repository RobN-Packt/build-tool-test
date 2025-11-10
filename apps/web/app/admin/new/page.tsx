import { NewBookForm } from '@/components/NewBookForm';

export const metadata = {
  title: 'Create Book — Book Store'
};

export default function NewBookPage() {
  return (
    <section style={{ display: 'grid', gap: '1.5rem' }}>
      <header>
        <h1>Create a new book</h1>
        <p style={{ color: '#4b5563' }}>
          Provide the book details below. All fields are required.
        </p>
      </header>
      <NewBookForm />
    </section>
  );
}
