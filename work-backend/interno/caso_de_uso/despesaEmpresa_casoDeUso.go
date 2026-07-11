package caso_de_uso

import (
	"errors"
	"time"
	"work-backend/interno/entidade"
	"work-backend/interno/regras"
)

type DespesaEmpresaCasoDeUso struct {
	repositorio   entidade.DespesaEmpresaRepositorio
	categoriaRepo entidade.CategoriaDespesaEmpresaRepositorio
	empresaRepo   entidade.EmpresaRepositorio
}

func NovoDespesaEmpresaCasoDeUso(
	repo entidade.DespesaEmpresaRepositorio,
	catRepo entidade.CategoriaDespesaEmpresaRepositorio,
	empRepo entidade.EmpresaRepositorio,
) *DespesaEmpresaCasoDeUso {
	return &DespesaEmpresaCasoDeUso{
		repositorio:   repo,
		categoriaRepo: catRepo,
		empresaRepo:   empRepo,
	}
}

func (uc *DespesaEmpresaCasoDeUso) Cadastrar(despesa *entidade.DespesaEmpresa) error {
	if err := regras.ValidarDespesaEmpresa(despesa); err != nil {
		return err
	}

	// Garante a existência da Empresa vinculada
	empresa, err := uc.empresaRepo.FindByID(despesa.EmpresaID)
	if err != nil || empresa == nil {
		return errors.New("a empresa informada para esta despesa não existe")
	}

	// Garante a existência e acopla o objeto relacional da Categoria
	categoria, err := uc.categoriaRepo.FindByID(despesa.CategoriaDespesaEmpresaID)
	if err != nil || categoria == nil {
		return errors.New("a categoria informada para esta despesa não existe")
	}
	despesa.CategoriaDespesaEmpresa = *categoria

	return uc.repositorio.Create(despesa)
}

func (uc *DespesaEmpresaCasoDeUso) ListarPorEmpresa(empresaID int) ([]*entidade.DespesaEmpresa, error) {
	if empresaID <= 0 {
		return nil, errors.New("ID da empresa inválido")
	}
	return uc.repositorio.FindAllByEmpresaID(empresaID)
}

func (uc *DespesaEmpresaCasoDeUso) Buscar(id int) (*entidade.DespesaEmpresa, error) {
	if id <= 0 {
		return nil, errors.New("o ID da despesa deve ser um número válido")
	}
	return uc.repositorio.FindByID(id)
}

func (uc *DespesaEmpresaCasoDeUso) Atualizar(id int, novaDescricao string, novoValor float64, novaCatID int, novaData time.Time) error {
	despesa, err := uc.Buscar(id)
	if err != nil {
		return err
	}

	if novaDescricao != "" { despesa.Descricao = novaDescricao }
	if novoValor > 0 { despesa.Valor = novoValor }
	if novaCatID > 0 { despesa.CategoriaDespesaEmpresaID = novaCatID }
	if !novaData.IsZero() { despesa.Data = novaData }

	if err := regras.ValidarDespesaEmpresa(despesa); err != nil {
		return err
	}

	if novaCatID > 0 {
		categoria, err := uc.categoriaRepo.FindByID(despesa.CategoriaDespesaEmpresaID)
		if err != nil || categoria == nil {
			return errors.New("a nova categoria informada não existe")
		}
		despesa.CategoriaDespesaEmpresa = *categoria
	}

	return uc.repositorio.Update(despesa)
}

func (uc *DespesaEmpresaCasoDeUso) Deletar(id int) error {
	if _, err := uc.Buscar(id); err != nil {
		return err
	}
	return uc.repositorio.Delete(id)
}