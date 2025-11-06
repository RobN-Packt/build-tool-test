import "./globals.css";
import type { Metadata } from "next";

export const metadata: Metadata = {
  title: "Book Store PoC",
  description: "Evaluate monorepo build tooling with a book store",
};

export default function RootLayout({ children }: { children: React.ReactNode }) {
  return (
    <html lang="en">
      <body className="bg-slate-50 text-slate-900">
        <header className="border-b border-slate-200 bg-white">
          <nav className="mx-auto flex max-w-5xl items-center justify-between px-6 py-4">
            <h1 className="text-xl font-semibold">Book Store PoC</h1>
            <a className="text-sm text-blue-600" href="/admin/books">
              Admin
            </a>
          </nav>
        </header>
        <main className="mx-auto min-h-screen max-w-5xl px-6 py-8">{children}</main>
      </body>
    </html>
  );
}
