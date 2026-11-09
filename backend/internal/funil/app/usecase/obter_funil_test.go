package usecase

import (
	"context"
	"testing"

	"github.com/George-Alves/desafio-enviepro/backend/internal/funil/domain/entity"
)

type repoFalso struct{ totais []entity.TotalEtapa }

func (r repoFalso) TotaisPorEtapa(context.Context) ([]entity.TotalEtapa, error) {
	return r.totais, nil
}

func TestObterFunil_TodasAsEtapasEmOrdemComZeroNasVazias(t *testing.T) {
	repo := repoFalso{totais: []entity.TotalEtapa{
		{Etapa: "Proposta", Quantidade: 2, Valor: 5000},
		{Etapa: "Lead", Quantidade: 3, Valor: 900},
		{Etapa: "Etapa inventada", Quantidade: 9, Valor: 1},
	}}

	got, err := NewObterFunil(repo).Executar(context.Background())
	if err != nil {
		t.Fatal(err)
	}

	esperado := []entity.TotalEtapa{
		{Etapa: "Lead", Quantidade: 3, Valor: 900},
		{Etapa: "Contato"},
		{Etapa: "Proposta", Quantidade: 2, Valor: 5000},
		{Etapa: "Negociação"},
		{Etapa: "Fechado"},
	}
	if len(got) != len(esperado) {
		t.Fatalf("esperava %d etapas, vieram %d: %+v", len(esperado), len(got), got)
	}
	for i := range esperado {
		if got[i] != esperado[i] {
			t.Errorf("etapa %d: veio %+v, esperado %+v", i, got[i], esperado[i])
		}
	}
}
