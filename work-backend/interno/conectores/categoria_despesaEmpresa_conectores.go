package conectores

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"work-backend/interno/caso_de_uso"
	"work-backend/interno/entidade"
)

type CategoriaDespesaEmpresaConector struct {
	casoDeUso *caso_de_uso.CategoriaDespesaEmpresaCasoDeUso
}

func NovoCategoriaDespesaEmpresaConector(uc *caso_de_uso.CategoriaDespesaEmpresaCasoDeUso) *CategoriaDespesaEmpresaConector {
	return &CategoriaDespesaEmpresaConector{casoDeUso: uc}
}

type EditarCategoriaDespesaEmpresaDTO struct {
	Nome string `json:"nome"`
}

func (conector *CategoriaDespesaEmpresaConector) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	partes := strings.Split(strings.Trim(r.URL.Path, "/"), "/")

	if len(partes) == 0 || partes[0] != "categorias-despesas-empresa" {
		http.NotFound(w, r)
		return
	}

	switch r.Method {
	case http.MethodGet:
		if len(partes) == 1 {
			conector.tratarListagem(w, r)
		} else if len(partes) == 2 {
			conector.tratarBuscaPorID(w, r, partes[1])
		} else {
			http.NotFound(w, r)
		}
	case http.MethodPost:
		if len(partes) == 1 { conector.tratarCadastro(w, r) } else { http.NotFound(w, r) }
	case http.MethodPut:
		if len(partes) == 2 { conector.tratarEdicao(w, r, partes[1]) } else { http.Error(w, "ID não fornecido", http.StatusBadRequest) }
	case http.MethodDelete:
		if len(partes) == 2 { conector.tratarExclusao(w, r, partes[1]) } else { http.Error(w, "ID não fornecido", http.StatusBadRequest) }
	default:
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
	}
}

func (conector *CategoriaDespesaEmpresaConector) tratarCadastro(w http.ResponseWriter, r *http.Request) {
	var cat entidade.CategoriaDespesaEmpresa
	if err := json.NewDecoder(r.Body).Decode(&cat); err != nil {
		http.Error(w, "JSON enviado inválido!", http.StatusBadRequest)
		return
	}
	if err := conector.casoDeUso.Cadastrar(&cat); err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(cat)
}

func (conector *CategoriaDespesaEmpresaConector) tratarListagem(w http.ResponseWriter, r *http.Request) {
	lista, err := conector.casoDeUso.Listar()
	if err != nil {
		http.Error(w, "Erro ao buscar registros", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(lista)
}

// CORRIGIDO: Tipo alterado de CategoriaConector para CategoriaDespesaEmpresaConector
func (conector *CategoriaDespesaEmpresaConector) tratarBuscaPorID(w http.ResponseWriter, r *http.Request, idStr string) {
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID deve ser um número inteiro.", http.StatusBadRequest)
		return
	}
	cat, err := conector.casoDeUso.Buscar(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(cat)
}

func (conector *CategoriaDespesaEmpresaConector) tratarEdicao(w http.ResponseWriter, r *http.Request, idStr string) {
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID inválido.", http.StatusBadRequest)
		return
	}
	var dto EditarCategoriaDespesaEmpresaDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		http.Error(w, "JSON inválido para edição!", http.StatusBadRequest)
		return
	}
	if err := conector.casoDeUso.Atualizar(id, dto.Nome); err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Categoria atualizada com sucesso!"})
}

func (conector *CategoriaDespesaEmpresaConector) tratarExclusao(w http.ResponseWriter, r *http.Request, idStr string) {
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