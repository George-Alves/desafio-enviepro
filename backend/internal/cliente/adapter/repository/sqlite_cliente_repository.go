package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/George-Alves/desafio-enviepro/backend/internal/cliente/domain/entity"
)

// FormatoData é como as datas são guardadas no SQLite (texto, sem hora).
const FormatoData = "2006-01-02"

type SQLiteClienteRepository struct {
	db *sql.DB
}

func NewSQLiteClienteRepository(db *sql.DB) *SQLiteClienteRepository {
	return &SQLiteClienteRepository{db: db}
}

func (r *SQLiteClienteRepository) Listar(ctx context.Context) ([]entity.Cliente, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, nome, empresa, email, telefone, status, responsavel, ultimo_contato
		FROM clientes ORDER BY nome`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var clientes []entity.Cliente
	for rows.Next() {
		var c entity.Cliente
		var ultimoContato string
		if err := rows.Scan(&c.ID, &c.Nome, &c.Empresa, &c.Email, &c.Telefone,
			&c.Status, &c.Responsavel, &ultimoContato); err != nil {
			return nil, err
		}
		if c.UltimoContato, err = time.Parse(FormatoData, ultimoContato); err != nil {
			return nil, err
		}
		clientes = append(clientes, c)
	}
	return clientes, rows.Err()
}
