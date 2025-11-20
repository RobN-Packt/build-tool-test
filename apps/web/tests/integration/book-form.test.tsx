import { act, render, screen, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { BookForm } from '@/components/BookForm';
import type { Book } from '@/lib/api/types';

const pushMock = vi.fn();

vi.mock('next/navigation', () => ({
  useRouter: () => ({
    push: pushMock,
    refresh: vi.fn()
  })
}));

describe('BookForm', () => {
  beforeEach(() => {
    pushMock.mockReset();
    vi.restoreAllMocks();
  });

  it('submits create form and redirects', async () => {
    const user = userEvent.setup();

    const mockBody = {
      id: '123',
      title: 'New Book',
      author: 'Jane Doe',
      price: 10,
      currency: 'USD',
      stock: 5,
      createdAt: new Date().toISOString(),
      updatedAt: new Date().toISOString()
    };

    const mockResponse = new Response(JSON.stringify(mockBody), {
      status: 201,
      headers: { 'content-type': 'application/json' }
    });

    const fetchSpy = vi.spyOn(global, 'fetch').mockResolvedValue(mockResponse);

    render(<BookForm mode="create" />);

    await act(async () => {
      await user.type(screen.getByLabelText(/title/i), 'New Book');
      await user.type(screen.getByLabelText(/author/i), 'Jane Doe');
      await user.type(screen.getByLabelText(/price/i), '10');
      await user.clear(screen.getByLabelText(/currency/i));
      await user.type(screen.getByLabelText(/currency/i), 'usd');
      await user.type(screen.getByLabelText(/stock/i), '5');
      await user.click(screen.getByRole('button', { name: /create book/i }));
    });

    await waitFor(() => {
      expect(fetchSpy).toHaveBeenCalledTimes(1);
      const [request] = fetchSpy.mock.calls[0] as [Request];
      expect(request.url).toMatch(/\/books$/);
      expect(request.method).toBe('POST');
      expect(pushMock).toHaveBeenCalledWith('/');
    });
  });

  it('submits edit form', async () => {
    const user = userEvent.setup();

    const existing: Book = {
      id: 'abc',
      title: 'Original',
      author: 'Author',
      price: 20,
      currency: 'USD',
      stock: 3,
      createdAt: new Date().toISOString(),
      updatedAt: new Date().toISOString()
    };

    const mockBody = { ...existing, title: 'Updated Title' };

    const mockResponse = new Response(JSON.stringify(mockBody), {
      status: 200,
      headers: { 'content-type': 'application/json' }
    });

    const fetchSpy = vi.spyOn(global, 'fetch').mockResolvedValue(mockResponse);

    render(<BookForm mode="edit" book={existing} />);

    await act(async () => {
      await user.clear(screen.getByLabelText(/title/i));
      await user.type(screen.getByLabelText(/title/i), 'Updated Title');
      await user.click(screen.getByRole('button', { name: /update book/i }));
    });

    await waitFor(() => {
      expect(fetchSpy).toHaveBeenCalledTimes(1);
      const [request] = fetchSpy.mock.calls[0] as [Request];
      expect(request.url).toMatch(/\/books\/abc$/);
      expect(request.method).toBe('PUT');
    });
  });
});

