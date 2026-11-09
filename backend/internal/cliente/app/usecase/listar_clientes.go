package usecase

import (
	"context"
	"time"

	"github.com/George-Alves/desafio-enviepro/backend/internal/cliente/domain/entity"
	"github.com/George-Alves/desafio-enviepro/backend/internal/cliente/domain/ports"
)

// LimiteDiasSemContato: acima disso o cliente é sinalizado como sem acompanhamento.
const LimiteDiasSemContato = 15

// ClienteAcompanhado é o cliente junto com a situação do acompanhamento.
type ClienteAcompanhado struct {
	entity.Cliente
	DiasSemContato int
	Alerta         bool
}

type ListarClientes struct {
	repo ports.ClienteRepository
}

func NewListarClientes(repo ports.ClienteRepository) *ListarClientes {
	return &ListarClientes{repo: repo}
}

// Executar devolve os clientes com os dias sem contato calculados em relação a agora.
func (u *ListarClientes) Executar(ctx context.Context, agora time.Time) ([]ClienteAcompanhado, error) {
	clientes, err := u.repo.Listar(ctx)
	if err != nil {
		return nil, err
	}
	resultado := make([]ClienteAcompanhado, 0, len(clientes))
	for _, c := range clientes {
		dias := diasEntre(c.UltimoContato, agora)
		resultado = append(resultado, ClienteAcompanhado{
			Cliente:        c,
			DiasSemContato: dias,
			Alerta:         dias > LimiteDiasSemContato,
		})
	}
	return resultado, nil
}

// diasEntre conta dias de calendário, ignorando a hora.
func diasEntre(de, ate time.Time) int {
	dia := func(t time.Time) time.Time {
		a, m, d := t.Date()
		return time.Date(a, m, d, 0, 0, 0, 0, time.UTC)
	}
	return int(dia(ate).Sub(dia(de)).Hours() / 24)
}
