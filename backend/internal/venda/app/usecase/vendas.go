package usecase

import (
	"context"
	"math"
	"time"

	"github.com/George-Alves/desafio-enviepro/backend/internal/venda/domain/entity"
	"github.com/George-Alves/desafio-enviepro/backend/internal/venda/domain/ports"
)

// Resumo são os números de vendas de um mês.
type Resumo struct {
	Total       float64
	Quantidade  int
	TicketMedio float64
}

type ResumoVendasDoMes struct {
	repo ports.VendaRepository
}

func NewResumoVendasDoMes(repo ports.VendaRepository) *ResumoVendasDoMes {
	return &ResumoVendasDoMes{repo: repo}
}

// Executar resume o mês de "agora". Sem vendas, o ticket médio é zero.
func (u *ResumoVendasDoMes) Executar(ctx context.Context, agora time.Time) (Resumo, error) {
	ano, mes, _ := agora.Date()
	inicio := time.Date(ano, mes, 1, 0, 0, 0, 0, time.UTC)
	total, qtd, err := u.repo.TotalDoPeriodo(ctx, inicio, inicio.AddDate(0, 1, 0))
	if err != nil {
		return Resumo{}, err
	}
	resumo := Resumo{Total: arredondar(total), Quantidade: qtd}
	if qtd > 0 {
		resumo.TicketMedio = arredondar(total / float64(qtd))
	}
	return resumo, nil
}

type ListarVendas struct {
	repo ports.VendaRepository
}

func NewListarVendas(repo ports.VendaRepository) *ListarVendas {
	return &ListarVendas{repo: repo}
}

func (u *ListarVendas) Executar(ctx context.Context) ([]entity.Venda, error) {
	return u.repo.Listar(ctx)
}

func arredondar(v float64) float64 { return math.Round(v*100) / 100 }
