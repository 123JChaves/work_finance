package conectores

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"work-backend/interno/caso_de_uso"
	"work-backend/interno/entidade"
)

type ClienteConector struct {
	casoDeUso *caso_de_uso.ClienteCasoDeUso
}

func NovoClienteConector(uc *caso_de_uso.ClienteCasoDeUso) *ClienteConector {
	return &ClienteConector{casoDeUso: uc}
}

type InserirClienteDTO struct {
	Nome      string  `json:"nome"`
	Email     string  `json:"email"`
	CpfCnpj   string  `json:"cpf_cnpj"`
	Contato   *string `json:"contato"`
	EmpresaID int     `json:"empresa_id"`
}

type EditarClienteDTO struct {
	Nome    string  `json:"nome"`
	Email   string  `json:"email"`
	CpfCnpj string  `json:"cpf_cnpj"`
	Contato *string `json:"contato"`
}

func (conector *ClienteConector) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	partes := strings.Split(strings.Trim(r.URL.Path, "/"), "/")

	if len(partes) == 0 || partes[0] != "clientes" {
		http.NotFound(w, r)
		return
	}

	switch r.Method {
	case http.MethodGet:
		if len(partes) == 1 {
			http.NotFound(w, r)
		} else if len(partes) == 2 {
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

func (conector *ClienteConector) tratarCadastro(w http.ResponseWriter, r *http.Request) {
	var dto InserirClienteDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		http.Error(w, "JSON enviado inválido!", http.StatusBadRequest)
		return
	}

	cliente := &entidade.Cliente{
		Nome:      dto.Nome,
		Email:     dto.Email,
		CpfCnpj:   dto.CpfCnpj,
		Contato:   dto.Contato,
		EmpresaID: dto.EmpresaID,
	}

	if err := conector.casoDeUso.Cadastrar(cliente); err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(cliente)
}

func (conector *ClienteConector) tratarListagemPorEmpresa(w http.ResponseWriter, r *http.Request, empresaIDStr string) {
	empresaID, err := strconv.Atoi(empresaIDStr)
	if err != nil {
		http.Error(w, "ID da empresa deve ser um número inteiro.", http.StatusBadRequest)
		return
	}

	lista, err := conector.casoDeUso.ListarPorEmpresa(empresaID)
	if err != nil {
		http.Error(w, "Erro ao buscar clientes", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(lista)
}

func (conector *ClienteConector) tratarBuscaPorID(w http.ResponseWriter, r *http.Request, idStr string) {
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID inválido.", http.StatusBadRequest)
		return
	}

	cliente, err := conector.casoDeUso.Buscar(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(cliente)
}

func (conector *ClienteConector) tratarEdicao(w http.ResponseWriter, r *http.Request, idStr string) {
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID inválido para edição.", http.StatusBadRequest)
		return
	}

	var dto EditarClienteDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		http.Error(w, "JSON inválido!", http.StatusBadRequest)
		return
	}

	err = conector.casoDeUso.Atualizar(id, dto.Nome, dto.Email, dto.CpfCnpj, dto.Contato)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Dados do cliente atualizados com sucesso!"})
}

func (conector *ClienteConector) tratarExclusao(w http.ResponseWriter, r *http.Request, idStr string) {
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