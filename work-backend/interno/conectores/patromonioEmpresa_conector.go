package conectores

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"work-backend/interno/caso_de_uso"
)

type PatrimonioEmpresaConector struct {
	casoDeUso *caso_de_uso.PatrimonioEmpresaCasoDeUso
}

func NovoPatrimonioEmpresaConector(uc *caso_de_uso.PatrimonioEmpresaCasoDeUso) *PatrimonioEmpresaConector {
	return &PatrimonioEmpresaConector{casoDeUso: uc}
}

type InicializarPatrimonioEmpresaDTO struct {
	EmpresaID int `json:"empresa_id"`
}

func (conector *PatrimonioEmpresaConector) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	partes := strings.Split(strings.Trim(r.URL.Path, "/"), "/")

	// CORRIGIDO: Valida rota base usando o índice inicial
	if len(partes) == 0 || partes[0] != "patrimonios-empresa" {
		http.NotFound(w, r)
		return
	}

	switch r.Method {
	case http.MethodGet:
		// CORRIGIDO: Mapeia rota /patrimonios-empresa/total/:empresa_id de forma correta
		if len(partes) == 3 && partes[1] == "total" {
			conector.tratarTotalAtivos(w, r, partes[2])
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

func (conector *PatrimonioEmpresaConector) tratarInicializacao(w http.ResponseWriter, r *http.Request) {
	var dto InicializarPatrimonioEmpresaDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		http.Error(w, "JSON enviado inválido!", http.StatusBadRequest)
		return
	}

	if err := conector.casoDeUso.Inicializar(dto.EmpresaID); err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "Carteira patrimonial corporativa ativada com sucesso!"})
}

func (conector *PatrimonioEmpresaConector) tratarTotalAtivos(w http.ResponseWriter, r *http.Request, empresaIDStr string) {
	empresaID, err := strconv.Atoi(empresaIDStr)
	if err != nil {
		http.Error(w, "ID da empresa inválido.", http.StatusBadRequest)
		return
	}

	total, err := conector.casoDeUso.ObterValorConsolidado(empresaID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{"empresa_id": empresaID, "soma_total_ativos_corporativos": total})
}

func (conector *PatrimonioEmpresaConector) tratarExclusao(w http.ResponseWriter, r *http.Request, idStr string) {
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