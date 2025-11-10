import { render, screen, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest';
import { BooksView } from '@/app/components/BooksView';
import type { Book } from '@/app/components/BookTable';

const mockPost = vi.fn();
const mockGet = vi.fn();

vi.mock('@/lib/api/client', () => ({
  apiClient: {
    POST: mockPost,
    GET: mockGet
  }
}));

describe('Book creation flow', () => {
  beforeEach(() => {
    mockPost.mockReset();
    mockGet.mockReset();
  });

  afterEach(() => {
    vi.clearAllMocks();
  });

  it('submits the create form and refreshes the list', async () => {
    const created: Book = {
      id: 'generated-id',
      title: 'Domain-Driven Design',
      author: 'Eric Evans',
      price: 59.99,
      currency: 'USD',
      stock: 7,
      createdAt: new Date('2024-05-01T12:00:00Z').toISOString(),
      updatedAt: new Date('2024-05-01T12:00:00Z').toISOString()
    };

    mockPost.mockResolvedValue({ data: created });
    mockGet.mockResolvedValue({ data: [created] });

    render(<BooksView initialBooks={[]} />);

    const user = userEvent.setup();

    await user.type(screen.getByLabelText('Title'), 'Domain-Driven Design');
    await user.type(screen.getByLabelText('Author'), 'Eric Evans');
    await user.clear(screen.getByLabelText('Price'));
    await user.type(screen.getByLabelText('Price'), '59.99');

    await user.click(screen.getByRole('button', { name: /create/i }));

    await waitFor(() => {
      expect(mockPost).toHaveBeenCalledWith('/books', {
        body: {
          title: 'Domain-Driven Design',
          author: 'Eric Evans',
          price: 59.99,
          currency: 'USD',
          stock: 0
        }
      });
    });

    await waitFor(() => {
      expect(mockGet).toHaveBeenCalled();
    });

    expect(await screen.findByText('Domain-Driven Design')).toBeVisible();
    expect(screen.getByText('Eric Evans')).toBeVisible();
  });
});
