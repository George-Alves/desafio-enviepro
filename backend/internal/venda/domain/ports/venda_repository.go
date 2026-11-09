package ports

import (
	"context"
	"time"

	"github.com/George-Alves/desafio-enviepro/backend/internal/venda/domain/entity"
)

// VendaRepository é o que a regra de negócio precisa para ler vendas.
type VendaRepository interface {
	// Listar devolve as vendas da mais recente para a mais antiga.
	Listar(ctx context.Context) ([]entity.Venda, error)
	// TotalDoPeriodo soma as vendas com data em [de, ate), ignorando a hora.
	TotalDoPeriodo(ctx context.Context, de, ate time.Time) (total float64, quantidade int, err error)
}
