package caso_de_uso

import (
	"errors"
	"work-backend/interno/entidade"
	"work-backend/interno/regras"
)

type CategoriaDespesaEmpresaCasoDeUso struct {
	repositorio entidade.CategoriaDespesaEmpresaRepositorio
}

func NovoCategoriaDespesaEmpresaCasoDeUso(repo entidade.CategoriaDespesaEmpresaRepositorio) *CategoriaDespesaEmpresaCasoDeUso {
	return &CategoriaDespesaEmpresaCasoDeUso{repositorio: repo}
}

func (uc *CategoriaDespesaEmpresaCasoDeUso) Cadastrar(cat *entidade.CategoriaDespesaEmpresa) error {
	if err := regras.ValidarCategoriaDespesaEmpresa(cat); err != nil {
		return err
	}

	existente, _ := uc.repositorio.FindByNome(cat.Nome)
	if existente != nil {
		return errors.New("já existe uma categoria de despesa corporativa cadastrada com este nome")
	}

	return uc.repositorio.Create(cat)
}

func (uc *CategoriaDespesaEmpresaCasoDeUso) Listar() ([]*entidade.CategoriaDespesaEmpresa, error) {
	return uc.repositorio.FindAll()
}

func (uc *CategoriaDespesaEmpresaCasoDeUso) Buscar(id int) (*entidade.CategoriaDespesaEmpresa, error) {
	if id <= 0 {
		return nil, errors.New("o ID da categoria deve ser um número válido")
	}
	return uc.repositorio.FindByID(id)
}

func (uc *CategoriaDespesaEmpresaCasoDeUso) Atualizar(id int, novoNome string) error {
	cat, err := uc.Buscar(id)
	if err != nil {
		return err
	}

	cat.Nome = novoNome
	if err := regras.ValidarCategoriaDespesaEmpresa(cat); err != nil {
		return err
	}

	existente, _ := uc.repositorio.FindByNome(cat.Nome)
	if existente != nil && existente.ID != cat.ID {
		return errors.New("já existe outra categoria de despesa com este nome")
	}

	return uc.repositorio.Update(cat)
}

func (uc *CategoriaDespesaEmpresaCasoDeUso) Deletar(id int) error {
	if _, err := uc.Buscar(id); err != nil {
		return err
	}
	return uc.repositorio.Delete(id)
}