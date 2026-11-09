package ports

import (
	"context"

	"github.com/George-Alves/desafio-enviepro/backend/internal/funil/domain/entity"
)

// FunilRepository entrega os totais de oportunidades agrupados por etapa.
// Só devolve as etapas que têm ao menos uma oportunidade.
type FunilRepository interface {
	TotaisPorEtapa(ctx context.Context) ([]entity.TotalEtapa, error)
}
