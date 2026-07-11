package caso_de_uso

import (
	"errors"
	"work-backend/interno/entidade"
	"work-backend/interno/regras"
)

type CategoriaCasoDeUso struct {
	repositorio entidade.CategoriaRepositorio
}

func NovoCategoriaCasoDeUso(repo entidade.CategoriaRepositorio) *CategoriaCasoDeUso {
	return &CategoriaCasoDeUso{repositorio: repo}
}

func (uc *CategoriaCasoDeUso) Cadastrar(cat *entidade.Categoria) error {
	// 1. Invoca a função pura de negócio isolada na pasta de regras
	if err := regras.ValidarCategoria(cat); err != nil {
		return err
	}

	// 2. Proteção de integridade contra duplicidade de nomes
	existente, _ := uc.repositorio.FindByNome(cat.Nome)
	if existente != nil {
		return errors.New("já existe uma categoria cadastrada com este nome")
	}

	return uc.repositorio.Create(cat)
}

func (uc *CategoriaCasoDeUso) Listar() ([]*entidade.Categoria, error) {
	return uc.repositorio.FindAll()
}

func (uc *CategoriaCasoDeUso) Buscar(id int) (*entidade.Categoria, error) {
	if id <= 0 {
		return nil, errors.New("o ID da categoria deve ser um número válido")
	}
	return uc.repositorio.FindByID(id)
}

func (uc *CategoriaCasoDeUso) Atualizar(id int, novoNome string) error {
	cat, err := uc.Buscar(id)
	if err != nil {
		return err
	}

	cat.Nome = novoNome
	
	// Executa a validação pura para o novo nome atribuído
	if err := regras.ValidarCategoria(cat); err != nil {
		return err
	}

	// Impede que altere o nome para o de outra categoria existente no banco
	existente, _ := uc.repositorio.FindByNome(cat.Nome)
	if existente != nil && existente.ID != cat.ID {
		return errors.New("já existe outra categoria com este nome")
	}

	return uc.repositorio.Update(cat)
}

func (uc *CategoriaCasoDeUso) Deletar(id int) error {
	if _, err := uc.Buscar(id); err != nil {
		return err
	}
	return uc.repositorio.Delete(id)
}