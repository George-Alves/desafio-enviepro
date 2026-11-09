package handler

import (
	"log"
	"net/http"
	"time"

	"github.com/George-Alves/desafio-enviepro/backend/internal/cliente/app/usecase"
	"github.com/gin-gonic/gin"
)

// clienteJSON é o formato público do cliente; a API nunca devolve a entity direto.
type clienteJSON struct {
	ID             int64  `json:"id"`
	Nome           string `json:"nome"`
	Empresa        string `json:"empresa"`
	Email          string `json:"email"`
	Telefone       string `json:"telefone"`
	Status         string `json:"status"`
	Responsavel    string `json:"responsavel"`
	UltimoContato  string `json:"ultimo_contato"`
	DiasSemContato int    `json:"dias_sem_contato"`
	Alerta         bool   `json:"alerta"`
}

type ClienteHandler struct {
	listar *usecase.ListarClientes
}

func NewClienteHandler(listar *usecase.ListarClientes) *ClienteHandler {
	return &ClienteHandler{listar: listar}
}

func (h *ClienteHandler) Registrar(r gin.IRouter) {
	r.GET("/clientes", h.Listar)
}

func (h *ClienteHandler) Listar(c *gin.Context) {
	clientes, err := h.listar.Executar(c.Request.Context(), time.Now())
	if err != nil {
		log.Printf("erro ao listar clientes: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"erro": "não foi possível listar os clientes"})
		return
	}
	resposta := make([]clienteJSON, 0, len(clientes))
	for _, cl := range clientes {
		resposta = append(resposta, clienteJSON{
			ID:             cl.ID,
			Nome:           cl.Nome,
			Empresa:        cl.Empresa,
			Email:          cl.Email,
			Telefone:       cl.Telefone,
			Status:         cl.Status,
			Responsavel:    cl.Responsavel,
			UltimoContato:  cl.UltimoContato.Format("2006-01-02"),
			DiasSemContato: cl.DiasSemContato,
			Alerta:         cl.Alerta,
		})
	}
	c.JSON(http.StatusOK, resposta)
}
