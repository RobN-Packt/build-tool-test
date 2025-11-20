'use client';

import Link from 'next/link';
import { useRouter } from 'next/navigation';
import { useState } from 'react';
import { deleteBook } from '@/lib/api/client';

interface BookActionsProps {
  id: string;
  title: string;
}

export function BookActions({ id, title }: BookActionsProps) {
  const router = useRouter();
  const [isDeleting, setIsDeleting] = useState(false);

  const handleDelete = async () => {
    if (isDeleting) return;

    const confirmed = window.confirm(
      `Are you sure you want to remove “${title}” from your Packt library?`
    );
    if (!confirmed) return;

    setIsDeleting(true);
    try {
      await deleteBook(id);
      router.refresh();
    } catch (error) {
      const message = error instanceof Error ? error.message : 'Failed to delete book.';
      alert(message);
    } finally {
      setIsDeleting(false);
    }
  };

  return (
    <div className="table-actions">
      <Link href={`/admin/${id}/edit`} className="action-link">
        Edit
      </Link>
      <button type="button" className="action-link danger" onClick={handleDelete} disabled={isDeleting}>
        {isDeleting ? 'Deleting…' : 'Delete'}
      </button>
    </div>
  );
}

