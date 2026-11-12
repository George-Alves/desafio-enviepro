package usecase

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/George-Alves/desafio-enviepro/backend/internal/usuario/domain/entity"
	"github.com/George-Alves/desafio-enviepro/backend/internal/usuario/domain/ports"
	"golang.org/x/crypto/bcrypt"
)

// repoFalso é um repositório em memória que respeita as mesmas regras do de verdade
// (e-mail único sem diferenciar maiúsculas, erros combinados na porta).
type repoFalso struct {
	usuarios      map[int64]entity.Usuario
	proximoID     int64
	chamadasCriar int
}

func novoRepoFalso() *repoFalso { return &repoFalso{usuarios: map[int64]entity.Usuario{}, proximoID: 1} }

func (r *repoFalso) emailEmUso(email string, ignorarID int64) bool {
	for _, u := range r.usuarios {
		if u.ID != ignorarID && strings.EqualFold(u.Email, email) {
			return true
		}
	}
	return false
}

func (r *repoFalso) Criar(_ context.Context, u entity.Usuario) (entity.Usuario, error) {
	r.chamadasCriar++
	if r.emailEmUso(u.Email, 0) {
		return entity.Usuario{}, ports.ErrEmailJaCadastrado
	}
	u.ID, u.CriadoEm = r.proximoID, time.Now().UTC()
	r.proximoID++
	r.usuarios[u.ID] = u
	return u, nil
}

func (r *repoFalso) Listar(context.Context) ([]entity.Usuario, error) {
	lista := []entity.Usuario{}
	for id := int64(1); id < r.proximoID; id++ {
		if u, ok := r.usuarios[id]; ok {
			lista = append(lista, u)
		}
	}
	return lista, nil
}

func (r *repoFalso) BuscarPorID(_ context.Context, id int64) (entity.Usuario, error) {
	u, ok := r.usuarios[id]
	if !ok {
		return entity.Usuario{}, ports.ErrUsuarioNaoEncontrado
	}
	return u, nil
}

func (r *repoFalso) Atualizar(_ context.Context, u entity.Usuario) error {
	if _, ok := r.usuarios[u.ID]; !ok {
		return ports.ErrUsuarioNaoEncontrado
	}
	if r.emailEmUso(u.Email, u.ID) {
		return ports.ErrEmailJaCadastrado
	}
	r.usuarios[u.ID] = u
	return nil
}

func (r *repoFalso) Apagar(_ context.Context, id int64) error {
	if _, ok := r.usuarios[id]; !ok {
		return ports.ErrUsuarioNaoEncontrado
	}
	delete(r.usuarios, id)
	return nil
}

func dadosValidos(email string) DadosNovoUsuario {
	return DadosNovoUsuario{Nome: "Ana", Email: email, Senha: "senha-forte-1", Papel: entity.PapelVendedor}
}

func TestCriarUsuario_GravaSoOHashBcrypt(t *testing.T) {
	repo := novoRepoFalso()
	criar := NewCriarUsuario(repo)

	criado, err := criar.Executar(context.Background(), DadosNovoUsuario{
		Nome: "  Ana  ", Email: "  Ana@Teste.COM ", Senha: "senha-forte-1", Papel: entity.PapelAdmin,
	})
	if err != nil {
		t.Fatal(err)
	}

	guardado := repo.usuarios[criado.ID]
	if guardado.Nome != "Ana" || guardado.Email != "ana@teste.com" {
		t.Errorf("nome e e-mail deveriam ser limpos: %q, %q", guardado.Nome, guardado.Email)
	}
	if !strings.HasPrefix(guardado.SenhaHash, "$2a$") {
		t.Errorf("o hash deveria começar com $2a$, veio %q", guardado.SenhaHash)
	}
	if strings.Contains(guardado.SenhaHash, "senha-forte-1") || guardado.SenhaHash == "senha-forte-1" {
		t.Error("a senha em texto não pode aparecer no que foi gravado")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(guardado.SenhaHash), []byte("senha-forte-1")); err != nil {
		t.Errorf("o hash deveria conferir com a senha original: %v", err)
	}
	if err := bcrypt.CompareHashAndPassword([]byte(guardado.SenhaHash), []byte("outra-senha-qualquer")); err == nil {
		t.Error("uma senha errada não deveria conferir")
	}
}

func TestCriarUsuario_MesmaSenhaGeraHashesDiferentes_PorCausaDoSalt(t *testing.T) {
	repo := novoRepoFalso()
	criar := NewCriarUsuario(repo)

	a, _ := criar.Executar(context.Background(), dadosValidos("a@teste.com"))
	b, _ := criar.Executar(context.Background(), dadosValidos("b@teste.com"))

	if a.SenhaHash == b.SenhaHash {
		t.Error("duas contas com a mesma senha deveriam ter hashes diferentes (salt aleatório)")
	}
}

func TestCriarUsuario_DadosInvalidosNaoChegamAoRepositorio(t *testing.T) {
	repo := novoRepoFalso()
	criar := NewCriarUsuario(repo)

	// o critério da issue: senha de 5 caracteres é recusada com mensagem clara
	in := dadosValidos("ana@teste.com")
	in.Senha = "12345"
	_, err := criar.Executar(context.Background(), in)

	var erros entity.ErrosValidacao
	if !errors.As(err, &erros) || !strings.Contains(err.Error(), "mínimo 8") {
		t.Fatalf("esperava erro de validação falando do mínimo de 8, veio %v", err)
	}
	if repo.chamadasCriar != 0 {
		t.Error("com dado inválido o repositório não deveria ser chamado")
	}

	in = dadosValidos("isso-nao-e-email")
	in.Papel = "chefe"
	if _, err := criar.Executar(context.Background(), in); !errors.As(err, &erros) || len(erros) != 2 {
		t.Errorf("esperava 2 erros (email e papel), veio %v", err)
	}
}

