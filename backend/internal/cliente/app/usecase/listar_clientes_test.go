package usecase

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/George-Alves/desafio-enviepro/backend/internal/cliente/domain/entity"
)

type repoFalso struct {
	clientes []entity.Cliente
	err      error
}

func (r repoFalso) Listar(context.Context) ([]entity.Cliente, error) { return r.clientes, r.err }

func TestListarClientes_AlertaSoAcimaDe15Dias(t *testing.T) {
	agora := time.Date(2026, 11, 9, 18, 30, 0, 0, time.UTC)
	repo := repoFalso{clientes: []entity.Cliente{
		{ID: 1, Nome: "hoje", UltimoContato: agora},
		{ID: 2, Nome: "quinze", UltimoContato: agora.AddDate(0, 0, -15)},
		{ID: 3, Nome: "dezesseis", UltimoContato: agora.AddDate(0, 0, -16)},
	}}

	got, err := NewListarClientes(repo).Executar(context.Background(), agora)
	if err != nil {
		t.Fatal(err)
	}
	esperado := []struct {
		dias   int
		alerta bool
	}{{0, false}, {15, false}, {16, true}}
	for i, e := range esperado {
		if got[i].DiasSemContato != e.dias || got[i].Alerta != e.alerta {
			t.Errorf("%s: dias=%d alerta=%v, esperado dias=%d alerta=%v",
				got[i].Nome, got[i].DiasSemContato, got[i].Alerta, e.dias, e.alerta)
		}
	}
}

func TestListarClientes_PropagaErroDoRepositorio(t *testing.T) {
	erro := errors.New("banco fora do ar")
	_, err := NewListarClientes(repoFalso{err: erro}).Executar(context.Background(), time.Now())
	if !errors.Is(err, erro) {
		t.Fatalf("esperava o erro do repositório, veio %v", err)
	}
}
