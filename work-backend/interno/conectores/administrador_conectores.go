package conectores

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"work-backend/interno/caso_de_uso"
	"work-backend/interno/entidade"
)

type AdministradorConector struct {
	casoDeUso *caso_de_uso.Administrador_casoDeUso
}

func NovoAdministradorConector(uc *caso_de_uso.Administrador_casoDeUso) *AdministradorConector {
	return &AdministradorConector{casoDeUso: uc}
}

type EditarAdministradorDTO struct {
	Nome      string `json:"nome"`
	Email     string `json:"email"`
	NovaSenha string `json:"novaSenha"`
}

func (conector *AdministradorConector) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	partes := strings.Split(strings.Trim(r.URL.Path, "/"), "/")

	if len(partes) == 0 || partes[0] != "administradores" {
		http.NotFound(w, r)
		return
	}

	switch r.Method {
	case http.MethodGet:
		if len(partes) == 1 { conector.tratarListagem(w, r) } else { conector.tratarBuscaPorID(w, r, partes[1]) }
	case http.MethodPost:
		conector.tratarCadastro(w, r)
	case http.MethodPut:
		conector.tratarEdicao(w, r, partes[1])
	case http.MethodDelete:
		conector.tratarExclusao(w, r, partes[1])
	default:
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
	}
}

func (conector *AdministradorConector) tratarCadastro(w http.ResponseWriter, r *http.Request) {
	var admin entidade.Administrador
	if json.NewDecoder(r.Body).Decode(&admin) != nil {
		http.Error(w, "JSON inválido!", http.StatusBadRequest)
		return
	}
	if err := conector.casoDeUso.Cadastrar(&admin); err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}
	admin.Senha = "[PROTEGIDO]"
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(admin)
}

func (conector *AdministradorConector) tratarListagem(w http.ResponseWriter, r *http.Request) {
	lista, _ := conector.casoDeUso.Listar()
	for _, a := range lista { a.Senha = "[PROTEGIDO]" }
	json.NewEncoder(w).Encode(lista)
}

func (conector *AdministradorConector) tratarBuscaPorID(w http.ResponseWriter, r *http.Request, idStr string) {
	id, _ := strconv.Atoi(idStr)
	admin, err := conector.casoDeUso.Buscar(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	admin.Senha = "[PROTEGIDO]"
	json.NewEncoder(w).Encode(admin)
}

func (conector *AdministradorConector) tratarEdicao(w http.ResponseWriter, r *http.Request, idStr string) {
	id, _ := strconv.Atoi(idStr)
	var dto EditarAdministradorDTO
	json.NewDecoder(r.Body).Decode(&dto)
	if err := conector.casoDeUso.Atualizar(id, dto.Nome, dto.Email, dto.NovaSenha); err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}
	json.NewEncoder(w).Encode(map[string]string{"message": "Sucesso!"})
}

func (conector *AdministradorConector) tratarExclusao(w http.ResponseWriter, r *http.Request, idStr string) {
	id, _ := strconv.Atoi(idStr)
	if err := conector.casoDeUso.Deletar(id); err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
