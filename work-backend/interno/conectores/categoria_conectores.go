package conectores

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"work-backend/interno/caso_de_uso"
	"work-backend/interno/entidade"
)

type CategoriaConector struct {
	casoDeUso *caso_de_uso.CategoriaCasoDeUso
}

func NovoCategoriaConector(uc *caso_de_uso.CategoriaCasoDeUso) *CategoriaConector {
	return &CategoriaConector{casoDeUso: uc}
}

type EditarCategoriaDTO struct {
	Nome string `json:"nome"`
}

func (conector *CategoriaConector) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	partes := strings.Split(strings.Trim(r.URL.Path, "/"), "/")

	if len(partes) == 0 || partes[0] != "categorias" {
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
		if len(partes) == 2 { conector.tratarEdicao(w, r, partes[1]) } else { http.Error(w, "ID da categoria não fornecido na URL", http.StatusBadRequest) }
	case http.MethodDelete:
		if len(partes) == 2 { conector.tratarExclusao(w, r, partes[1]) } else { http.Error(w, "ID da categoria não fornecido na URL", http.StatusBadRequest) }
	default:
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
	}
}

func (conector *CategoriaConector) tratarCadastro(w http.ResponseWriter, r *http.Request) {
	var cat entidade.Categoria
	
	// PROTEÇÃO ADICIONADA: Valida falha na decodificação do corpo da requisição
	if erro := json.NewDecoder(r.Body).Decode(&cat); erro != nil {
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

func (conector *CategoriaConector) tratarListagem(w http.ResponseWriter, r *http.Request) {
	lista, erro := conector.casoDeUso.Listar()
	
	// PROTEÇÃO ADICIONADA: Trata falhas na busca interna
	if erro != nil {
		http.Error(w, "Erro ao buscar categorias", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(lista)
}

func (conector *CategoriaConector) tratarBuscaPorID(w http.ResponseWriter, r *http.Request, idStr string) {
	// PROTEÇÃO ADICIONADA: Valida se a string extraída da URL é um inteiro válido
	id, erro := strconv.Atoi(idStr)
	if erro != nil {
		http.Error(w, "ID inválido. Deve ser um número inteiro.", http.StatusBadRequest)
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

func (conector *CategoriaConector) tratarEdicao(w http.ResponseWriter, r *http.Request, idStr string) {
	// PROTEÇÃO ADICIONADA: Valida inteiros e decodificação do DTO de entrada
	id, erro := strconv.Atoi(idStr)
	if erro != nil {
		http.Error(w, "ID inválido. Deve ser um número inteiro.", http.StatusBadRequest)
		return
	}

	var dto EditarCategoriaDTO
	if erro := json.NewDecoder(r.Body).Decode(&dto); erro != nil {
		http.Error(w, "JSON enviado inválido para edição!", http.StatusBadRequest)
		return
	}

	if err := conector.casoDeUso.Atualizar(id, dto.Nome); err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Categoria atualizada com sucesso!"})
}

func (conector *CategoriaConector) tratarExclusao(w http.ResponseWriter, r *http.Request, idStr string) {
	// PROTEÇÃO ADICIONADA: Valida string do ID antes de invocar o caso de uso
	id, erro := strconv.Atoi(idStr)
	if erro != nil {
		http.Error(w, "ID inválido. Deve ser um número inteiro.", http.StatusBadRequest)
		return
	}

	if err := conector.casoDeUso.Deletar(id); err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}