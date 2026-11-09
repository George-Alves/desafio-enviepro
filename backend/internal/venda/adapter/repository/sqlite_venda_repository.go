package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/George-Alves/desafio-enviepro/backend/internal/venda/domain/entity"
)

// FormatoData é como as datas são guardadas no SQLite (texto, sem hora).
const FormatoData = "2006-01-02"

type SQLiteVendaRepository struct {
	db *sql.DB
}

func NewSQLiteVendaRepository(db *sql.DB) *SQLiteVendaRepository {
	return &SQLiteVendaRepository{db: db}
}

func (r *SQLiteVendaRepository) Listar(ctx context.Context) ([]entity.Venda, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, cliente_id, valor, data, produto
		FROM vendas ORDER BY data DESC, id DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var vendas []entity.Venda
	for rows.Next() {
		var v entity.Venda
		var data string
		if err := rows.Scan(&v.ID, &v.ClienteID, &v.Valor, &data, &v.Produto); err != nil {
			return nil, err
		}
		if v.Data, err = time.Parse(FormatoData, data); err != nil {
			return nil, err
		}
		vendas = append(vendas, v)
	}
	return vendas, rows.Err()
}

func (r *SQLiteVendaRepository) TotalDoPeriodo(ctx context.Context, de, ate time.Time) (float64, int, error) {
	var total float64
	var quantidade int
	err := r.db.QueryRowContext(ctx, `
		SELECT COALESCE(SUM(valor), 0), COUNT(*)
		FROM vendas WHERE data >= ? AND data < ?`,
		de.Format(FormatoData), ate.Format(FormatoData)).Scan(&total, &quantidade)
	return total, quantidade, err
}
