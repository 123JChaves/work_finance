package conectores

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"work-backend/interno/caso_de_uso"
)

type PatrimonioConector struct {
	casoDeUso *caso_de_uso.PatrimonioCasoDeUso
}

func NovoPatrimonioConector(uc *caso_de_uso.PatrimonioCasoDeUso) *PatrimonioConector {
	return &PatrimonioConector{casoDeUso: uc}
}

type InicializarPatrimonioDTO struct {
	UsuarioID int `json:"usuario_id"`
}

func (conector *PatrimonioConector) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	partes := strings.Split(strings.Trim(r.URL.Path, "/"), "/")

	if len(partes) == 0 || partes[0] != "patrimonios" {
		http.NotFound(w, r)
		return
	}

	switch r.Method {
	case http.MethodGet:
		if len(partes) == 3 && partes[2] == "total" {
			conector.tratarTotalBens(w, r, partes[1])
		} else {
			http.NotFound(w, r)
		}
	case http.MethodPost:
		if len(partes) == 1 { conector.tratarInicializacao(w, r) } else { http.NotFound(w, r) }
	case http.MethodDelete:
		if len(partes) == 2 { conector.tratarExclusao(w, r, partes[1]) } else { http.Error(w, "ID ausente", http.StatusBadRequest) }
	default:
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
	}
}

func (conector *PatrimonioConector) tratarInicializacao(w http.ResponseWriter, r *http.Request) {
	var dto InicializarPatrimonioDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		http.Error(w, "JSON enviado inválido!", http.StatusBadRequest)
		return
	}

	if err := conector.casoDeUso.Inicializar(dto.UsuarioID); err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "Carteira patrimonial pessoal ativada com sucesso!"})
}

func (conector *PatrimonioConector) tratarTotalBens(w http.ResponseWriter, r *http.Request, usuarioIDStr string) {
	usuarioID, err := strconv.Atoi(usuarioIDStr)
	if err != nil {
		http.Error(w, "ID do usuário inválido.", http.StatusBadRequest)
		return
	}

	total, err := conector.casoDeUso.ObterValorConsolidado(usuarioID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{"usuario_id": usuarioID, "soma_total_bens_pessoais": total})
}

func (conector *PatrimonioConector) tratarExclusao(w http.ResponseWriter, r *http.Request, idStr string) {
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