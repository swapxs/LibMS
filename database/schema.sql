-- Create Library table
CREATE TABLE libraries (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) UNIQUE NOT NULL
);

-- Create Users table
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    contact_number VARCHAR(20),
    role VARCHAR(50) NOT NULL,
    library_id INT REFERENCES libraries(id)
);

-- Create BookInventory table
CREATE TABLE book_inventory (
    isbn VARCHAR(20) PRIMARY KEY,
    library_id INT REFERENCES libraries(id),
    title VARCHAR(255) NOT NULL,
    authors TEXT NOT NULL,
    publisher VARCHAR(255),
    version VARCHAR(50),
    total_copies INT NOT NULL,
    available_copies INT NOT NULL
);

-- Create RequestEvents table
CREATE TABLE request_events (
    req_id SERIAL PRIMARY KEY,
    book_id VARCHAR(20) REFERENCES book_inventory(isbn),
    reader_id INT REFERENCES users(id),
    request_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    approval_date TIMESTAMP,
    approver_id INT REFERENCES users(id),
    request_type VARCHAR(50)
);

-- Create IssueRegistery table
CREATE TABLE issue_registery (
    issue_id SERIAL PRIMARY KEY,
    isbn VARCHAR(20) REFERENCES book_inventory(isbn),
    reader_id INT REFERENCES users(id),
    issue_approver_id INT REFERENCES users(id),
    issue_status VARCHAR(50),
    issue_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    expected_return_date TIMESTAMP,
    return_date TIMESTAMP,
    return_approver_id INT REFERENCES users(id)
);

