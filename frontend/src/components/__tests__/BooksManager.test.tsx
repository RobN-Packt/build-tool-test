import { render, screen, waitFor, within } from "@testing-library/react";
import userEvent from "@testing-library/user-event";

import { BooksManager } from "@/components/BooksManager";
import { Book } from "@/types/book";

jest.mock("@/lib/api", () => ({
  fetchBooks: jest.fn(),
  createBook: jest.fn(),
  updateBook: jest.fn(),
  deleteBook: jest.fn(),
}));

const api = jest.requireMock("@/lib/api");

const sampleBooks: Book[] = [
  {
    id: 1,
    title: "Clean Code",
    author: "Robert C. Martin",
    isbn: "9780132350884",
    price: 39.5,
    stock: 10,
    description: "Agile craftsmanship",
    publishedDate: "2008-08-01T00:00:00Z",
    createdAt: "2024-01-01T00:00:00Z",
    updatedAt: "2024-01-01T00:00:00Z",
  },
];

describe("BooksManager", () => {
  beforeEach(() => {
    jest.resetAllMocks();
  });

  it("renders books after initial load", async () => {
    api.fetchBooks.mockResolvedValue(sampleBooks);

    render(<BooksManager />);

    expect(screen.getByText(/loading books/i)).toBeInTheDocument();

    await waitFor(() => expect(screen.getByText("Clean Code")).toBeInTheDocument());

    const summary = screen.getByTestId("inventory-summary");
    expect(within(summary).getByText("Books").nextSibling?.textContent).toBe("1");
  });

  it("submits create book form", async () => {
    api.fetchBooks.mockResolvedValue([]);
    api.createBook.mockImplementation(async () => ({
      id: 2,
      title: "New Book",
      author: "Jane Doe",
      isbn: "9991112223334",
      price: 25,
      stock: 5,
      description: "",
      publishedDate: "2024-05-01T00:00:00Z",
      createdAt: "2024-05-01T00:00:00Z",
      updatedAt: "2024-05-01T00:00:00Z",
    }));

    render(<BooksManager />);

    await waitFor(() => expect(screen.getByText(/no books found/i)).toBeInTheDocument());

    await userEvent.type(screen.getByLabelText(/title/i), "New Book");
    await userEvent.type(screen.getByLabelText(/author/i), "Jane Doe");
    await userEvent.type(screen.getByLabelText(/isbn/i), "9991112223334");
    await userEvent.clear(screen.getByLabelText(/price/i));
    await userEvent.type(screen.getByLabelText(/price/i), "25");
    await userEvent.clear(screen.getByLabelText(/stock/i));
    await userEvent.type(screen.getByLabelText(/stock/i), "5");
    await userEvent.type(screen.getByLabelText(/published date/i), "2024-05-01");

    await userEvent.click(screen.getByRole("button", { name: /create/i }));

    await waitFor(() => expect(api.createBook).toHaveBeenCalled());

    expect(api.createBook).toHaveBeenCalledWith({
      title: "New Book",
      author: "Jane Doe",
      isbn: "9991112223334",
      price: "25",
      stock: "5",
      description: "",
      publishedDate: "2024-05-01",
    });

    await waitFor(() => expect(screen.getByText("New Book")).toBeInTheDocument());
  });
});

