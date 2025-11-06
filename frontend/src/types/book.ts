export type Book = {
  id: number;
  title: string;
  author: string;
  isbn: string;
  price: number;
  stock: number;
  description: string;
  publishedDate: string;
  createdAt: string;
  updatedAt: string;
};

export type BookFormValues = {
  title: string;
  author: string;
  isbn: string;
  price: string;
  stock: string;
  description: string;
  publishedDate: string;
};

