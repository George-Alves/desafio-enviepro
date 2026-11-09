package handler

import (
	"log"
	"net/http"
	"time"

	"github.com/George-Alves/desafio-enviepro/backend/internal/venda/app/usecase"
	"github.com/gin-gonic/gin"
)

type resumoJSON struct {
	Total       float64 `json:"total"`
	Quantidade  int     `json:"quantidade"`
	TicketMedio float64 `json:"ticket_medio"`
}

type vendaJSON struct {
	ID        int64   `json:"id"`
	ClienteID int64   `json:"cliente_id"`
	Valor     float64 `json:"valor"`
	Data      string  `json:"data"`
	Produto   string  `json:"produto"`
}

type VendaHandler struct {
	resumo *usecase.ResumoVendasDoMes
	listar *usecase.ListarVendas
}

func NewVendaHandler(resumo *usecase.ResumoVendasDoMes, listar *usecase.ListarVendas) *VendaHandler {
	return &VendaHandler{resumo: resumo, listar: listar}
}

func (h *VendaHandler) Registrar(r gin.IRouter) {
	r.GET("/vendas/resumo", h.Resumo)
	r.GET("/vendas", h.Listar)
}

func (h *VendaHandler) Resumo(c *gin.Context) {
	r, err := h.resumo.Executar(c.Request.Context(), time.Now())
	if err != nil {
		log.Printf("erro ao resumir vendas: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"erro": "não foi possível resumir as vendas"})
		return
	}
	c.JSON(http.StatusOK, resumoJSON{Total: r.Total, Quantidade: r.Quantidade, TicketMedio: r.TicketMedio})
}

func (h *VendaHandler) Listar(c *gin.Context) {
	vendas, err := h.listar.Executar(c.Request.Context())
	if err != nil {
		log.Printf("erro ao listar vendas: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"erro": "não foi possível listar as vendas"})
		return
	}
	resposta := make([]vendaJSON, 0, len(vendas))
	for _, v := range vendas {
		resposta = append(resposta, vendaJSON{
			ID: v.ID, ClienteID: v.ClienteID, Valor: v.Valor,
			Data: v.Data.Format("2006-01-02"), Produto: v.Produto,
		})
	}
	c.JSON(http.StatusOK, resposta)
}
