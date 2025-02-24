-- Create Library table
CREATE TABLE libraries (
    id SERIAL PRIMARY KEY,
    name TEXT UNIQUE NOT NULL
);

-- Create Users table
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    email TEXT UNIQUE NOT NULL,
    contact_number TEXT,
    role TEXT NOT NULL,
    library_id INT REFERENCES libraries(id)
);

-- Create BookInventory table
CREATE TABLE book_inventory (
    isbn TEXT PRIMARY KEY,
    library_id INT REFERENCES libraries(id),
    title TEXT NOT NULL,
    authors TEXT NOT NULL,
    publisher TEXT,
    version TEXT,
    total_copies INT NOT NULL,
    available_copies INT NOT NULL
);

-- Create RequestEvents table
CREATE TABLE request_events (
    req_id SERIAL PRIMARY KEY,
    book_id TEXT REFERENCES book_inventory(isbn),
    reader_id INT REFERENCES users(id),
    request_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    approval_date TIMESTAMP,
    approver_id INT REFERENCES users(id),
    request_type TEXT
);

-- Create IssueRegistery table
CREATE TABLE issue_registery (
    issue_id SERIAL PRIMARY KEY,
    isbn TEXT REFERENCES book_inventory(isbn),
    reader_id INT REFERENCES users(id),
    issue_approver_id INT REFERENCES users(id),
    issue_status TEXT,
    issue_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    expected_return_date TIMESTAMP,
    return_date TIMESTAMP,
    return_approver_id INT REFERENCES users(id)
);

