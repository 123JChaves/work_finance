package conectores

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"work-backend/interno/caso_de_uso"
	"work-backend/interno/entidade"
)

type CategoriaDespesaUsuarioConector struct {
	casoDeUso *caso_de_uso.CategoriaDespesaUsuarioCasoDeUso
}

func NovoCategoriaDespesaUsuarioConector(uc *caso_de_uso.CategoriaDespesaUsuarioCasoDeUso) *CategoriaDespesaUsuarioConector {
	return &CategoriaDespesaUsuarioConector{casoDeUso: uc}
}

type EditarCategoriaDespesaUsuarioDTO struct {
	Nome string `json:"nome"`
}

func (conector *CategoriaDespesaUsuarioConector) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	partes := strings.Split(strings.Trim(r.URL.Path, "/"), "/")

	// Rota base esperada: /categorias-despesas-usuario
	if len(partes) == 0 || partes[0] != "categorias-despesas-usuario" {
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

func (conector *CategoriaDespesaUsuarioConector) tratarCadastro(w http.ResponseWriter, r *http.Request) {
	var cat entidade.CategoriaDespesaUsuario
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

func (conector *CategoriaDespesaUsuarioConector) tratarListagem(w http.ResponseWriter, r *http.Request) {
	lista, err := conector.casoDeUso.Listar()
	if err != nil {
		http.Error(w, "Erro ao buscar registros", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(lista)
}

func (conector *CategoriaDespesaUsuarioConector) tratarBuscaPorID(w http.ResponseWriter, r *http.Request, idStr string) {
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

func (conector *CategoriaDespesaUsuarioConector) tratarEdicao(w http.ResponseWriter, r *http.Request, idStr string) {
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID inválido.", http.StatusBadRequest)
		return
	}

	var dto EditarCategoriaDespesaUsuarioDTO
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

func (conector *CategoriaDespesaUsuarioConector) tratarExclusao(w http.ResponseWriter, r *http.Request, idStr string) {
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