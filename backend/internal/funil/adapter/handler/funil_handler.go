package handler

import (
	"log"
	"net/http"

	"github.com/George-Alves/desafio-enviepro/backend/internal/funil/app/usecase"
	"github.com/gin-gonic/gin"
)

type etapaJSON struct {
	Etapa      string  `json:"etapa"`
	Quantidade int     `json:"quantidade"`
	Valor      float64 `json:"valor"`
}

type FunilHandler struct {
	obter *usecase.ObterFunil
}

func NewFunilHandler(obter *usecase.ObterFunil) *FunilHandler {
	return &FunilHandler{obter: obter}
}

func (h *FunilHandler) Registrar(r gin.IRouter) {
	r.GET("/funil", h.Obter)
}

func (h *FunilHandler) Obter(c *gin.Context) {
	etapas, err := h.obter.Executar(c.Request.Context())
	if err != nil {
		log.Printf("erro ao obter funil: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"erro": "não foi possível obter o funil"})
		return
	}
	resposta := make([]etapaJSON, 0, len(etapas))
	for _, e := range etapas {
		resposta = append(resposta, etapaJSON{Etapa: e.Etapa, Quantidade: e.Quantidade, Valor: e.Valor})
	}
	c.JSON(http.StatusOK, resposta)
}
