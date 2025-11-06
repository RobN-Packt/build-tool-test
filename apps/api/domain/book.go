package domain

import "time"

type Book struct {
    ID        string    `json:"id"`
    Title     string    `json:"title"`
    Author    string    `json:"author"`
    Price     float64   `json:"price"`
    Currency  string    `json:"currency"`
    Stock     int       `json:"stock"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}

type CreateBookInput struct {
    Title    string  `json:"title"`
    Author   string  `json:"author"`
    Price    float64 `json:"price"`
    Currency string  `json:"currency"`
    Stock    int     `json:"stock"`
}

type UpdateBookInput struct {
    Title    string  `json:"title"`
    Author   string  `json:"author"`
    Price    float64 `json:"price"`
    Currency string  `json:"currency"`
    Stock    int     `json:"stock"`
}

type PurchaseInput struct {
    BookID     string `json:"book_id"`
    Quantity   int    `json:"quantity"`
    CustomerID string `json:"customer_id"`
}
