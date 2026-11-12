package entity

import (
	"net/mail"
	"strings"
	"time"
	"unicode/utf8"
)

const (
	PapelAdmin    = "admin"
	PapelVendedor = "vendedor"

	NomeMaxCaracteres  = 100
	EmailMaxCaracteres = 254
	SenhaMinCaracteres = 8
	// SenhaMaxBytes é o limite do bcrypt: acima disso ele recusa a senha.
	SenhaMaxBytes = 72
)

// Usuario é quem acessa o CRM. Guarda só o hash da senha, nunca a senha em si.
type Usuario struct {
	ID        int64
	Nome      string
	Email     string
	SenhaHash string
	Papel     string
	CriadoEm  time.Time
}

// ErroCampo diz qual campo falhou e por quê, em linguagem que o usuário entende.
type ErroCampo struct {
	Campo    string
	Mensagem string
}

func (e ErroCampo) Error() string { return e.Mensagem }

// ErrosValidacao reúne todos os campos inválidos de uma vez.
type ErrosValidacao []ErroCampo

func (e ErrosValidacao) Error() string {
	mensagens := make([]string, len(e))
	for i, erro := range e {
		mensagens[i] = erro.Mensagem
	}
	return strings.Join(mensagens, " ")
}

// NormalizarEmail tira espaços das pontas e põe em minúsculas, para que
// " Ana@X.com " e "ana@x.com" sejam tratados como o mesmo e-mail.
func NormalizarEmail(email string) string {
	return strings.ToLower(strings.TrimSpace(email))
}

func ValidarNome(nome string) error {
	nome = strings.TrimSpace(nome)
	if nome == "" {
		return ErroCampo{"nome", "O nome é obrigatório."}
	}
	if utf8.RuneCountInString(nome) > NomeMaxCaracteres {
		return ErroCampo{"nome", "O nome deve ter no máximo 100 caracteres."}
	}
	return nil
}

func ValidarEmail(email string) error {
	email = NormalizarEmail(email)
	if email == "" {
		return ErroCampo{"email", "O e-mail é obrigatório."}
	}
	if len(email) > EmailMaxCaracteres {
		return ErroCampo{"email", "O e-mail deve ter no máximo 254 caracteres."}
	}
	invalido := ErroCampo{"email", "Informe um e-mail válido, como nome@empresa.com.br."}

	// ParseAddress aceita formas como "Ana <ana@x.com>"; exigir que o endereço seja idêntico ao texto recusa isso.
	endereco, err := mail.ParseAddress(email)
	if err != nil || endereco.Address != email {
		return invalido
	}
	dominio := email[strings.LastIndex(email, "@")+1:]
	if !strings.Contains(dominio, ".") || strings.HasPrefix(dominio, ".") || strings.HasSuffix(dominio, ".") {
		return invalido
	}
	return nil
}

func ValidarSenha(senha string) error {
	if utf8.RuneCountInString(senha) < SenhaMinCaracteres {
		return ErroCampo{"senha", "A senha deve ter no mínimo 8 caracteres."}
	}
	if len(senha) > SenhaMaxBytes {
		return ErroCampo{"senha", "A senha deve ter no máximo 72 bytes (limite do bcrypt)."}
	}
	return nil
}

func ValidarPapel(papel string) error {
	if papel != PapelAdmin && papel != PapelVendedor {
		return ErroCampo{"papel", "O papel deve ser 'admin' ou 'vendedor'."}
	}
	return nil
}

// ValidarNovoUsuario valida o cadastro: todos os campos são obrigatórios.
func ValidarNovoUsuario(nome, email, senha, papel string) error {
	return reunir(ValidarNome(nome), ValidarEmail(email), ValidarSenha(senha), ValidarPapel(papel))
}

// ValidarAtualizacao valida a edição: a senha é opcional (vazia = manter a atual).
func ValidarAtualizacao(nome, email, senha, papel string) error {
	var erroSenha error
	if senha != "" {
		erroSenha = ValidarSenha(senha)
	}
	return reunir(ValidarNome(nome), ValidarEmail(email), erroSenha, ValidarPapel(papel))
}

// reunir junta os erros não nulos. Devolve nil de verdade quando não há nenhum.
func reunir(erros ...error) error {
	var todos ErrosValidacao
	for _, err := range erros {
		if campo, ok := err.(ErroCampo); ok {
			todos = append(todos, campo)
		}
	}
	if len(todos) == 0 {
		return nil
	}
	return todos
}
