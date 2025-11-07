"use client";

import { FormEvent, useCallback, useEffect, useMemo, useState } from "react";

import { createBook, deleteBook, fetchBooks, updateBook } from "@/lib/api";
import { Book, BookFormValues } from "@/types/book";

import styles from "./BooksManager.module.css";

const initialFormValues: BookFormValues = {
  title: "",
  author: "",
  isbn: "",
  price: "0",
  stock: "0",
  description: "",
  publishedDate: "",
};

const formatPublishedDate = (value: string) => {
  if (!value) {
    return "—";
  }

  const isoValue = value.length === 10 ? `${value}T00:00:00Z` : value;
  const date = new Date(isoValue);

  if (Number.isNaN(date.getTime())) {
    return value;
  }

  return date.toLocaleDateString();
};

export function BooksManager() {
  const [books, setBooks] = useState<Book[]>([]);
  const [formValues, setFormValues] = useState<BookFormValues>(initialFormValues);
  const [editingId, setEditingId] = useState<number | null>(null);
  const [loading, setLoading] = useState<boolean>(true);
  const [error, setError] = useState<string | null>(null);
  const [submitting, setSubmitting] = useState<boolean>(false);

  const loadBooks = useCallback(async () => {
    try {
      setLoading(true);
      const data = await fetchBooks();
      setBooks(Array.isArray(data) ? data : []);
      setError(null);
    } catch (err) {
      setError(err instanceof Error ? err.message : "Failed to load books");
    } finally {
      setLoading(false);
    }
  }, []);

  useEffect(() => {
    loadBooks();
  }, [loadBooks]);

  const summary = useMemo(() => {
    if (!Array.isArray(books) || books.length === 0) {
      return { totalBooks: 0, totalStock: 0, inventoryValue: 0 };
    }

    const totals = books.reduce(
      (accumulator, book) => {
        const next = { ...accumulator };
        next.totalStock += book.stock;
        next.inventoryValue += book.price * book.stock;
        return next;
      },
      { totalStock: 0, inventoryValue: 0 },
    );

    return {
      totalBooks: books.length,
      totalStock: totals.totalStock,
      inventoryValue: totals.inventoryValue,
    };
  }, [books]);

  const formatter = useMemo(
    () =>
      new Intl.NumberFormat("en-US", {
        style: "currency",
        currency: "USD",
        minimumFractionDigits: 2,
      }),
    [],
  );

  const handleInputChange = (field: keyof BookFormValues, value: string) => {
    setFormValues((previous: BookFormValues) => ({ ...previous, [field]: value }));
  };

  const resetForm = () => {
    setEditingId(null);
    setFormValues(initialFormValues);
    setError(null);
  };

  const handleSubmit = async (event: FormEvent<HTMLFormElement>) => {
    event.preventDefault();
    setSubmitting(true);
    setError(null);

    try {
      if (editingId) {
        await updateBook(editingId, formValues);
      } else {
        await createBook(formValues);
      }

      await loadBooks();
      resetForm();
    } catch (err) {
      setError(err instanceof Error ? err.message : "Failed to save book");
    } finally {
      setSubmitting(false);
    }
  };

  const handleEdit = (book: Book) => {
    setEditingId(book.id);
    setFormValues({
      title: book.title,
      author: book.author,
      isbn: book.isbn,
      price: book.price.toString(),
      stock: book.stock.toString(),
      description: book.description ?? "",
      publishedDate: book.publishedDate.slice(0, 10),
    });
  };

  const handleDelete = async (id: number) => {
    setError(null);
    try {
      await deleteBook(id);
      await loadBooks();

      if (editingId === id) {
        resetForm();
      }
    } catch (err) {
      setError(err instanceof Error ? err.message : "Failed to delete book");
    }
  };

  return (
    <div className={styles.container}>
      <header className={styles.header}>
        <div>
          <h1 className={styles.title}>Bookshop Manager</h1>
          <p className={styles.subtitle}>
            Manage inventory, pricing, and catalogue entries with live updates.
          </p>
        </div>
        <dl className={styles.summary} data-testid="inventory-summary">
          <div>
            <dt>Books</dt>
            <dd>{summary.totalBooks}</dd>
          </div>
          <div>
            <dt>Units in Stock</dt>
            <dd>{summary.totalStock}</dd>
          </div>
          <div>
            <dt>Inventory Value</dt>
            <dd>{formatter.format(summary.inventoryValue)}</dd>
          </div>
        </dl>
      </header>

      <section className={styles.layout}>
        <form className={styles.form} onSubmit={handleSubmit} data-testid="book-form">
          <h2>{editingId ? "Update Book" : "Add New Book"}</h2>

          <div className={styles.fieldGroup}>
            <label htmlFor="title">Title *</label>
            <input
              id="title"
              name="title"
              value={formValues.title}
              onChange={(event) => handleInputChange("title", event.target.value)}
              required
            />
          </div>

          <div className={styles.fieldGroup}>
            <label htmlFor="author">Author *</label>
            <input
              id="author"
              name="author"
              value={formValues.author}
              onChange={(event) => handleInputChange("author", event.target.value)}
              required
            />
          </div>

          <div className={styles.fieldGrid}>
            <div className={styles.fieldGroup}>
              <label htmlFor="isbn">ISBN *</label>
              <input
                id="isbn"
                name="isbn"
                value={formValues.isbn}
                onChange={(event) => handleInputChange("isbn", event.target.value)}
                required
              />
            </div>

            <div className={styles.fieldGroup}>
              <label htmlFor="price">Price *</label>
              <input
                id="price"
                name="price"
                type="number"
                min="0"
                step="0.01"
                value={formValues.price}
                onChange={(event) => handleInputChange("price", event.target.value)}
                required
              />
            </div>

            <div className={styles.fieldGroup}>
              <label htmlFor="stock">Stock *</label>
              <input
                id="stock"
                name="stock"
                type="number"
                min="0"
                value={formValues.stock}
                onChange={(event) => handleInputChange("stock", event.target.value)}
                required
              />
            </div>
          </div>

          <div className={styles.fieldGroup}>
            <label htmlFor="publishedDate">Published Date *</label>
            <input
              id="publishedDate"
              name="publishedDate"
              type="date"
              value={formValues.publishedDate}
              onChange={(event) => handleInputChange("publishedDate", event.target.value)}
              required
            />
          </div>

          <div className={styles.fieldGroup}>
            <label htmlFor="description">Description</label>
            <textarea
              id="description"
              name="description"
              rows={3}
              value={formValues.description}
              onChange={(event) => handleInputChange("description", event.target.value)}
            />
          </div>

          <div className={styles.formActions}>
            <button type="submit" disabled={submitting}>
              {submitting ? "Saving..." : editingId ? "Update" : "Create"}
            </button>
            <button type="button" onClick={resetForm} disabled={submitting}>
              Reset
            </button>
          </div>

          {error && (
            <p role="alert" className={styles.error}>
              {error}
            </p>
          )}
        </form>

        <section className={styles.listSection}>
          <header className={styles.listHeader}>
            <h2>Catalogue</h2>
            <p>{loading ? "Refreshing inventory..." : `${books.length} books available`}</p>
          </header>

            {loading ? (
            <p className={styles.placeholder}>Loading books...</p>
          ) : books.length === 0 ? (
            <p className={styles.placeholder}>No books found. Add your first entry!</p>
          ) : (
            <table className={styles.table} data-testid="books-table">
              <thead>
                <tr>
                  <th>Title</th>
                  <th>Author</th>
                  <th>ISBN</th>
                  <th>Price</th>
                  <th>Stock</th>
                  <th>Published</th>
                  <th>Actions</th>
                </tr>
              </thead>
              <tbody>
                {books.map((book) => (
                  <tr key={book.id}>
                    <td>{book.title}</td>
                    <td>{book.author}</td>
                    <td>{book.isbn}</td>
                    <td>{formatter.format(book.price)}</td>
                    <td>{book.stock}</td>
                    <td>{formatPublishedDate(book.publishedDate)}</td>
                    <td className={styles.rowActions}>
                      <button type="button" onClick={() => handleEdit(book)}>
                        Edit
                      </button>
                      <button type="button" onClick={() => handleDelete(book.id)}>
                        Delete
                      </button>
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          )}
        </section>
      </section>
    </div>
  );
}

