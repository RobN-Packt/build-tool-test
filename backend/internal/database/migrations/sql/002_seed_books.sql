INSERT INTO books (title, author, isbn, price, stock, description, published_date)
VALUES
    ('The Go Programming Language', 'Alan A. A. Donovan', '9780134190440', 45.99, 12,
     'Comprehensive guide to Go programming fundamentals and best practices.', '2015-10-26'),
    ('Clean Code', 'Robert C. Martin', '9780132350884', 39.50, 8,
     'Handbook of agile software craftsmanship with practical guidelines.', '2008-08-01'),
    ('Domain-Driven Design', 'Eric Evans', '9780321125217', 59.99, 5,
     'A foundational text on tackling complexity in software design.', '2003-08-30'),
    ('Designing Data-Intensive Applications', 'Martin Kleppmann', '9781449373320', 54.25, 10,
     'In-depth exploration of modern data systems and architectural principles.', '2017-03-14'),
    ('Working Effectively with Legacy Code', 'Michael Feathers', '9780131177055', 49.00, 6,
     'Strategies for refactoring and improving existing codebases.', '2004-09-22')
ON CONFLICT (isbn) DO NOTHING;

