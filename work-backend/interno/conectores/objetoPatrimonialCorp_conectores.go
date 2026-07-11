package conectores

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"
	"work-backend/interno/caso_de_uso"
	"work-backend/interno/entidade"
)

type ObjetoPatrimonialCorpConector struct {
	casoDeUso *caso_de_uso.ObjetoPatrimonialCorpCasoDeUso
}

func NovoObjetoPatrimonialCorpConector(uc *caso_de_uso.ObjetoPatrimonialCorpCasoDeUso) *ObjetoPatrimonialCorpConector {
	return &ObjetoPatrimonialCorpConector{casoDeUso: uc}
}

type InserirObjetoPatrimonialCorpDTO struct {
	Nome             string    `json:"nome"`
	ValorPatrimonial float64   `json:"valor_patrimonial"`
	DataAquisicao    time.Time `json:"data_aquisicao"`
	EmpresaID        int       `json:"empresa_id"`
}

type EditarObjetoPatrimonialCorpDTO struct {
	Nome             string    `json:"nome"`
	ValorPatrimonial float64   `json:"valor_patrimonial"`
	DataAquisicao    time.Time `json:"data_aquisicao"`
}

func (conector *ObjetoPatrimonialCorpConector) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	partes := strings.Split(strings.Trim(r.URL.Path, "/"), "/")

	if len(partes) == 0 || partes[0] != "objetos-patrimoniais-corp" {
		http.NotFound(w, r)
		return
	}

	switch r.Method {
	case http.MethodGet:
		if len(partes) == 2 {
			conector.tratarBuscaPorID(w, r, partes[1])
		} else if len(partes) == 3 && partes[2] == "depreciacao" {
			conector.tratarCalculoDepreciacao(w, r, partes[1])
		} else if len(partes) == 3 && partes[1] == "empresa" {
			conector.tratarListagemPorEmpresa(w, r, partes[2])
		} else {
			http.NotFound(w, r)
		}
	case http.MethodPost:
		if len(partes) == 1 { conector.tratarCadastro(w, r) } else { http.NotFound(w, r) }
	case http.MethodPut:
		if len(partes) == 2 { conector.tratarEdicao(w, r, partes[1]) } else { http.Error(w, "ID ausente", http.StatusBadRequest) }
	case http.MethodDelete:
		if len(partes) == 2 { conector.tratarExclusao(w, r, partes[1]) } else { http.Error(w, "ID ausente", http.StatusBadRequest) }
	default:
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
	}
}

func (conector *ObjetoPatrimonialCorpConector) tratarCadastro(w http.ResponseWriter, r *http.Request) {
	var dto InserirObjetoPatrimonialCorpDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		http.Error(w, "JSON enviado inválido!", http.StatusBadRequest)
		return
	}

	objeto := &entidade.ObjetoPatrimonialCorp{
		Nome:             dto.Nome,
		ValorPatrimonial: dto.ValorPatrimonial,
		DataAquisicao:    dto.DataAquisicao,
		EmpresaID:        dto.EmpresaID,
	}

	if err := conector.casoDeUso.Cadastrar(objeto); err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(objeto)
}

func (conector *ObjetoPatrimonialCorpConector) tratarListagemPorEmpresa(w http.ResponseWriter, r *http.Request, empresaIDStr string) {
	empresaID, err := strconv.Atoi(empresaIDStr)
	if err != nil {
		http.Error(w, "ID da empresa deve ser um inteiro.", http.StatusBadRequest)
		return
	}

	lista, err := conector.casoDeUso.ListarPorEmpresa(empresaID)
	if err != nil {
		http.Error(w, "Erro ao buscar ativos", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(lista)
}

func (conector *ObjetoPatrimonialCorpConector) tratarBuscaPorID(w http.ResponseWriter, r *http.Request, idStr string) {
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID inválido.", http.StatusBadRequest)
		return
	}

	objeto, err := conector.casoDeUso.Buscar(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(objeto)
}

func (conector *ObjetoPatrimonialCorpConector) tratarCalculoDepreciacao(w http.ResponseWriter, r *http.Request, idStr string) {
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID inválido.", http.StatusBadRequest)
		return
	}

	// Extrai a taxa enviada via Query Param (ex: ?taxa=0.10)
	taxaStr := r.URL.Query().Get("taxa")
	taxa, err := strconv.ParseFloat(taxaStr, 64)
	if err != nil || taxaStr == "" {
		http.Error(w, "É obrigatório fornecer uma taxa válida (ex: ?taxa=0.10).", http.StatusBadRequest)
		return
	}

	valorAtualizado, err := conector.casoDeUso.ObterValorDepreciado(id, taxa)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{"objeto_id": id, "valor_depreciado": valorAtualizado})
}

func (conector *ObjetoPatrimonialCorpConector) tratarEdicao(w http.ResponseWriter, r *http.Request, idStr string) {
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID inválido para edição.", http.StatusBadRequest)
		return
	}

	var dto EditarObjetoPatrimonialCorpDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		http.Error(w, "JSON inválido!", http.StatusBadRequest)
		return
	}

	err = conector.casoDeUso.Atualizar(id, dto.Nome, dto.ValorPatrimonial, dto.DataAquisicao)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Ativo corporativo atualizado com sucesso!"})
}

func (conector *ObjetoPatrimonialCorpConector) tratarExclusao(w http.ResponseWriter, r *http.Request, idStr string) {
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID inválido.", http.StatusBadRequest)
		return
	}

	if err := conector.casoDeUso.Deletar(id); err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}