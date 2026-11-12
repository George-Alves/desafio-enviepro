package main

import (
	"context"
	"errors"
	"path/filepath"
	"testing"
	"time"

	usuarioRepo "github.com/George-Alves/desafio-enviepro/backend/internal/usuario/adapter/repository"
	"github.com/George-Alves/desafio-enviepro/backend/internal/usuario/domain/entity"
	"github.com/George-Alves/desafio-enviepro/backend/internal/usuario/domain/ports"
)

// Estes testes usam um banco SQLite de verdade, num arquivo temporário, criado por abrirBanco.
// Ficam no pacote principal porque o esquema das tabelas mora em db.go: assim há uma única fonte.

func novoRepoDeUsuarios(t *testing.T) *usuarioRepo.SQLiteUsuarioRepository {
	t.Helper()
	db, err := abrirBanco(filepath.Join(t.TempDir(), "usuarios.db"))
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { db.Close() })
	return usuarioRepo.NewSQLiteUsuarioRepository(db)
}

func novoUsuario(nome, email string) entity.Usuario {
	return entity.Usuario{Nome: nome, Email: email, SenhaHash: "$2a$10$hash-de-mentira", Papel: entity.PapelVendedor}
}

func TestUsuarioRepo_CriarEBuscarPorID(t *testing.T) {
	repo := novoRepoDeUsuarios(t)
	ctx := context.Background()

	criado, err := repo.Criar(ctx, novoUsuario("Ana", "ana@teste.com"))
	if err != nil {
		t.Fatal(err)
	}
	if criado.ID == 0 {
		t.Error("o ID deveria vir preenchido")
	}
	if time.Since(criado.CriadoEm) > time.Minute || time.Since(criado.CriadoEm) < -time.Minute {
		t.Errorf("CriadoEm deveria ser agora, veio %v", criado.CriadoEm)
	}

	lido, err := repo.BuscarPorID(ctx, criado.ID)
	if err != nil {
		t.Fatal(err)
	}
	if lido.Nome != "Ana" || lido.Email != "ana@teste.com" || lido.Papel != entity.PapelVendedor ||
		lido.SenhaHash != "$2a$10$hash-de-mentira" || !lido.CriadoEm.Equal(criado.CriadoEm) {
		t.Errorf("o que foi lido difere do que foi criado: %+v", lido)
	}
}

func TestUsuarioRepo_EmailRepetidoFalha_MesmoComMaiusculas(t *testing.T) {
	repo := novoRepoDeUsuarios(t)
	ctx := context.Background()

	if _, err := repo.Criar(ctx, novoUsuario("Ana", "ana@teste.com")); err != nil {
		t.Fatal(err)
	}
	_, err := repo.Criar(ctx, novoUsuario("Outra Ana", "ANA@Teste.com"))
	if !errors.Is(err, ports.ErrEmailJaCadastrado) {
		t.Fatalf("esperava ErrEmailJaCadastrado, veio %v", err)
	}

	todos, _ := repo.Listar(ctx)
	if len(todos) != 1 {
		t.Errorf("o segundo usuário não deveria ter sido gravado; há %d", len(todos))
	}
}

func TestUsuarioRepo_ListarEmOrdemDeCadastro(t *testing.T) {
	repo := novoRepoDeUsuarios(t)
	ctx := context.Background()

	if vazio, err := repo.Listar(ctx); err != nil || len(vazio) != 0 {
		t.Fatalf("banco novo deveria listar zero usuários, veio %d (erro %v)", len(vazio), err)
	}
	for _, n := range []string{"Zeca", "Ana", "Bia"} {
		if _, err := repo.Criar(ctx, novoUsuario(n, n+"@teste.com")); err != nil {
			t.Fatal(err)
		}
	}
	todos, err := repo.Listar(ctx)
	if err != nil || len(todos) != 3 {
		t.Fatalf("esperava 3 usuários, veio %d (erro %v)", len(todos), err)
	}
	if todos[0].Nome != "Zeca" || todos[1].Nome != "Ana" || todos[2].Nome != "Bia" {
		t.Errorf("ordem deveria ser a do cadastro: %s, %s, %s", todos[0].Nome, todos[1].Nome, todos[2].Nome)
	}
}

func TestUsuarioRepo_Atualizar(t *testing.T) {
	repo := novoRepoDeUsuarios(t)
	ctx := context.Background()
	ana, _ := repo.Criar(ctx, novoUsuario("Ana", "ana@teste.com"))
	bia, _ := repo.Criar(ctx, novoUsuario("Bia", "bia@teste.com"))

	ana.Nome, ana.Papel, ana.SenhaHash = "Ana Souza", entity.PapelAdmin, "$2a$10$outro"
	if err := repo.Atualizar(ctx, ana); err != nil {
		t.Fatal(err)
	}
	lida, _ := repo.BuscarPorID(ctx, ana.ID)
	if lida.Nome != "Ana Souza" || lida.Papel != entity.PapelAdmin || lida.SenhaHash != "$2a$10$outro" {
		t.Errorf("atualização não foi gravada: %+v", lida)
	}

	// trocar para o e-mail de outro usuário
	ana.Email = "BIA@teste.com"
	if err := repo.Atualizar(ctx, ana); !errors.Is(err, ports.ErrEmailJaCadastrado) {
		t.Errorf("esperava ErrEmailJaCadastrado, veio %v", err)
	}
	if intacta, _ := repo.BuscarPorID(ctx, ana.ID); intacta.Email != "ana@teste.com" {
		t.Errorf("o e-mail da Ana não deveria ter mudado: %s", intacta.Email)
	}

	// id que não existe
	fantasma := bia
	fantasma.ID = 9999
	if err := repo.Atualizar(ctx, fantasma); !errors.Is(err, ports.ErrUsuarioNaoEncontrado) {
		t.Errorf("esperava ErrUsuarioNaoEncontrado, veio %v", err)
	}
}

func TestUsuarioRepo_ApagarEIdInexistente(t *testing.T) {
	repo := novoRepoDeUsuarios(t)
	ctx := context.Background()
	ana, _ := repo.Criar(ctx, novoUsuario("Ana", "ana@teste.com"))

	if err := repo.Apagar(ctx, ana.ID); err != nil {
		t.Fatal(err)
	}
	if _, err := repo.BuscarPorID(ctx, ana.ID); !errors.Is(err, ports.ErrUsuarioNaoEncontrado) {
		t.Errorf("depois de apagar, BuscarPorID deveria dar ErrUsuarioNaoEncontrado, veio %v", err)
	}
	if err := repo.Apagar(ctx, ana.ID); !errors.Is(err, ports.ErrUsuarioNaoEncontrado) {
		t.Errorf("apagar de novo deveria dar ErrUsuarioNaoEncontrado, veio %v", err)
	}
	if _, err := repo.BuscarPorID(ctx, 9999); !errors.Is(err, ports.ErrUsuarioNaoEncontrado) {
		t.Errorf("id inexistente deveria dar ErrUsuarioNaoEncontrado, veio %v", err)
	}

	// o e-mail apagado pode ser usado de novo
	if _, err := repo.Criar(ctx, novoUsuario("Ana Nova", "ana@teste.com")); err != nil {
		t.Errorf("e-mail de usuário apagado deveria poder ser reutilizado: %v", err)
	}
}
