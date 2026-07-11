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

type ObjetoPatrimonialPessConector struct {
	casoDeUso *caso_de_uso.ObjetoPatrimonialPessCasoDeUso
}

func NovoObjetoPatrimonialPessConector(uc *caso_de_uso.ObjetoPatrimonialPessCasoDeUso) *ObjetoPatrimonialPessConector {
	return &ObjetoPatrimonialPessConector{casoDeUso: uc}
}

type InserirObjetoPatrimonialPessDTO struct {
	Nome                         string    `json:"nome"`
	ValorPatrimonial             float64   `json:"valor_patrimonial"`
	DataAquisicao                time.Time `json:"data_aquisicao"`
	CategoriaObjetoPatrimonialID int       `json:"categoria_objeto_id"`
	UsuarioID                    int       `json:"usuario_id"`
}

type EditarObjetoPatrimonialPessDTO struct {
	Nome                         string    `json:"nome"`
	ValorPatrimonial             float64   `json:"valor_patrimonial"`
	DataAquisicao                time.Time `json:"data_aquisicao"`
	CategoriaObjetoPatrimonialID int       `json:"categoria_objeto_id"`
}

func (conector *ObjetoPatrimonialPessConector) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	partes := strings.Split(strings.Trim(r.URL.Path, "/"), "/")

	// CORRIGIDO: Valida o nome da rota base no primeiro índice
	if len(partes) == 0 || partes[0] != "objetos-patrimoniais-pess" {
		http.NotFound(w, r)
		return
	}

	switch r.Method {
	case http.MethodGet:
		// CORRIGIDO: Filtra rotas com base na posição correta dos elementos
		if len(partes) == 2 {
			conector.tratarBuscaPorID(w, r, partes[1])
		} else if len(partes) == 3 && partes[2] == "valorizacao" {
			conector.tratarCalculoValorizacao(w, r, partes[1])
		} else if len(partes) == 3 && partes[2] == "patrimonio" {
			conector.tratarListagemPorPatrimonio(w, r, partes[1])
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

func (conector *ObjetoPatrimonialPessConector) tratarCadastro(w http.ResponseWriter, r *http.Request) {
	var dto InserirObjetoPatrimonialPessDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		http.Error(w, "JSON enviado inválido!", http.StatusBadRequest)
		return
	}

	objeto := &entidade.ObjetoPatrimonialPess{
		Nome:                         dto.Nome,
		ValorPatrimonial:             dto.ValorPatrimonial,
		DataAquisicao:                dto.DataAquisicao,
		CategoriaObjetoPatrimonialID: dto.CategoriaObjetoPatrimonialID,
	}

	if err := conector.casoDeUso.Cadastrar(objeto, dto.UsuarioID); err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(objeto)
}

func (conector *ObjetoPatrimonialPessConector) tratarListagemPorPatrimonio(w http.ResponseWriter, r *http.Request, patrimonioIDStr string) {
	patrimonioID, err := strconv.Atoi(patrimonioIDStr)
	if err != nil {
		http.Error(w, "ID do patrimônio deve ser um número inteiro.", http.StatusBadRequest)
		return
	}

	lista, err := conector.casoDeUso.ListarPorPatrimonio(patrimonioID)
	if err != nil {
		http.Error(w, "Erro ao processar busca", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(lista)
}

func (conector *ObjetoPatrimonialPessConector) tratarBuscaPorID(w http.ResponseWriter, r *http.Request, idStr string) {
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID inválido.", http.StatusBadRequest)
		return
	}

	objeto, err := conector.casoDeUso.Buscar(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(objeto)
}

func (conector *ObjetoPatrimonialPessConector) tratarCalculoValorizacao(w http.ResponseWriter, r *http.Request, idStr string) {
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID inválido.", http.StatusBadRequest)
		return
	}

	taxaStr := r.URL.Query().Get("taxa")
	taxa, err := strconv.ParseFloat(taxaStr, 64)
	if err != nil || taxaStr == "" {
		http.Error(w, "É obrigatório fornecer o parâmetro de taxa anual (ex: ?taxa=0.08).", http.StatusBadRequest)
		return
	}

	valorProjetado, err := conector.casoDeUso.ObterProjecaoValorizacao(id, taxa)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{"objeto_id": id, "valor_projetado_acumulado": valorProjetado})
}

func (conector *ObjetoPatrimonialPessConector) tratarEdicao(w http.ResponseWriter, r *http.Request, idStr string) {
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID inválido para edição.", http.StatusBadRequest)
		return
	}

	var dto EditarObjetoPatrimonialPessDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		http.Error(w, "JSON inválido!", http.StatusBadRequest)
		return
	}

	err = conector.casoDeUso.Atualizar(id, dto.Nome, dto.ValorPatrimonial, dto.DataAquisicao, dto.CategoriaObjetoPatrimonialID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Ativo pessoal atualizado com sucesso!"})
}

func (conector *ObjetoPatrimonialPessConector) tratarExclusao(w http.ResponseWriter, r *http.Request, idStr string) {
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