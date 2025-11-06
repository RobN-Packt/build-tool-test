import { render, screen } from "@testing-library/react";

import { BookTable } from "@/app/components/book-table";

describe("BookTable", () => {
  it("renders empty state", () => {
    render(<BookTable books={[]} />);
    expect(screen.getByText(/no books yet/i)).toBeInTheDocument();
  });

  it("renders rows", () => {
    render(
      <BookTable
        books={[
          {
            id: "1",
            title: "Domain-Driven Design",
            author: "Eric Evans",
            price: 50,
            currency: "USD",
            stock: 5,
          },
        ]}
      />,
    );
    expect(screen.getByText(/Domain-Driven Design/)).toBeInTheDocument();
    expect(screen.getByRole("link", { name: /Domain-Driven Design/ })).toHaveAttribute("href", "/books/1");
  });
});
