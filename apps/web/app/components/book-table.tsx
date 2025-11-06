"use client";

import Link from "next/link";

type Book = {
  id: string;
  title: string;
  author: string;
  price: number;
  currency: string;
  stock: number;
};

interface Props {
  books: Book[];
}

export function BookTable({ books }: Props) {
  return (
    <div className="rounded border border-slate-200 bg-white">
      <table className="w-full table-auto border-collapse">
        <thead className="bg-slate-100 text-left text-sm">
          <tr>
            <th className="px-4 py-3">Title</th>
            <th className="px-4 py-3">Author</th>
            <th className="px-4 py-3">Price</th>
            <th className="px-4 py-3">Stock</th>
          </tr>
        </thead>
        <tbody>
          {books.map((book) => (
            <tr key={book.id} className="border-t border-slate-200 text-sm">
              <td className="px-4 py-3">
                <Link className="text-blue-600" href={`/books/${book.id}`}>
                  {book.title}
                </Link>
              </td>
              <td className="px-4 py-3">{book.author}</td>
              <td className="px-4 py-3">
                {new Intl.NumberFormat("en-US", { style: "currency", currency: book.currency }).format(book.price)}
              </td>
              <td className="px-4 py-3">{book.stock}</td>
            </tr>
          ))}
          {books.length === 0 && (
            <tr>
              <td className="px-4 py-8 text-center text-slate-500" colSpan={4}>
                No books yet. Add one from the admin dashboard.
              </td>
            </tr>
          )}
        </tbody>
      </table>
    </div>
  );
}
