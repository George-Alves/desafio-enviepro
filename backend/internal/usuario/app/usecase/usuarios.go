package usecase

import (
	"context"
	"strings"

	"github.com/George-Alves/desafio-enviepro/backend/internal/usuario/domain/entity"
	"github.com/George-Alves/desafio-enviepro/backend/internal/usuario/domain/ports"
	"golang.org/x/crypto/bcrypt"
)

// DadosNovoUsuario é o que se precisa para cadastrar um usuário. A senha chega em texto
// e nunca sai desta camada nesse formato: o que vai para o banco é só o hash.
type DadosNovoUsuario struct {
	Nome  string
	Email string
	Senha string
	Papel string
}

// DadosAtualizacao é o que se pode alterar. Senha vazia significa "manter a atual".
type DadosAtualizacao struct {
	Nome  string
	Email string
	Senha string
	Papel string
}

// gerarHash é o único lugar do projeto que chama o bcrypt (custo padrão).
// O salt aleatório é sorteado a cada chamada e fica dentro do próprio hash.
func gerarHash(senha string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(senha), bcrypt.DefaultCost)
	return string(hash), err
}

type CriarUsuario struct {
	repo ports.UsuarioRepository
}

func NewCriarUsuario(repo ports.UsuarioRepository) *CriarUsuario {
	return &CriarUsuario{repo: repo}
}

// Executar limpa e valida os dados, gera o hash da senha e grava. Se a validação falhar,
// nada é gravado. Devolve ports.ErrEmailJaCadastrado se o e-mail já existir.
func (u *CriarUsuario) Executar(ctx context.Context, in DadosNovoUsuario) (entity.Usuario, error) {
	nome := strings.TrimSpace(in.Nome)
	email := entity.NormalizarEmail(in.Email)

	if err := entity.ValidarNovoUsuario(nome, email, in.Senha, in.Papel); err != nil {
		return entity.Usuario{}, err
	}
	hash, err := gerarHash(in.Senha)
	if err != nil {
		return entity.Usuario{}, err
	}
	return u.repo.Criar(ctx, entity.Usuario{Nome: nome, Email: email, SenhaHash: hash, Papel: in.Papel})
}

type ListarUsuarios struct {
	repo ports.UsuarioRepository
}

func NewListarUsuarios(repo ports.UsuarioRepository) *ListarUsuarios {
	return &ListarUsuarios{repo: repo}
}

func (u *ListarUsuarios) Executar(ctx context.Context) ([]entity.Usuario, error) {
	return u.repo.Listar(ctx)
}

type BuscarUsuario struct {
	repo ports.UsuarioRepository
}

func NewBuscarUsuario(repo ports.UsuarioRepository) *BuscarUsuario {
	return &BuscarUsuario{repo: repo}
}

// Executar devolve ports.ErrUsuarioNaoEncontrado se o id não existir.
func (u *BuscarUsuario) Executar(ctx context.Context, id int64) (entity.Usuario, error) {
	return u.repo.BuscarPorID(ctx, id)
}

type AtualizarUsuario struct {
	repo ports.UsuarioRepository
}

func NewAtualizarUsuario(repo ports.UsuarioRepository) *AtualizarUsuario {
	return &AtualizarUsuario{repo: repo}
}

// Executar valida, troca nome, e-mail e papel e, só se vier uma senha nova, gera um hash novo.
// Sem senha, o hash atual fica exatamente como está.
func (u *AtualizarUsuario) Executar(ctx context.Context, id int64, in DadosAtualizacao) (entity.Usuario, error) {
	nome := strings.TrimSpace(in.Nome)
	email := entity.NormalizarEmail(in.Email)

	if err := entity.ValidarAtualizacao(nome, email, in.Senha, in.Papel); err != nil {
		return entity.Usuario{}, err
	}
	atual, err := u.repo.BuscarPorID(ctx, id)
	if err != nil {
		return entity.Usuario{}, err
	}
	atual.Nome, atual.Email, atual.Papel = nome, email, in.Papel
	if in.Senha != "" {
		if atual.SenhaHash, err = gerarHash(in.Senha); err != nil {
			return entity.Usuario{}, err
		}
	}
	if err := u.repo.Atualizar(ctx, atual); err != nil {
		return entity.Usuario{}, err
	}
	return atual, nil
}

type ApagarUsuario struct {
	repo ports.UsuarioRepository
}

func NewApagarUsuario(repo ports.UsuarioRepository) *ApagarUsuario {
	return &ApagarUsuario{repo: repo}
}

// Executar devolve ports.ErrUsuarioNaoEncontrado se o id não existir.
func (u *ApagarUsuario) Executar(ctx context.Context, id int64) error {
	return u.repo.Apagar(ctx, id)
}
