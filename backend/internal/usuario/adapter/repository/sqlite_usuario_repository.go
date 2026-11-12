package repository

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/George-Alves/desafio-enviepro/backend/internal/usuario/domain/entity"
	"github.com/George-Alves/desafio-enviepro/backend/internal/usuario/domain/ports"
	"modernc.org/sqlite"
	sqlite3 "modernc.org/sqlite/lib"
)

// Trava de compilação: se o repositório deixar de cumprir o contrato, o projeto não compila.
var _ ports.UsuarioRepository = (*SQLiteUsuarioRepository)(nil)

const colunas = "id, nome, email, senha_hash, papel, criado_em"

type SQLiteUsuarioRepository struct {
	db *sql.DB
}

func NewSQLiteUsuarioRepository(db *sql.DB) *SQLiteUsuarioRepository {
	return &SQLiteUsuarioRepository{db: db}
}

func (r *SQLiteUsuarioRepository) Criar(ctx context.Context, u entity.Usuario) (entity.Usuario, error) {
	// RETURNING pede ao SQLite o id e a data que ele mesmo gerou, sem uma segunda consulta.
	var criadoEm string
	err := r.db.QueryRowContext(ctx, `
		INSERT INTO usuarios (nome, email, senha_hash, papel)
		VALUES (?, ?, ?, ?)
		RETURNING id, criado_em`,
		u.Nome, u.Email, u.SenhaHash, u.Papel).Scan(&u.ID, &criadoEm)
	if err != nil {
		if eEmailRepetido(err) {
			return entity.Usuario{}, ports.ErrEmailJaCadastrado
		}
		return entity.Usuario{}, err
	}
	if u.CriadoEm, err = time.Parse(time.RFC3339, criadoEm); err != nil {
		return entity.Usuario{}, err
	}
	return u, nil
}

func (r *SQLiteUsuarioRepository) Listar(ctx context.Context) ([]entity.Usuario, error) {
	rows, err := r.db.QueryContext(ctx, "SELECT "+colunas+" FROM usuarios ORDER BY id")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	usuarios := []entity.Usuario{}
	for rows.Next() {
		u, err := escanear(rows)
		if err != nil {
			return nil, err
		}
		usuarios = append(usuarios, u)
	}
	return usuarios, rows.Err()
}

func (r *SQLiteUsuarioRepository) BuscarPorID(ctx context.Context, id int64) (entity.Usuario, error) {
	u, err := escanear(r.db.QueryRowContext(ctx, "SELECT "+colunas+" FROM usuarios WHERE id = ?", id))
	if errors.Is(err, sql.ErrNoRows) {
		return entity.Usuario{}, ports.ErrUsuarioNaoEncontrado
	}
	return u, err
}

func (r *SQLiteUsuarioRepository) Atualizar(ctx context.Context, u entity.Usuario) error {
	res, err := r.db.ExecContext(ctx, `
		UPDATE usuarios SET nome = ?, email = ?, senha_hash = ?, papel = ?
		WHERE id = ?`,
		u.Nome, u.Email, u.SenhaHash, u.Papel, u.ID)
	if err != nil {
		if eEmailRepetido(err) {
			return ports.ErrEmailJaCadastrado
		}
		return err
	}
	return exigirUmaLinha(res)
}

func (r *SQLiteUsuarioRepository) Apagar(ctx context.Context, id int64) error {
	res, err := r.db.ExecContext(ctx, "DELETE FROM usuarios WHERE id = ?", id)
	if err != nil {
		return err
	}
	return exigirUmaLinha(res)
}

// exigirUmaLinha devolve ErrUsuarioNaoEncontrado quando o comando não alterou nenhuma linha.
func exigirUmaLinha(res sql.Result) error {
	n, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if n == 0 {
		return ports.ErrUsuarioNaoEncontrado
	}
	return nil
}

// eEmailRepetido reconhece, pelo código numérico do SQLite (e não pelo texto), a violação do UNIQUE.
func eEmailRepetido(err error) bool {
	var erroSQLite *sqlite.Error
	return errors.As(err, &erroSQLite) && erroSQLite.Code() == sqlite3.SQLITE_CONSTRAINT_UNIQUE
}

// escanear lê uma linha (de QueryRow ou de Rows) para dentro de um Usuario.
func escanear(linha interface{ Scan(...any) error }) (entity.Usuario, error) {
	var u entity.Usuario
	var criadoEm string
	if err := linha.Scan(&u.ID, &u.Nome, &u.Email, &u.SenhaHash, &u.Papel, &criadoEm); err != nil {
		return entity.Usuario{}, err
	}
	var err error
	u.CriadoEm, err = time.Parse(time.RFC3339, criadoEm)
	return u, err
}
