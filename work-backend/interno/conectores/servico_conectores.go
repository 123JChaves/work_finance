package conectores

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"work-backend/interno/caso_de_uso"
	"work-backend/interno/entidade"
)

type ServicoConector struct {
	casoDeUso *caso_de_uso.Servico_casoDeUso
}

func NovoServicoConector(uc *caso_de_uso.Servico_casoDeUso) *ServicoConector {
	return &ServicoConector{casoDeUso: uc}
}

type InserirServicoDTO struct {
	Nome        string  `json:"nome"`
	Valor       float64 `json:"valor"`
	CategoriaID int     `json:"categoria_id"`
	EmpresaID   int     `json:"empresa_id"`
}

type EditarServicoDTO struct {
	Nome        string  `json:"nome"`
	Valor       float64 `json:"valor"` // Alterado de string para float64 para bater com as regras
	CategoriaID int     `json:"categoria_id"`
}

func (conector *ServicoConector) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	partes := strings.Split(strings.Trim(r.URL.Path, "/"), "/")

	if len(partes) == 0 || partes[0] != "servicos" {
		http.NotFound(w, r)
		return
	}

	switch r.Method {
	case http.MethodGet:
		// Exemplo de rota: /servicos/relatorio/2026/07
		if len(partes) == 4 && partes[1] == "relatorio" {
			conector.tratarRelatorio(w, r, partes[2], partes[3]) // CORRIGIDO: Passando ano e mes extraídos
		} else if len(partes) == 1 {
			conector.tratarListagem(w, r)
		} else if len(partes) == 2 {
			conector.tratarBuscaPorID(w, r, partes[1])
		} else {
			http.NotFound(w, r)
		}
	case http.MethodPost:
		if len(partes) == 1 {
			conector.tratarCadastro(w, r)
		} else {
			http.NotFound(w, r)
		}
	case http.MethodPut:
		if len(partes) == 2 {
			conector.tratarEdicao(w, r, partes[1])
		} else {
			http.Error(w, "ID do serviço não fornecido na URL", http.StatusBadRequest)
		}
	case http.MethodDelete:
		if len(partes) == 2 {
			conector.tratarExclusao(w, r, partes[1])
		} else {
			http.Error(w, "ID do serviço não fornecido na URL", http.StatusBadRequest)
		}
	default:
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
	}
}

func (conector *ServicoConector) tratarCadastro(w http.ResponseWriter, r *http.Request) {
	var dto InserirServicoDTO
	if erro := json.NewDecoder(r.Body).Decode(&dto); erro != nil {
		http.Error(w, "JSON enviado inválido!", http.StatusBadRequest)
		return
	}

	servico := &entidade.Servico{
		Nome:        dto.Nome,
		Valor:       dto.Valor,
		CategoriaID: dto.CategoriaID,
		EmpresaID:   dto.EmpresaID,
	}

	if err := conector.casoDeUso.Cadastrar(servico); err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(servico)
}

func (conector *ServicoConector) tratarListagem(w http.ResponseWriter, r *http.Request) {
	lista, err := conector.casoDeUso.Listar()
	if err != nil {
		http.Error(w, "Erro ao buscar serviços", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(lista)
}

func (conector *ServicoConector) tratarBuscaPorID(w http.ResponseWriter, r *http.Request, idStr string) {
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID inválido.", http.StatusBadRequest)
		return
	}

	servico, err := conector.casoDeUso.Buscar(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(servico)
}

func (conector *ServicoConector) tratarRelatorio(w http.ResponseWriter, r *http.Request, anoStr, mesStr string) {
	ano, errAno := strconv.Atoi(anoStr)
	mes, errMes := strconv.Atoi(mesStr)
	if errAno != nil || errMes != nil {
		http.Error(w, "Ano e mês devem ser números inteiros válidos na URL.", http.StatusBadRequest)
		return
	}

	relatorio, err := conector.casoDeUso.ObterDadosRelatorio(ano, mes)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(relatorio)
}

func (conector *ServicoConector) tratarEdicao(w http.ResponseWriter, r *http.Request, idStr string) {
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID inválido.", http.StatusBadRequest)
		return
	}

	var dto EditarServicoDTO
	if erro := json.NewDecoder(r.Body).Decode(&dto); erro != nil {
		http.Error(w, "JSON inválido!", http.StatusBadRequest)
		return
	}

	// CORRIGIDO: Repassando o valor como float64 diretamente para bater com a assinatura do caso de uso
	err = conector.casoDeUso.Atualizar(id, dto.Nome, dto.Valor, dto.CategoriaID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Serviço atualizado com sucesso!"})
}

func (conector *ServicoConector) tratarExclusao(w http.ResponseWriter, r *http.Request, idStr string) {
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