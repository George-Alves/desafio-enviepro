package ports

import (
	"context"

	"github.com/George-Alves/desafio-enviepro/backend/internal/cliente/domain/entity"
)

// ClienteRepository é o que a regra de negócio precisa para ler clientes.
type ClienteRepository interface {
	Listar(ctx context.Context) ([]entity.Cliente, error)
}
