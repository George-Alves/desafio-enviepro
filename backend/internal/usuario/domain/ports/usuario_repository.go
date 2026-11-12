package ports

import (
	"context"
	"errors"

	"github.com/George-Alves/desafio-enviepro/backend/internal/usuario/domain/entity"
)

// Erros que qualquer implementação do repositório deve devolver nessas situações,
// para o resto do sistema não depender de mensagens do banco de dados.
var (
	ErrUsuarioNaoEncontrado = errors.New("usuário não encontrado")
	ErrEmailJaCadastrado    = errors.New("e-mail já cadastrado")
)

// UsuarioRepository é o que a regra de negócio precisa para guardar e ler usuários.
type UsuarioRepository interface {
	// Criar grava um usuário novo e o devolve como ficou (com ID e CriadoEm preenchidos).
	// Devolve ErrEmailJaCadastrado se o e-mail já existir.
	Criar(ctx context.Context, usuario entity.Usuario) (entity.Usuario, error)

	Listar(ctx context.Context) ([]entity.Usuario, error)

	// BuscarPorID devolve ErrUsuarioNaoEncontrado se o id não existir.
	BuscarPorID(ctx context.Context, id int64) (entity.Usuario, error)

	// Atualizar regrava nome, e-mail, hash da senha e papel.
	// Devolve ErrUsuarioNaoEncontrado ou ErrEmailJaCadastrado (e-mail de outro usuário).
	Atualizar(ctx context.Context, usuario entity.Usuario) error

	// Apagar devolve ErrUsuarioNaoEncontrado se o id não existir.
	Apagar(ctx context.Context, id int64) error
}
