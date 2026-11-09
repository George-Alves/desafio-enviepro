package usecase

import (
	"context"

	"github.com/George-Alves/desafio-enviepro/backend/internal/funil/domain/entity"
	"github.com/George-Alves/desafio-enviepro/backend/internal/funil/domain/ports"
)

type ObterFunil struct {
	repo ports.FunilRepository
}

func NewObterFunil(repo ports.FunilRepository) *ObterFunil {
	return &ObterFunil{repo: repo}
}

// Executar devolve todas as etapas na ordem do funil, com zero nas que não têm oportunidades.
// Etapas desconhecidas vindas do repositório são ignoradas.
func (u *ObterFunil) Executar(ctx context.Context) ([]entity.TotalEtapa, error) {
	totais, err := u.repo.TotaisPorEtapa(ctx)
	if err != nil {
		return nil, err
	}
	porEtapa := make(map[string]entity.TotalEtapa, len(totais))
	for _, t := range totais {
		porEtapa[t.Etapa] = t
	}
	funil := make([]entity.TotalEtapa, 0, len(entity.Etapas))
	for _, etapa := range entity.Etapas {
		t := porEtapa[etapa]
		t.Etapa = etapa
		funil = append(funil, t)
	}
	return funil, nil
}
