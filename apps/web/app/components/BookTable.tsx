import type { components } from '@/lib/api/types';
import { formatCurrency } from '@/lib/utils/format';

export type Book = components['schemas']['Book'];

interface BookTableProps {
  books: Book[];
}

export function BookTable({ books }: BookTableProps) {
  if (books.length === 0) {
    return <p>No books available yet. Add your first title to get started.</p>;
  }

  return (
    <table aria-label="Book inventory">
      <thead>
        <tr>
          <th scope="col">Title</th>
          <th scope="col">Author</th>
          <th scope="col">Price</th>
          <th scope="col">Stock</th>
          <th scope="col">Added</th>
        </tr>
      </thead>
      <tbody>
        {books.map((book) => (
          <tr key={book.id}>
            <td>{book.title}</td>
            <td>{book.author}</td>
            <td>{formatCurrency(book.price, book.currency)}</td>
            <td>{book.stock}</td>
            <td>{new Date(book.createdAt).toLocaleString()}</td>
          </tr>
        ))}
      </tbody>
    </table>
  );
}
