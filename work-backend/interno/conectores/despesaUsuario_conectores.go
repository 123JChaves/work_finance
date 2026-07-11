package conectores

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"
	"work-backend/interno/caso_de_uso"
	"work-backend/interno/entidade"
)

type DespesaUsuarioConector struct {
	casoDeUso *caso_de_uso.DespesaUsuarioCasoDeUso
}

func NovoDespesaUsuarioConector(uc *caso_de_uso.DespesaUsuarioCasoDeUso) *DespesaUsuarioConector {
	return &DespesaUsuarioConector{casoDeUso: uc}
}

type InserirDespesaUsuarioDTO struct {
	Descricao                 string    `json:"descricao"`
	Valor                     float64   `json:"valor"`
	Data                      time.Time `json:"data"`
	CategoriaDespesaUsuarioID int       `json:"categoria_despesa_usuario_id"`
	UsuarioID                 int       `json:"usuario_id"`
}

type EditarDespesaUsuarioDTO struct {
	Descricao                 string    `json:"descricao"`
	Valor                     float64   `json:"valor"`
	Data                      time.Time `json:"data"`
	CategoriaDespesaUsuarioID int       `json:"categoria_despesa_usuario_id"`
}

func (conector *DespesaUsuarioConector) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	partes := strings.Split(strings.Trim(r.URL.Path, "/"), "/")

	if len(partes) == 0 || partes[0] != "despesas-usuario" {
		http.NotFound(w, r)
		return
	}

	switch r.Method {
	case http.MethodGet:
		if len(partes) == 2 {
			conector.tratarBuscaPorID(w, r, partes[1])
		} else if len(partes) == 3 && partes[1] == "usuario" {
			conector.tratarListagemPorUsuario(w, r, partes[2])
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

func (conector *DespesaUsuarioConector) tratarCadastro(w http.ResponseWriter, r *http.Request) {
	var dto InserirDespesaUsuarioDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		http.Error(w, "JSON enviado inválido!", http.StatusBadRequest)
		return
	}

	despesa := &entidade.DespesaUsuario{
		Descricao:                 dto.Descricao,
		Valor:                     dto.Valor,
		Data:                      dto.Data,
		CategoriaDespesaUsuarioID: dto.CategoriaDespesaUsuarioID,
		UsuarioID:                 dto.UsuarioID,
	}

	if err := conector.casoDeUso.Cadastrar(despesa); err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(despesa)
}

func (conector *DespesaUsuarioConector) tratarListagemPorUsuario(w http.ResponseWriter, r *http.Request, usuarioIDStr string) {
	usuarioID, err := strconv.Atoi(usuarioIDStr)
	if err != nil {
		http.Error(w, "ID do usuário inválido.", http.StatusBadRequest)
		return
	}

	lista, err := conector.casoDeUso.ListarPorUsuario(usuarioID)
	if err != nil {
		http.Error(w, "Erro ao processar busca", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(lista)
}

func (conector *DespesaUsuarioConector) tratarBuscaPorID(w http.ResponseWriter, r *http.Request, idStr string) {
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID inválido.", http.StatusBadRequest)
		return
	}

	despesa, err := conector.casoDeUso.Buscar(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(despesa)
}

func (conector *DespesaUsuarioConector) tratarEdicao(w http.ResponseWriter, r *http.Request, idStr string) {
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID com formato incorreto.", http.StatusBadRequest)
		return
	}

	var dto EditarDespesaUsuarioDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		http.Error(w, "JSON inválido!", http.StatusBadRequest)
		return
	}

	err = conector.casoDeUso.Atualizar(id, dto.Descricao, dto.Valor, dto.CategoriaDespesaUsuarioID, dto.Data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Despesa atualizada com sucesso!"})
}

func (conector *DespesaUsuarioConector) tratarExclusao(w http.ResponseWriter, r *http.Request, idStr string) {
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