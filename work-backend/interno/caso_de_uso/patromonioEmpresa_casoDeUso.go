package caso_de_uso

import (
	"errors"
	"work-backend/interno/entidade"
	"work-backend/interno/regras"
)

type PatrimonioEmpresaCasoDeUso struct {
	repositorio   entidade.PatrimonioEmpresaRepositorio
	empresaRepo   entidade.EmpresaRepositorio
}

func NovoPatrimonioEmpresaCasoDeUso(repo entidade.PatrimonioEmpresaRepositorio, empRepo entidade.EmpresaRepositorio) *PatrimonioEmpresaCasoDeUso {
	return &PatrimonioEmpresaCasoDeUso{
		repositorio:   repo,
		empresaRepo:   empRepo,
	}
}

func (uc *PatrimonioEmpresaCasoDeUso) Inicializar(empresaID int) error {
	if empresaID <= 0 {
		return errors.New("ID de empresa inválido")
	}

	// Garante a existência da empresa correspondente
	empresa, err := uc.empresaRepo.FindByID(empresaID)
	if err != nil || empresa == nil {
		return errors.New("a empresa informada não existe no sistema")
	}

	// Impede a duplicação
	existente, _ := uc.repositorio.FindByEmpresaID(empresaID)
	if existente != nil {
		return errors.New("esta empresa já possui um inventário de patrimônio inicializado")
	}

	patrimonio := &entidade.PatrimonioEmpresa{EmpresaID: empresaID}
	return uc.repositorio.Create(patrimonio)
}

func (uc *PatrimonioEmpresaCasoDeUso) ObterValorConsolidado(empresaID int) (float64, error) {
	patrimonio, err := uc.repositorio.FindByEmpresaID(empresaID)
	if err != nil || patrimonio == nil {
		return 0, errors.New("patrimônio corporativo não localizado para esta empresa")
	}

	return regras.CalcularTotalPatrimonioEmpresa(patrimonio), nil
}

func (uc *PatrimonioEmpresaCasoDeUso) Deletar(id int) error {
	if id <= 0 {
		return errors.New("ID inválido")
	}
	pat, err := uc.repositorio.FindByID(id)
	if err != nil || pat == nil {
		return errors.New("inventário patrimonial corporativo não localizado")
	}
	return uc.repositorio.Delete(id)
}