package conectores

import (
	"encoding/json"
	"net/http"
	"work-backend/interno/caso_de_uso"
	"work-backend/interno/entidade"
)

type AdministradorConector struct {
	casoDeUso *caso_de_uso.Administrador_casoDeUso
}

func NovoAdministradorConector(uc *caso_de_uso.Administrador_casoDeUso) *AdministradorConector {
	return &AdministradorConector{casoDeUso: uc}
}

func (conector *AdministradorConector) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Tolerância para barras extras na rota (ex: /administradores/)
	if r.URL.Path != "/administradores" && r.URL.Path != "/administradores/" {
		http.NotFound(w, r)
		return
	}

	// Direciona a requisição baseada no método HTTP
	switch r.Method {
	case http.MethodGet:
		conector.tratarListagem(w, r)
	case http.MethodPost:
		conector.tratarCadastro(w, r)
	default:
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
	}
}

func (conector *AdministradorConector) tratarCadastro(w http.ResponseWriter, r *http.Request) {
	var administrador entidade.Administrador

	if erro := json.NewDecoder(r.Body).Decode(&administrador); erro != nil {
		http.Error(w, "JSON enviado inválido!", http.StatusBadRequest)
		return
	}

	if erro := conector.casoDeUso.Cadastrar(&administrador); erro != nil {
		http.Error(w, erro.Error(), http.StatusUnprocessableEntity)
		return
	}

	administrador.Senha = "[PROTEGIDO]"

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(administrador)
}

func (conector *AdministradorConector) tratarListagem(w http.ResponseWriter, r *http.Request) {
	administradores, erro := conector.casoDeUso.Listar()
	if erro != nil {
		http.Error(w, "Erro ao buscar administradores", http.StatusInternalServerError)
		return
	}

	// Mascara as senhas de todos os administradores listados por segurança
	for _, admin := range administradores {
		admin.Senha = "[PROTEGIDO]"
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(administradores)
}