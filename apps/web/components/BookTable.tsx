import { Book } from '@/lib/api';

interface BookTableProps {
  books: Book[];
}

export function BookTable({ books }: BookTableProps) {
  if (!books.length) {
    return <div className="empty-state">No Packt titles in your catalog yet.</div>;
  }

  return (
    <div className="table-wrapper">
      <table className="table">
        <thead>
          <tr>
            <th>Title</th>
            <th>Author</th>
            <th>Price</th>
            <th>Stock</th>
          </tr>
        </thead>
        <tbody>
          {books.map((book) => (
            <tr key={book.id}>
              <td>{book.title}</td>
              <td>{book.author}</td>
              <td>
                <strong>{book.currency}</strong> {Number(book.price).toFixed(2)}
              </td>
              <td>
                <span className="status-pill">{book.stock} in stock</span>
              </td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
}
