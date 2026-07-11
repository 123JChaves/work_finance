package caso_de_uso

import (
	"errors"
	"work-backend/interno/entidade"
	"work-backend/interno/regras"
)

type CategoriaDespesaUsuarioCasoDeUso struct {
	repositorio entidade.CategoriaDespesaUsuarioRepositorio
}

func NovoCategoriaDespesaUsuarioCasoDeUso(repo entidade.CategoriaDespesaUsuarioRepositorio) *CategoriaDespesaUsuarioCasoDeUso {
	return &CategoriaDespesaUsuarioCasoDeUso{repositorio: repo}
}

func (uc *CategoriaDespesaUsuarioCasoDeUso) Cadastrar(cat *entidade.CategoriaDespesaUsuario) error {
	if err := regras.ValidarCategoriaDespesaUsuario(cat); err != nil {
		return err
	}

	existente, _ := uc.repositorio.FindByNome(cat.Nome)
	if existente != nil {
		return errors.New("já existe uma categoria cadastrada com este nome")
	}

	return uc.repositorio.Create(cat)
}

func (uc *CategoriaDespesaUsuarioCasoDeUso) Listar() ([]*entidade.CategoriaDespesaUsuario, error) {
	return uc.repositorio.FindAll()
}

func (uc *CategoriaDespesaUsuarioCasoDeUso) Buscar(id int) (*entidade.CategoriaDespesaUsuario, error) {
	if id <= 0 {
		return nil, errors.New("o ID da categoria deve ser um número válido")
	}
	return uc.repositorio.FindByID(id)
}

func (uc *CategoriaDespesaUsuarioCasoDeUso) Atualizar(id int, novoNome string) error {
	cat, err := uc.Buscar(id)
	if err != nil {
		return err
	}

	cat.Nome = novoNome
	if err := regras.ValidarCategoriaDespesaUsuario(cat); err != nil {
		return err
	}

	existente, _ := uc.repositorio.FindByNome(cat.Nome)
	if existente != nil && existente.ID != cat.ID {
		return errors.New("já existe outra categoria com este nome")
	}

	return uc.repositorio.Update(cat)
}

func (uc *CategoriaDespesaUsuarioCasoDeUso) Deletar(id int) error {
	if _, err := uc.Buscar(id); err != nil {
		return err
	}
	return uc.repositorio.Delete(id)
}