func TestCriarUsuario_EmailRepetidoDevolveErroDaPorta(t *testing.T) {
	repo := novoRepoFalso()
	criar := NewCriarUsuario(repo)
	if _, err := criar.Executar(context.Background(), dadosValidos("ana@teste.com")); err != nil {
		t.Fatal(err)
	}
	_, err := criar.Executar(context.Background(), dadosValidos("ANA@teste.com"))
	if !errors.Is(err, ports.ErrEmailJaCadastrado) {
		t.Errorf("esperava ErrEmailJaCadastrado, veio %v", err)
	}
}

func TestAtualizarUsuario_SenhaOpcional(t *testing.T) {
	repo := novoRepoFalso()
	ctx := context.Background()
	criado, _ := NewCriarUsuario(repo).Executar(ctx, dadosValidos("ana@teste.com"))
	atualizar := NewAtualizarUsuario(repo)

	// sem senha: troca nome e papel, o hash continua exatamente o mesmo
	feito, err := atualizar.Executar(ctx, criado.ID, DadosAtualizacao{
		Nome: "Ana Souza", Email: "ana@teste.com", Papel: entity.PapelAdmin,
	})
	if err != nil {
		t.Fatal(err)
	}
	if feito.Nome != "Ana Souza" || feito.Papel != entity.PapelAdmin || feito.SenhaHash != criado.SenhaHash {
		t.Errorf("sem senha nova, o hash deveria ficar igual: %+v", feito)
	}

	// com senha: hash novo, que confere com a nova senha e não com a antiga
	feito, err = atualizar.Executar(ctx, criado.ID, DadosAtualizacao{
		Nome: "Ana Souza", Email: "ana@teste.com", Senha: "nova-senha-9", Papel: entity.PapelAdmin,
	})
	if err != nil {
		t.Fatal(err)
	}
	if feito.SenhaHash == criado.SenhaHash {
		t.Error("com senha nova, o hash deveria mudar")
	}
	if bcrypt.CompareHashAndPassword([]byte(feito.SenhaHash), []byte("nova-senha-9")) != nil {
		t.Error("o hash novo deveria conferir com a senha nova")
	}
	if bcrypt.CompareHashAndPassword([]byte(feito.SenhaHash), []byte("senha-forte-1")) == nil {
		t.Error("o hash novo não deveria conferir com a senha antiga")
	}
}

func TestAtualizarUsuario_Erros(t *testing.T) {
	repo := novoRepoFalso()
	ctx := context.Background()
	criar := NewCriarUsuario(repo)
	ana, _ := criar.Executar(ctx, dadosValidos("ana@teste.com"))
	_, _ = criar.Executar(ctx, dadosValidos("bia@teste.com"))
	atualizar := NewAtualizarUsuario(repo)
	base := DadosAtualizacao{Nome: "Ana", Email: "ana@teste.com", Papel: entity.PapelVendedor}

	curta := base
	curta.Senha = "123"
	if _, err := atualizar.Executar(ctx, ana.ID, curta); !strings.Contains(errString(err), "mínimo 8") {
		t.Errorf("senha curta deveria falar do mínimo de 8, veio %v", err)
	}
	if repo.usuarios[ana.ID].SenhaHash != ana.SenhaHash {
		t.Error("uma atualização inválida não pode alterar o hash")
	}

	emailDeOutro := base
	emailDeOutro.Email = "BIA@teste.com"
	if _, err := atualizar.Executar(ctx, ana.ID, emailDeOutro); !errors.Is(err, ports.ErrEmailJaCadastrado) {
		t.Errorf("esperava ErrEmailJaCadastrado, veio %v", err)
	}

	if _, err := atualizar.Executar(ctx, 9999, base); !errors.Is(err, ports.ErrUsuarioNaoEncontrado) {
		t.Errorf("esperava ErrUsuarioNaoEncontrado, veio %v", err)
	}
}

func errString(err error) string {
	if err == nil {
		return ""
	}
	return err.Error()
}

func TestListarBuscarEApagarUsuario(t *testing.T) {
	repo := novoRepoFalso()
	ctx := context.Background()
	criar := NewCriarUsuario(repo)
	ana, _ := criar.Executar(ctx, dadosValidos("ana@teste.com"))
	_, _ = criar.Executar(ctx, dadosValidos("bia@teste.com"))

	todos, err := NewListarUsuarios(repo).Executar(ctx)
	if err != nil || len(todos) != 2 {
		t.Fatalf("esperava 2 usuários, veio %d (erro %v)", len(todos), err)
	}
	if u, err := NewBuscarUsuario(repo).Executar(ctx, ana.ID); err != nil || u.Email != "ana@teste.com" {
		t.Errorf("busca falhou: %+v, %v", u, err)
	}
	if _, err := NewBuscarUsuario(repo).Executar(ctx, 9999); !errors.Is(err, ports.ErrUsuarioNaoEncontrado) {
		t.Errorf("esperava ErrUsuarioNaoEncontrado, veio %v", err)
	}

	apagar := NewApagarUsuario(repo)
	if err := apagar.Executar(ctx, ana.ID); err != nil {
		t.Fatal(err)
	}
	if err := apagar.Executar(ctx, ana.ID); !errors.Is(err, ports.ErrUsuarioNaoEncontrado) {
		t.Errorf("apagar de novo deveria dar ErrUsuarioNaoEncontrado, veio %v", err)
	}
}
