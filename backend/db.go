package main

import (
	"database/sql"
	"fmt"

	_ "modernc.org/sqlite" // driver SQLite em Go puro, sem precisar de compilador C
)

const esquema = `
CREATE TABLE IF NOT EXISTS clientes (
	id             INTEGER PRIMARY KEY AUTOINCREMENT,
	nome           TEXT NOT NULL,
	empresa        TEXT NOT NULL,
	email          TEXT NOT NULL,
	telefone       TEXT NOT NULL,
	status         TEXT NOT NULL,
	responsavel    TEXT NOT NULL,
	ultimo_contato TEXT NOT NULL
);
CREATE TABLE IF NOT EXISTS oportunidades (
	id           INTEGER PRIMARY KEY AUTOINCREMENT,
	cliente_id   INTEGER NOT NULL REFERENCES clientes(id),
	titulo       TEXT NOT NULL,
	etapa_funil  TEXT NOT NULL,
	valor        REAL NOT NULL,
	data_criacao TEXT NOT NULL
);
CREATE TABLE IF NOT EXISTS vendas (
	id         INTEGER PRIMARY KEY AUTOINCREMENT,
	cliente_id INTEGER NOT NULL REFERENCES clientes(id),
	valor      REAL NOT NULL,
	data       TEXT NOT NULL,
	produto    TEXT NOT NULL
);
`

// abrirBanco abre (ou cria) o arquivo SQLite e garante que as tabelas existem.
func abrirBanco(caminho string) (*sql.DB, error) {
	db, err := sql.Open("sqlite", caminho)
	if err != nil {
		return nil, err
	}
	// Uma única conexão: o SQLite grava um por vez, e assim o PRAGMA vale para todas as consultas.
	db.SetMaxOpenConns(1)
	if _, err := db.Exec("PRAGMA foreign_keys = ON"); err != nil {
		db.Close()
		return nil, fmt.Errorf("ativando chaves estrangeiras: %w", err)
	}
	if _, err := db.Exec(esquema); err != nil {
		db.Close()
		return nil, fmt.Errorf("criando tabelas: %w", err)
	}
	return db, nil
}

// bancoVazio diz se ainda não há nenhum cliente (banco recém-criado).
func bancoVazio(db *sql.DB) (bool, error) {
	var n int
	if err := db.QueryRow("SELECT COUNT(*) FROM clientes").Scan(&n); err != nil {
		return false, err
	}
	return n == 0, nil
}
