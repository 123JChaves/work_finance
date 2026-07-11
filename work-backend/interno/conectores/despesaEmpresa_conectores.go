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

type DespesaEmpresaConector struct {
	casoDeUso *caso_de_uso.DespesaEmpresaCasoDeUso
}

func NovoDespesaEmpresaConector(uc *caso_de_uso.DespesaEmpresaCasoDeUso) *DespesaEmpresaConector {
	return &DespesaEmpresaConector{casoDeUso: uc}
}

type InserirDespesaEmpresaDTO struct {
	Descricao                 string    `json:"descricao"`
	Valor                     float64   `json:"valor"`
	Data                      time.Time `json:"data"`
	CategoriaDespesaEmpresaID int       `json:"categoria_despesa_empresa_id"`
	EmpresaID                 int       `json:"empresa_id"`
}

type EditarDespesaEmpresaDTO struct {
	Descricao                 string    `json:"descricao"`
	Valor                     float64   `json:"valor"`
	Data                      time.Time `json:"data"`
	CategoriaDespesaEmpresaID int       `json:"categoria_despesa_empresa_id"`
}

func (conector *DespesaEmpresaConector) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	partes := strings.Split(strings.Trim(r.URL.Path, "/"), "/")

	// CORRIGIDO: Acessando a primeira posição da string [0] para validação da rota base
	if len(partes) == 0 || partes[0] != "despesas-empresa" {
		http.NotFound(w, r)
		return
	}

	switch r.Method {
	case http.MethodGet:
		// CORRIGIDO: Passando o índice correto partes[1] e tratando comparação na posição partes[1]
		if len(partes) == 2 {
			conector.tratarBuscaPorID(w, r, partes[1])
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

func (conector *DespesaEmpresaConector) tratarCadastro(w http.ResponseWriter, r *http.Request) {
	var dto InserirDespesaEmpresaDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		http.Error(w, "JSON enviado inválido!", http.StatusBadRequest)
		return
	}

	despesa := &entidade.DespesaEmpresa{
		Descricao:                 dto.Descricao,
		Valor:                     dto.Valor,
		Data:                      dto.Data,
		CategoriaDespesaEmpresaID: dto.CategoriaDespesaEmpresaID,
		EmpresaID:                 dto.EmpresaID,
	}

	if err := conector.casoDeUso.Cadastrar(despesa); err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(despesa)
}

func (conector *DespesaEmpresaConector) tratarListagemPorEmpresa(w http.ResponseWriter, r *http.Request, empresaIDStr string) {
	empresaID, err := strconv.Atoi(empresaIDStr)
	if err != nil {
		http.Error(w, "ID da empresa deve ser um número inteiro.", http.StatusBadRequest)
		return
	}

	lista, err := conector.casoDeUso.ListarPorEmpresa(empresaID)
	if err != nil {
		http.Error(w, "Erro ao buscar despesas", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(lista)
}

func (conector *DespesaEmpresaConector) tratarBuscaPorID(w http.ResponseWriter, r *http.Request, idStr string) {
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID inválido.", http.StatusBadRequest)
		return
	}

	despesa, err := conector.casoDeUso.Buscar(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(despesa)
}

func (conector *DespesaEmpresaConector) tratarEdicao(w http.ResponseWriter, r *http.Request, idStr string) {
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID inválido para edição.", http.StatusBadRequest)
		return
	}

	var dto EditarDespesaEmpresaDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		http.Error(w, "JSON inválido!", http.StatusBadRequest)
		return
	}

	err = conector.casoDeUso.Atualizar(id, dto.Descricao, dto.Valor, dto.CategoriaDespesaEmpresaID, dto.Data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Despesa corporativa atualizada com sucesso!"})
}

func (conector *DespesaEmpresaConector) tratarExclusao(w http.ResponseWriter, r *http.Request, idStr string) {
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID inválido para exclusão.", http.StatusBadRequest)
		return
	}

	if err := conector.casoDeUso.Deletar(id); err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}