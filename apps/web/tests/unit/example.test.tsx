import { render, screen } from '@testing-library/react';
import { BookTable, type Book } from '@/app/components/BookTable';

describe('BookTable', () => {
  it('renders book rows with currency formatting', () => {
    const books: Book[] = [
      {
        id: '1',
        title: 'Clean Code',
        author: 'Robert C. Martin',
        price: 49.99,
        currency: 'USD',
        stock: 3,
        createdAt: new Date('2024-01-01T12:00:00Z').toISOString(),
        updatedAt: new Date('2024-01-01T12:00:00Z').toISOString()
      }
    ];

    render(<BookTable books={books} />);

    expect(screen.getByRole('table', { name: 'Book inventory' })).toBeInTheDocument();
    expect(screen.getByText('Clean Code')).toBeVisible();
    expect(screen.getByText('Robert C. Martin')).toBeVisible();
    expect(screen.getByText('$49.99')).toBeVisible();
  });
});
