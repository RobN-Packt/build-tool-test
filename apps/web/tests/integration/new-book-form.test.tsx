import { act, render, screen, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { NewBookForm } from '@/components/NewBookForm';

const pushMock = vi.fn();

vi.mock('next/navigation', () => ({
  useRouter: () => ({
    push: pushMock
  })
}));

describe('NewBookForm', () => {
  beforeEach(() => {
    pushMock.mockReset();
    vi.restoreAllMocks();
  });

  it('submits form data via fetch and redirects', async () => {
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

    render(<NewBookForm />);

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
});
