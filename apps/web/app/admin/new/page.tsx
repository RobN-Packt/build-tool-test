import Link from 'next/link';
import { NewBookForm } from '@/components/NewBookForm';

export const metadata = {
  title: 'Add Book — Packt Library'
};

export default function NewBookPage() {
  return (
    <div className="page">
      <Link href="/" className="back-link">
        Back to inventory
      </Link>
      <section className="card form-card">
        <div>
          <h1 className="page-title">Create a new Packt book</h1>
          <p className="text-muted">
            Enter the details below to add a new title to your Packt library. All fields are required.
          </p>
        </div>
        <NewBookForm />
      </section>
    </div>
  );
}
