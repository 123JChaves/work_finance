package caso_de_uso

import (
	"errors"
	"time"
	"work-backend/interno/entidade"
	"work-backend/interno/regras"
)

type ObjetoPatrimonialCorpCasoDeUso struct {
	repositorio           entidade.ObjetoPatrimonialCorpRepositorio
	empresaRepo           entidade.EmpresaRepositorio
	patrimonioEmpresaRepo entidade.PatrimonioEmpresaRepositorio
}

func NovoObjetoPatrimonialCorpCasoDeUso(
	repo entidade.ObjetoPatrimonialCorpRepositorio,
	empRepo entidade.EmpresaRepositorio,
	patRepo entidade.PatrimonioEmpresaRepositorio,
) *ObjetoPatrimonialCorpCasoDeUso {
	return &ObjetoPatrimonialCorpCasoDeUso{
		repositorio:           repo,
		empresaRepo:           empRepo,
		patrimonioEmpresaRepo: patRepo,
	}
}

func (uc *ObjetoPatrimonialCorpCasoDeUso) Cadastrar(objeto *entidade.ObjetoPatrimonialCorp) error {
	if err := regras.ValidarObjetoPatrimonialCorp(objeto); err != nil {
		return err
	}

	// Garante que a empresa informada existe
	empresa, err := uc.empresaRepo.FindByID(objeto.EmpresaID)
	if err != nil || empresa == nil {
		return errors.New("a empresa informada para este ativo corporativo não existe")
	}

	// Garante que o inventário de patrimônio da empresa existe
	patrimonio, err := uc.patrimonioEmpresaRepo.FindByEmpresaID(objeto.EmpresaID)
	if err != nil || patrimonio == nil {
		return errors.New("o patrimônio corporativo da empresa não está inicializado")
	}
	objeto.PatrimonioEmpresaID = patrimonio.ID

	return uc.repositorio.Create(objeto)
}

func (uc *ObjetoPatrimonialCorpCasoDeUso) ListarPorEmpresa(empresaID int) ([]*entidade.ObjetoPatrimonialCorp, error) {
	if empresaID <= 0 {
		return nil, errors.New("ID da empresa inválido")
	}
	return uc.repositorio.FindAllByEmpresaID(empresaID)
}

func (uc *ObjetoPatrimonialCorpCasoDeUso) Buscar(id int) (*entidade.ObjetoPatrimonialCorp, error) {
	if id <= 0 {
		return nil, errors.New("o ID do objeto patrimonial deve ser um número válido")
	}
	return uc.repositorio.FindByID(id)
}

func (uc *ObjetoPatrimonialCorpCasoDeUso) ObterValorDepreciado(id int, taxaAnual float64) (float64, error) {
	objeto, err := uc.Buscar(id)
	if err != nil || objeto == nil {
		return 0, errors.New("objeto patrimonial corporativo não localizado")
	}
	if taxaAnual < 0 || taxaAnual > 1 {
		return 0, errors.New("a taxa de depreciação deve ser um valor decimal entre 0.00 e 1.00")
	}
	return regras.CalcularDepreciacaoSimplesCorp(objeto, taxaAnual), nil
}

func (uc *ObjetoPatrimonialCorpCasoDeUso) Atualizar(id int, nome string, valor float64, data time.Time) error {
	objeto, err := uc.Buscar(id)
	if err != nil {
		return err
	}

	if nome != "" { objeto.Nome = nome }
	if valor > 0 { objeto.ValorPatrimonial = valor }
	if !data.IsZero() { objeto.DataAquisicao = data }

	if err := regras.ValidarObjetoPatrimonialCorp(objeto); err != nil {
		return err
	}

	return uc.repositorio.Update(objeto)
}

func (uc *ObjetoPatrimonialCorpCasoDeUso) Deletar(id int) error {
	if _, err := uc.Buscar(id); err != nil {
		return err
	}
	return uc.repositorio.Delete(id)
}