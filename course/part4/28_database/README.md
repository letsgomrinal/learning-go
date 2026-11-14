# Lesson 28: Databases

**Goal:** store and query data in a real SQL database with `database/sql`: open a connection, run
queries safely with placeholders, read rows into structs, handle "not found", and use transactions.

**Time:** about 2 hours. You should know basic SQL (`SELECT`, `INSERT`, `UPDATE`, `CREATE TABLE`).

```bash
cd ~/Desktop/learning-go/course/part4/28_database
```

---

## 0. This lesson is its own module

Look in this folder: it has its **own `go.mod`**. That makes it a separate module from the rest of the
course (lesson 09), with one third-party dependency: **`modernc.org/sqlite`**, a SQLite database
written in pure Go. No database server to install, and no C compiler needed. It stores everything in
a single file.

Open `go.mod`: the `require` lines list the dependency, and `go.sum` records checksums so everyone gets
exactly the same code. This is how you'd add it to your own project:

```bash
go get modernc.org/sqlite
```

(It's already done here. The first `go test` downloads the module, so you need an internet connection once.)

The same `database/sql` code works with PostgreSQL or MySQL. You only change the driver and the
connection string. Section 9 shows that.

---

## 1. `database/sql` and drivers

Go's standard library includes `database/sql`, a common interface for SQL databases. It knows **how to
talk to databases in general**, but not to any specific one. A **driver** package adds support for one
database. You import the driver for its side effect only (it registers itself), using `_`:

```go
import (
	"database/sql"

	_ "modernc.org/sqlite"      // registers the "sqlite" driver; we never call it directly
)
```

The blank import `_` means "run this package's setup code, but I won't use its name". You'll see it
with every database driver.

Python comparison: `database/sql` is like the DB-API (`sqlite3`, `psycopg`), and connection pooling is
built in.

---

## 2. Opening a database

```go
db, err := sql.Open("sqlite", "todos.db")     // driver name, data source
if err != nil {
	return err
}
defer db.Close()

if err := db.Ping(); err != nil {             // actually connect, to check it works
	return err
}
```

Important facts about `*sql.DB`:
- It's **not one connection**. It's a **pool** of connections that Go manages for you, opening and
  reusing them as needed. It's safe to use from many goroutines (every HTTP request) at once.
- **Open it once** when your program starts, and pass it around. Don't open and close it per query.
- `sql.Open` doesn't connect yet; `Ping` checks the connection works.

⚠️ SQLite special case: `":memory:"` gives an in-memory database, but **each connection in the pool
gets its own, separate, empty one**. Use a real file (the tests use one in `t.TempDir()`).

---

## 3. Running statements: `ExecContext`

For statements that don't return rows (`CREATE`, `INSERT`, `UPDATE`, `DELETE`):

```go
res, err := db.ExecContext(ctx,
	"INSERT INTO todos (title, done) VALUES (?, ?)",
	"buy milk", false)
if err != nil {
	return err
}
id, err := res.LastInsertId()       // the new row's id
n, err := res.RowsAffected()        // how many rows changed (useful for UPDATE / DELETE)
```

Every method has a `...Context` version that takes a `ctx` (lesson 22). Use those, so slow queries get
cancelled when a request is cancelled.

### ⚠️ Placeholders, never string formatting

The `?` marks are **placeholders**. The values are sent separately from the SQL, so the database never
confuses data with SQL code. **Never build SQL with `fmt.Sprintf` or `+`:**

```go
// ❌ SQL INJECTION: if title is  x'); DROP TABLE todos; --  you lose everything
db.ExecContext(ctx, "INSERT INTO todos (title) VALUES ('" + title + "')")

// ✅ safe, always
db.ExecContext(ctx, "INSERT INTO todos (title) VALUES (?)", title)
```

Same rule as Python's `cursor.execute("... ?", (title,))`. (PostgreSQL drivers use `$1, $2` instead
of `?`.)

---

## 4. One row: `QueryRowContext`

```go
var t Todo
err := db.QueryRowContext(ctx,
	"SELECT id, title, done FROM todos WHERE id = ?", id,
).Scan(&t.ID, &t.Title, &t.Done)
```

`Scan` copies the columns, **in order**, into the variables you point it at (`&`, lesson 08). The number
and order of `&` arguments must match the `SELECT` columns.

### Not found: `sql.ErrNoRows`

If no row matches, `Scan` returns the sentinel error **`sql.ErrNoRows`** (lesson 12):

```go
if errors.Is(err, sql.ErrNoRows) {
	return Todo{}, ErrNotFound          // translate it into your own error
}
if err != nil {
	return Todo{}, err
}
```

