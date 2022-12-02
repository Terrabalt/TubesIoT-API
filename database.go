package main

import (
	"context"
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
	"github.com/rs/zerolog/log"
	sqldblogger "github.com/simukti/sqldb-logger"
	"github.com/simukti/sqldb-logger/logadapter/zerologadapter"
)

type DB struct {
	*sql.DB
}

var insertStatement, getStatement *sql.Stmt

func StartDB(dbInfo string) (DB, error) {
	db, err := sql.Open("sqlite3", dbInfo)
	if err != nil {
		return DB{nil}, err
	}

	loggerAdapter := zerologadapter.New(log.Logger)
	db = sqldblogger.OpenDriver(dbInfo, db.Driver(), loggerAdapter, sqldblogger.WithSQLQueryAsMessage(true))
	err = db.Ping()
	if err != nil {
		return DB{nil}, err
	}

	log.Info().Str("dsn", dbInfo).Msg("Database connected")
	return DB{db}, nil
}

func (db DB) InitDB(ctx context.Context) error {
	_, err := db.ExecContext(ctx, `
	CREATE TABLE IF NOT EXISTS monitor (
		isOn BOOLEAN,
		suhu DECIMAL
		)`)
	if err != nil {
		return err
	}
	insertStatement, err = db.PrepareContext(ctx, `
	INSERT INTO monitor (
		isOn, suhu
	)
	VALUES
		(?, ?)`)
	if err != nil {
		return err
	}
	getStatement, err = db.PrepareContext(ctx, `
	SELECT
		isOn, 
		suhu
	FROM
		monitor
	LIMIT
		?`)
	return err
}

func (db DB) GetData(ctx context.Context) ([]Status, error) {
	var m []Status

	rows, err := getStatement.QueryContext(ctx, -20)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		d := Status{}
		if err := rows.Scan(
			&d.IsOn,
			&d.Suhu,
		); err != nil {
			return nil, err
		}
		m = append(m, d)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return m, nil
}
func (db DB) AddData(ctx context.Context, isOn bool, suhu float32) (int64, error) {
	res, err := insertStatement.ExecContext(ctx, isOn, suhu)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}
