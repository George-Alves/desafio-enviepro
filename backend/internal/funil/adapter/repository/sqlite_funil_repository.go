package repository

import (
	"context"
	"database/sql"

	"github.com/George-Alves/desafio-enviepro/backend/internal/funil/domain/entity"
)

type SQLiteFunilRepository struct {
	db *sql.DB
}

func NewSQLiteFunilRepository(db *sql.DB) *SQLiteFunilRepository {
	return &SQLiteFunilRepository{db: db}
}

func (r *SQLiteFunilRepository) TotaisPorEtapa(ctx context.Context) ([]entity.TotalEtapa, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT etapa_funil, COUNT(*), COALESCE(SUM(valor), 0)
		FROM oportunidades GROUP BY etapa_funil`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var totais []entity.TotalEtapa
	for rows.Next() {
		var t entity.TotalEtapa
		if err := rows.Scan(&t.Etapa, &t.Quantidade, &t.Valor); err != nil {
			return nil, err
		}
		totais = append(totais, t)
	}
	return totais, rows.Err()
}
