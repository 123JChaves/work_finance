package conectores

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"work-backend/interno/caso_de_uso"
	"work-backend/interno/entidade"
)

type UsuarioConector struct {
	casoDeUso *caso_de_uso.UsuarioCasoDeUso
}

func NovoUsuarioConector(uc *caso_de_uso.UsuarioCasoDeUso) *UsuarioConector {
	return &UsuarioConector{casoDeUso: uc}
}

type EditarUsuarioDTO struct {
	Nome        string  `json:"nome"`
	Email       string  `json:"email"`
	Cpf         string  `json:"cpf"`
	SalarioBase float64 `json:"salario_base"`
}

func (conector *UsuarioConector) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	partes := strings.Split(strings.Trim(r.URL.Path, "/"), "/")

	if len(partes) == 0 || partes[0] != "usuarios" {
		http.NotFound(w, r)
		return
	}

	switch r.Method {
	case http.MethodGet:
		if len(partes) == 1 {
			conector.tratarListagem(w, r)
		} else if len(partes) == 2 {
			conector.tratarBuscaPorID(w, r, partes[1])
		} else if len(partes) == 3 && partes[2] == "balanco" {
			conector.tratarBalancoPatrimonial(w, r, partes[1])
		} else {
			http.NotFound(w, r)
		}
	case http.MethodPost:
		if len(partes) == 1 {
			conector.tratarCadastro(w, r)
		} else if len(partes) == 3 && partes[2] == "fechamento" {
			conector.tratarFechamentoMensal(w, r, partes[1])
		} else {
			http.NotFound(w, r)
		}
	case http.MethodPut:
		if len(partes) == 2 {
			conector.tratarEdicao(w, r, partes[1])
		} else {
			http.Error(w, "ID ausente", http.StatusBadRequest)
		}
	case http.MethodDelete:
		if len(partes) == 2 {
			conector.tratarExclusao(w, r, partes[1])
		} else {
			http.Error(w, "ID ausente", http.StatusBadRequest)
		}
	default:
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
	}
}

func (conector *UsuarioConector) tratarCadastro(w http.ResponseWriter, r *http.Request) {
	var usuario entidade.Usuario
	if err := json.NewDecoder(r.Body).Decode(&usuario); err != nil {
		http.Error(w, "JSON enviado inválido!", http.StatusBadRequest)
		return
	}
	if err := conector.casoDeUso.Cadastrar(&usuario); err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}
	usuario.Senha = "[PROTEGIDO]"
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(usuario)
}

func (conector *UsuarioConector) tratarListagem(w http.ResponseWriter, r *http.Request) {
	lista, err := conector.casoDeUso.Listar()
	if err != nil {
		http.Error(w, "Erro ao listar", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(lista)
}

func (conector *UsuarioConector) tratarBuscaPorID(w http.ResponseWriter, r *http.Request, idStr string) {
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID deve ser numérico.", http.StatusBadRequest)
		return
	}
	usuario, err := conector.casoDeUso.Buscar(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	usuario.Senha = "[PROTEGIDO]"
	json.NewEncoder(w).Encode(usuario)
}

func (conector *UsuarioConector) tratarEdicao(w http.ResponseWriter, r *http.Request, idStr string) {
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID deve ser numérico.", http.StatusBadRequest)
		return
	}
	var dto EditarUsuarioDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		http.Error(w, "JSON inválido!", http.StatusBadRequest)
		return
	}
	if err := conector.casoDeUso.Atualizar(id, dto.Nome, dto.Email, dto.Cpf, dto.SalarioBase); err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}
	json.NewEncoder(w).Encode(map[string]string{"message": "Usuário atualizado com sucesso!"})
}

func (conector *UsuarioConector) tratarFechamentoMensal(w http.ResponseWriter, r *http.Request, idStr string) {
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}
	novaLiquidez, err := conector.casoDeUso.ProcessarFechamentoMensal(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}
	json.NewEncoder(w).Encode(map[string]interface{}{"usuario_id": id, "novo_total_liquido": novaLiquidez})
}

func (conector *UsuarioConector) tratarBalancoPatrimonial(w http.ResponseWriter, r *http.Request, idStr string) {
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}
	balanco, err := conector.casoDeUso.ObterBalançoPatrimonial(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(balanco)
}

func (conector *UsuarioConector) tratarExclusao(w http.ResponseWriter, r *http.Request, idStr string) {
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