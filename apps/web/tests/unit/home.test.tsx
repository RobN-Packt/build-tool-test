import { render, screen } from '@testing-library/react';
import { BookTable } from '@/components/BookTable';
import type { Book } from '@/lib/api';

vi.mock('next/navigation', () => ({
  useRouter: () => ({
    push: vi.fn(),
    refresh: vi.fn()
  })
}));

describe('Home page table', () => {
  it('renders book rows', () => {
    const books: Book[] = [
      {
        id: '1',
        title: 'Go in Action',
        author: 'William Kennedy',
        price: 39.99,
        currency: 'USD',
        stock: 3,
        createdAt: new Date().toISOString(),
        updatedAt: new Date().toISOString()
      },
      {
        id: '2',
        title: 'Concurrency in Go',
        author: 'Katherine Cox-Buday',
        price: 29.99,
        currency: 'USD',
        stock: 2,
        createdAt: new Date().toISOString(),
        updatedAt: new Date().toISOString()
      }
    ];

    render(<BookTable books={books} />);

    expect(screen.getByText('Go in Action')).toBeInTheDocument();
    expect(screen.getByText('Concurrency in Go')).toBeInTheDocument();
    expect(screen.getAllByRole('row')).toHaveLength(books.length + 1); // includes header row
  });
});
