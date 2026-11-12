package handler

import (
	"errors"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/George-Alves/desafio-enviepro/backend/internal/usuario/app/usecase"
	"github.com/George-Alves/desafio-enviepro/backend/internal/usuario/domain/entity"
	"github.com/George-Alves/desafio-enviepro/backend/internal/usuario/domain/ports"
	"github.com/gin-gonic/gin"
)

// limiteDoCorpo evita que alguém mande um JSON gigantesco.
const limiteDoCorpo = 1 << 20 // 1 MB

// usuarioJSON é o formato público do usuário. Não tem campo de senha nem de hash:
// como o handler só devolve este formato, o hash não tem por onde sair.
type usuarioJSON struct {
	ID       int64  `json:"id"`
	Nome     string `json:"nome"`
	Email    string `json:"email"`
	Papel    string `json:"papel"`
	CriadoEm string `json:"criado_em"`
}

// Formatos de entrada. Campos extras no JSON (como "id" ou "senha_hash") são ignorados,
// porque estas structs não os têm.
type novoUsuarioJSON struct {
	Nome  string `json:"nome"`
	Email string `json:"email"`
	Senha string `json:"senha"`
	Papel string `json:"papel"`
}

type atualizarUsuarioJSON struct {
	Nome  string `json:"nome"`
	Email string `json:"email"`
	Senha string `json:"senha"` // opcional: vazia = manter a senha atual
	Papel string `json:"papel"`
}

type UsuarioHandler struct {
	criar     *usecase.CriarUsuario
	listar    *usecase.ListarUsuarios
	buscar    *usecase.BuscarUsuario
	atualizar *usecase.AtualizarUsuario
	apagar    *usecase.ApagarUsuario
}

func NewUsuarioHandler(
	criar *usecase.CriarUsuario,
	listar *usecase.ListarUsuarios,
	buscar *usecase.BuscarUsuario,
	atualizar *usecase.AtualizarUsuario,
	apagar *usecase.ApagarUsuario,
) *UsuarioHandler {
	return &UsuarioHandler{criar: criar, listar: listar, buscar: buscar, atualizar: atualizar, apagar: apagar}
}

func (h *UsuarioHandler) Registrar(r gin.IRouter) {
	r.POST("/usuarios", h.Criar)
	r.GET("/usuarios", h.Listar)
	r.GET("/usuarios/:id", h.Buscar)
	r.PUT("/usuarios/:id", h.Atualizar)
	r.DELETE("/usuarios/:id", h.Apagar)
}

func (h *UsuarioHandler) Criar(c *gin.Context) {
	var corpo novoUsuarioJSON
	if !lerCorpo(c, &corpo) {
		return
	}
	u, err := h.criar.Executar(c.Request.Context(), usecase.DadosNovoUsuario{
		Nome: corpo.Nome, Email: corpo.Email, Senha: corpo.Senha, Papel: corpo.Papel,
	})
	if err != nil {
		responderErro(c, err)
		return
	}
	c.JSON(http.StatusCreated, paraJSON(u))
}

func (h *UsuarioHandler) Listar(c *gin.Context) {
	usuarios, err := h.listar.Executar(c.Request.Context())
	if err != nil {
		responderErro(c, err)
		return
	}
	resposta := make([]usuarioJSON, 0, len(usuarios))
	for _, u := range usuarios {
		resposta = append(resposta, paraJSON(u))
	}
	c.JSON(http.StatusOK, resposta)
}

func (h *UsuarioHandler) Buscar(c *gin.Context) {
	id, ok := lerID(c)
	if !ok {
		return
	}
	u, err := h.buscar.Executar(c.Request.Context(), id)
	if err != nil {
		responderErro(c, err)
		return
	}
	c.JSON(http.StatusOK, paraJSON(u))
}

func (h *UsuarioHandler) Atualizar(c *gin.Context) {
	id, ok := lerID(c)
	if !ok {
		return
	}
	var corpo atualizarUsuarioJSON
	if !lerCorpo(c, &corpo) {
		return
	}
	u, err := h.atualizar.Executar(c.Request.Context(), id, usecase.DadosAtualizacao{
		Nome: corpo.Nome, Email: corpo.Email, Senha: corpo.Senha, Papel: corpo.Papel,
	})
	if err != nil {
		responderErro(c, err)
		return
	}
	c.JSON(http.StatusOK, paraJSON(u))
}

func (h *UsuarioHandler) Apagar(c *gin.Context) {
	id, ok := lerID(c)
	if !ok {
		return
	}
	if err := h.apagar.Executar(c.Request.Context(), id); err != nil {
		responderErro(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}

func paraJSON(u entity.Usuario) usuarioJSON {
	return usuarioJSON{
		ID: u.ID, Nome: u.Nome, Email: u.Email, Papel: u.Papel,
		CriadoEm: u.CriadoEm.UTC().Format(time.RFC3339),
	}
}

// lerCorpo lê o JSON do pedido (com limite de tamanho). Se for inválido, já responde 400 e devolve false.
func lerCorpo(c *gin.Context, destino any) bool {
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, limiteDoCorpo)
	if err := c.ShouldBindJSON(destino); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"erro": "O corpo da requisição precisa ser um JSON válido."})
		return false
	}
	return true
}

// lerID converte o :id da URL. Se não for um número positivo, já responde 400 e devolve false.
func lerID(c *gin.Context) (int64, bool) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"erro": "O id deve ser um número positivo."})
		return 0, false
	}
	return id, true
}

// responderErro é o único lugar que decide o código HTTP de cada erro.
// Erros inesperados vão para o log do servidor; quem chamou recebe uma mensagem genérica.
func responderErro(c *gin.Context, err error) {
	var validacao entity.ErrosValidacao
	switch {
	case errors.As(err, &validacao):
		c.JSON(http.StatusBadRequest, gin.H{"erro": validacao.Error()})
	case errors.Is(err, ports.ErrUsuarioNaoEncontrado):
		c.JSON(http.StatusNotFound, gin.H{"erro": "Usuário não encontrado."})
	case errors.Is(err, ports.ErrEmailJaCadastrado):
		c.JSON(http.StatusConflict, gin.H{"erro": "Este e-mail já está cadastrado."})
	default:
		log.Printf("erro inesperado em %s %s: %v", c.Request.Method, c.Request.URL.Path, err)
		c.JSON(http.StatusInternalServerError, gin.H{"erro": "Erro interno do servidor."})
	}
}
