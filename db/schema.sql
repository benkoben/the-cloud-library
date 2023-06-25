CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    email TEXT UNIQUE NOT NULL,
    password TEXT NOT NULL
);

CREATE TABLE books (
    id SERIAL PRIMARY KEY,
    isbn TEXT UNIQUE NOT NULL,
    title TEXT NOT NULL,
    lang TEXT NOT NULL,
    translator TEXT,
    author TEXT NOT NULL,
    pages INTEGER NOT NULL,
    publisher TEXT NOT NULL,
    published_date DATE NOT NULL,
    added_date DATE NOT NULL
);

CREATE TABLE rentals (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(id) on DELETE CASCADE,
    book_id INTEGER REFERENCES books(id) on DELETE CASCADE,
    rental_date DATE NOT null,
    return_date DATE
);

INSERT INTO books (id, isbn, title, lang, translator, author, pages, publisher, published_date, added_date)
VALUES (1, '9789100187934', 'Pesten', 'english', 'Jan Stolpe', 'Albert Camus', 254, 'Albert Bonniers Förlag', '2021-01-07', '2023-06-03');
