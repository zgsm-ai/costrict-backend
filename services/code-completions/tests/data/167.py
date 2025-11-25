import sqlite3
from contextlib import contextmanager

@contextmanager
def db_connection(db_path):
    conn = sqlite3.connect(db_path)
    try:
        yield conn
    finally:
        conn.close()

def create_tables(conn):
    cursor = conn.cursor()
    
    cursor.execute('''
        CREATE TABLE IF NOT EXISTS users (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            name TEXT NOT NULL,
            email TEXT UNIQUE NOT NULL
        )
    ''')
    
    cursor.execute('''
        CREATE TABLE IF NOT EXISTS orders (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            user_id INTEGER NOT NULL,
            product TEXT NOT NULL,
            amount REAL NOT NULL,
            FOREIGN KEY (user_id) REFERENCES users (id)
        )
    ''')
    
    conn.commit()

def insert_sample_data(conn):
    cursor = conn.cursor()
    
    users = [('Alice', 'alice@example.com'), ('Bob', 'bob@example.com')]
    cursor.executemany('INSERT INTO users (name, email) VALUES (?, ?)', users)
    
    cursor.execute('SELECT id, name FROM users')
    users_data = cursor.fetchall()
    
    orders = []
    for user_id, name in users_data:
        <｜fim▁hole｜>if name == 'Alice':
            orders.append((user_id, 'Laptop', 1200.0))
        elif name == 'Bob':
            orders.append((user_id, 'Keyboard', 75.0))
    
    cursor.executemany('INSERT INTO orders (user_id, product, amount) VALUES (?, ?, ?)', orders)
    conn.commit()

def main():
    db_path = 'example.db'
    
    with db_connection(db_path) as conn:
        create_tables(conn)
        insert_sample_data(conn)
        
        cursor = conn.cursor()
        cursor.execute('''
            SELECT u.name, o.product, o.amount
            FROM users u
            JOIN orders o ON u.id = o.user_id
        ''')
        
        for row in cursor.fetchall():
            print(f"{row[0]} bought {row[1]} for ${row[2]}")

if __name__ == "__main__":
    main()