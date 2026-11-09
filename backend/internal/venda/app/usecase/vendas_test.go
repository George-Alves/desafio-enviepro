package usecase

import (
	"context"
	"testing"
	"time"

	"github.com/George-Alves/desafio-enviepro/backend/internal/venda/domain/entity"
)

type repoFalso struct {
	total      float64
	quantidade int
	de, ate    time.Time
}

func (r *repoFalso) Listar(context.Context) ([]entity.Venda, error) { return nil, nil }

func (r *repoFalso) TotalDoPeriodo(_ context.Context, de, ate time.Time) (float64, int, error) {
	r.de, r.ate = de, ate
	return r.total, r.quantidade, nil
}

func TestResumoVendasDoMes_CalculaTicketMedioEConsultaOMesInteiro(t *testing.T) {
	repo := &repoFalso{total: 10000, quantidade: 3}
	agora := time.Date(2026, 12, 31, 23, 59, 0, 0, time.UTC)

	got, err := NewResumoVendasDoMes(repo).Executar(context.Background(), agora)
	if err != nil {
		t.Fatal(err)
	}

	if got.Total != 10000 || got.Quantidade != 3 || got.TicketMedio != 3333.33 {
		t.Errorf("resumo inesperado: %+v", got)
	}
	if !repo.de.Equal(time.Date(2026, 12, 1, 0, 0, 0, 0, time.UTC)) ||
		!repo.ate.Equal(time.Date(2027, 1, 1, 0, 0, 0, 0, time.UTC)) {
		t.Errorf("período errado: de %v até %v", repo.de, repo.ate)
	}
}

func TestResumoVendasDoMes_SemVendasTicketEZero(t *testing.T) {
	got, err := NewResumoVendasDoMes(&repoFalso{}).Executar(context.Background(), time.Now())
	if err != nil {
		t.Fatal(err)
	}
	if got != (Resumo{}) {
		t.Errorf("esperava resumo zerado, veio %+v", got)
	}
}
