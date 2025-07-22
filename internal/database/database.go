package database

import "database/sql"

type TableInitializer func(database *sql.DB, channel chan bool)
