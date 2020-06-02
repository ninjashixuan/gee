package dialect

type sqlite3 struct {}

var _ Dialect = (*sqlite3)(nil)

