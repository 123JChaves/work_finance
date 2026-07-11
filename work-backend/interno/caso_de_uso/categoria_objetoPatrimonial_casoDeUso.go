package caso_de_uso

import (
	"errors"
	"work-backend/interno/entidade"
	"work-backend/interno/regras"
)

type CategoriaObjetoPatrimonialCasoDeUso struct {
	repositorio entidade.CategoriaObjetoPatrimonialRepositorio
}

func NovoCategoriaObjetoPatrimonialCasoDeUso(repo entidade.CategoriaObjetoPatrimonialRepositorio) *CategoriaObjetoPatrimonialCasoDeUso {
	return &CategoriaObjetoPatrimonialCasoDeUso{repositorio: repo}
}

func (uc *CategoriaObjetoPatrimonialCasoDeUso) Cadastrar(cat *entidade.CategoriaObjetoPatrimonial) error {
	if err := regras.ValidarCategoriaObjetoPatrimonial(cat); err != nil {
		return err
	}

	existente, _ := uc.repositorio.FindByNome(cat.Nome)
	if existente != nil {
		return errors.New("já existe uma categoria patrimonial cadastrada com este nome")
	}

	return uc.repositorio.Create(cat)
}

func (uc *CategoriaObjetoPatrimonialCasoDeUso) Listar() ([]*entidade.CategoriaObjetoPatrimonial, error) {
	return uc.repositorio.FindAll()
}

func (uc *CategoriaObjetoPatrimonialCasoDeUso) Buscar(id int) (*entidade.CategoriaObjetoPatrimonial, error) {
	if id <= 0 {
		return nil, errors.New("o ID da categoria patrimonial deve ser um número válido")
	}
	return uc.repositorio.FindByID(id)
}

func (uc *CategoriaObjetoPatrimonialCasoDeUso) Atualizar(id int, novoNome string) error {
	cat, err := uc.Buscar(id)
	if err != nil {
		return err
	}

	cat.Nome = novoNome
	if err := regras.ValidarCategoriaObjetoPatrimonial(cat); err != nil {
		return err
	}

	existente, _ := uc.repositorio.FindByNome(cat.Nome)
	if existente != nil && existente.ID != cat.ID {
		return errors.New("já existe outra categoria patrimonial com este nome")
	}

	return uc.repositorio.Update(cat)
}

func (uc *CategoriaObjetoPatrimonialCasoDeUso) Deletar(id int) error {
	if _, err := uc.Buscar(id); err != nil {
		return err
	}
	return uc.repositorio.Delete(id)
}