import { Book } from '@/lib/api';

interface BookTableProps {
  books: Book[];
}

export function BookTable({ books }: BookTableProps) {
  if (!books.length) {
    return <p>No books available yet.</p>;
  }

  return (
    <table style={{ borderCollapse: 'collapse', width: '100%' }}>
      <thead>
        <tr>
          <th style={cellHeadStyle}>Title</th>
          <th style={cellHeadStyle}>Author</th>
          <th style={cellHeadStyle}>Price</th>
          <th style={cellHeadStyle}>Stock</th>
        </tr>
      </thead>
      <tbody>
        {books.map((book) => (
          <tr key={book.id}>
            <td style={cellBodyStyle}>{book.title}</td>
            <td style={cellBodyStyle}>{book.author}</td>
            <td style={cellBodyStyle}>
              {book.currency} {Number(book.price).toFixed(2)}
            </td>
            <td style={cellBodyStyle}>{book.stock}</td>
          </tr>
        ))}
      </tbody>
    </table>
  );
}

const cellHeadStyle: React.CSSProperties = {
  borderBottom: '1px solid #ddd',
  padding: '0.5rem',
  textAlign: 'left',
  fontWeight: 600
};

const cellBodyStyle: React.CSSProperties = {
  borderBottom: '1px solid #f0f0f0',
  padding: '0.5rem'
};
