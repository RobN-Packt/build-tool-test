package books

import "time"

// Book represents a record in the books catalog.
type Book struct {
	ID            int64     `json:"id" db:"id"`
	Title         string    `json:"title" db:"title"`
	Author        string    `json:"author" db:"author"`
	ISBN          string    `json:"isbn" db:"isbn"`
	Price         float64   `json:"price" db:"price"`
	Stock         int       `json:"stock" db:"stock"`
	Description   string    `json:"description" db:"description"`
	PublishedDate time.Time `json:"publishedDate" db:"published_date"`
	CreatedAt     time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt     time.Time `json:"updatedAt" db:"updated_at"`
}

// CreateBookInput represents the payload required to create a book.
type CreateBookInput struct {
	Title         string  `json:"title"`
	Author        string  `json:"author"`
	ISBN          string  `json:"isbn"`
	Price         float64 `json:"price"`
	Stock         int     `json:"stock"`
	Description   string  `json:"description"`
	PublishedDate string  `json:"publishedDate"`
}

// UpdateBookInput represents the payload required to update a book.
type UpdateBookInput struct {
	Title         string  `json:"title"`
	Author        string  `json:"author"`
	ISBN          string  `json:"isbn"`
	Price         float64 `json:"price"`
	Stock         int     `json:"stock"`
	Description   string  `json:"description"`
	PublishedDate string  `json:"publishedDate"`
}