Translating `sql.ErrNoRows` into your own `ErrNotFound` keeps database details out of the rest of the
program. An HTTP handler can then turn `ErrNotFound` into a 404 (lesson 13's `StatusCode`).

---

## 5. Many rows: `QueryContext`

```go
rows, err := db.QueryContext(ctx, "SELECT id, title, done FROM todos ORDER BY id")
if err != nil {
	return nil, err
}
defer rows.Close()                       // ⚠️ always, or you leak a connection

todos := []Todo{}
for rows.Next() {                        // move to the next row; false at the end
	var t Todo
	if err := rows.Scan(&t.ID, &t.Title, &t.Done); err != nil {
		return nil, err
	}
	todos = append(todos, t)
}
return todos, rows.Err()                 // an error that stopped the loop early
```

The shape is like `bufio.Scanner` from lesson 23: `Next()` in a loop, then check the error afterwards.
Every `QueryContext` needs its `defer rows.Close()`.

---

## 6. Transactions: all or nothing

Sometimes several statements must **all** succeed or **all** be undone. Moving money between two
accounts is the classic example. That's a **transaction**:

```go
tx, err := db.BeginTx(ctx, nil)
if err != nil {
	return err
}
defer tx.Rollback()      // undo everything, unless Commit happened first (then it does nothing)

if _, err := tx.ExecContext(ctx, "UPDATE accounts SET balance = balance - ? WHERE id = ?", 50, from); err != nil {
	return err           // the deferred Rollback undoes the changes
}
if _, err := tx.ExecContext(ctx, "UPDATE accounts SET balance = balance + ? WHERE id = ?", 50, to); err != nil {
	return err
}
return tx.Commit()       // make it permanent
```

The `defer tx.Rollback()` pattern means you can't forget to undo on an error path: every early `return`
rolls back, and after a successful `Commit` the rollback is a harmless no-op.

Inside a transaction, use `tx.ExecContext` / `tx.QueryContext`, not `db....`. Otherwise the statement
runs outside the transaction.

---

## 7. NULL values

SQL `NULL` can't go into a plain Go `string` or `int` (they always have a value). If a column can be
NULL, scan into `sql.NullString`, `sql.NullInt64`, etc., or a pointer (`*string`), and check `.Valid`.
The easiest fix is often to declare columns `NOT NULL` with a default, as this lesson's table does.

---

## 8. A repository type

Real programs wrap the database in a type with methods, so the rest of the code calls
`store.Add(ctx, "buy milk")` and never sees SQL:

```go
type Store struct {
	db *sql.DB
}

func (s *Store) Add(ctx context.Context, title string) (int64, error) { ... }
func (s *Store) Get(ctx context.Context, id int64) (Todo, error)      { ... }
```

That's what you'll build. Later you could replace the HTTP server's in-memory map from lesson 26 with
this store, and the notes would survive a restart.

### 🤔 Predict the output

```go
res, _ := db.ExecContext(ctx, "UPDATE todos SET done = 1 WHERE id = ?", 999)
n, _ := res.RowsAffected()
fmt.Println(n)
```

(There's no todo with id 999.)

<details markdown="1">
<summary>Answer</summary>

`0`. An `UPDATE` that matches nothing is **not** an error in SQL. It just changes zero rows. If "no such
todo" matters, check `RowsAffected()` yourself and return `ErrNotFound`. You'll do that in Ex5.
</details>

---

## 9. Other databases

Switching to PostgreSQL is mostly a different driver and connection string:

```go
import _ "github.com/jackc/pgx/v5/stdlib"

db, err := sql.Open("pgx", "postgres://user:pass@localhost:5432/mydb")
// placeholders become $1, $2, ... instead of ?
```

Popular tools you'll meet later: **`sqlc`** (write SQL, it generates type-safe Go code), **`pgx`**'s own
API for PostgreSQL-specific features, and migration tools like **`goose`** or **`golang-migrate`** for
versioning your schema.

---

## 10. Common mistakes

| You wrote | What happens | Fix |
|---|---|---|
| SQL built with `+` or `Sprintf` from user input | SQL injection | `?` placeholders |
| Forgot `defer rows.Close()` | connections leak until the pool is exhausted | close right after checking `err` |
| Opening the database for every query | slow, defeats the pool | open once at startup |
| Didn't check `rows.Err()` | a failure halfway looks like a short list | `return list, rows.Err()` |
| `Scan` with columns in a different order than `SELECT` | wrong data in fields, or type errors | match the order exactly |
| `db.ExecContext` inside a transaction | runs outside it | `tx.ExecContext` |
| `":memory:"` SQLite with a pool | each connection sees a different empty database | use a file |

---

## 11. Exercises: a todo store

| | Exercise | Practises |
|---|---|---|
| Ex1 | `Open`: connect and create the table | `sql.Open`, `Ping`, `ExecContext` |
| Ex2 | `Add` a todo | `INSERT` with placeholders, `LastInsertId` |
| Ex3 | `Get` one todo | `QueryRowContext`, `Scan`, `sql.ErrNoRows` |
| Ex4 | `List` all todos | `QueryContext`, `rows.Next`, `rows.Err` |
| Ex5 | `MarkDone` | `UPDATE`, `RowsAffected` |
| Ex6 | `Search` by text | `LIKE` with a placeholder (and an injection test!) |
| Ex7 | `AddAll`: all or nothing | transactions |

---

## ✅ Checklist

- [ ] I know `*sql.DB` is a pool: open once, share it
- [ ] I always use placeholders, never string building, for values in SQL
- [ ] I can read one row with `QueryRowContext().Scan` and handle `sql.ErrNoRows`
- [ ] I can loop over rows with `defer rows.Close()` and `rows.Err()`
- [ ] I can use a transaction with `defer tx.Rollback()` and `tx.Commit()`

**Next → [Lesson 29: Testing in depth](../29_testing/README.md)**
