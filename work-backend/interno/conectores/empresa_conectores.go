package conectores

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"work-backend/interno/caso_de_uso"
	"work-backend/interno/entidade"
)

type EmpresaConector struct {
	casoDeUso *caso_de_uso.EmpresaCasoDeUso
}

func NovoEmpresaConector(uc *caso_de_uso.EmpresaCasoDeUso) *EmpresaConector {
	return &EmpresaConector{casoDeUso: uc}
}

type InserirEmpresaDTO struct {
	Nome string `json:"nome"`
	Cnpj string `json:"cnpj"`
}

func (conector *EmpresaConector) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	partes := strings.Split(strings.Trim(r.URL.Path, "/"), "/")

	if len(partes) == 0 || partes[0] != "empresas" {
		http.NotFound(w, r)
		return
	}

	switch r.Method {
	case http.MethodGet:
		if len(partes) == 1 {
			conector.tratarListagem(w, r)
		} else if len(partes) == 2 {
			conector.tratarBuscaPorID(w, r, partes[1])
		} else if len(partes) == 3 && partes[2] == "lucro" {
			conector.tratarLucroLiquido(w, r, partes[1])
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

func (conector *EmpresaConector) tratarCadastro(w http.ResponseWriter, r *http.Request) {
	var dto InserirEmpresaDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		http.Error(w, "JSON enviado inválido!", http.StatusBadRequest)
		return
	}

	empresa := &entidade.Empresa{Nome: dto.Nome, Cnpj: dto.Cnpj}
	if err := conector.casoDeUso.Cadastrar(empresa); err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(empresa)
}

func (conector *EmpresaConector) tratarListagem(w http.ResponseWriter, r *http.Request) {
	lista, err := conector.casoDeUso.Listar()
	if err != nil {
		http.Error(w, "Erro interno", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(lista)
}

func (conector *EmpresaConector) tratarBuscaPorID(w http.ResponseWriter, r *http.Request, idStr string) {
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID inválido.", http.StatusBadRequest)
		return
	}

	empresa, err := conector.casoDeUso.Buscar(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(empresa)
}

func (conector *EmpresaConector) tratarLucroLiquido(w http.ResponseWriter, r *http.Request, idStr string) {
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID inválido.", http.StatusBadRequest)
		return
	}

	lucro, err := conector.casoDeUso.ObterDRE(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{"empresa_id": id, "lucro_liquido": lucro})
}

func (conector *EmpresaConector) tratarEdicao(w http.ResponseWriter, r *http.Request, idStr string) {
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID inválido para edição.", http.StatusBadRequest)
		return
	}

	var dto InserirEmpresaDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		http.Error(w, "JSON inválido!", http.StatusBadRequest)
		return
	}

	if err := conector.casoDeUso.Atualizar(id, dto.Nome, dto.Cnpj); err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Dados corporativos atualizados!"})
}

func (conector *EmpresaConector) tratarExclusao(w http.ResponseWriter, r *http.Request, idStr string) {
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