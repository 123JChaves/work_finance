package caso_de_uso

import (
	"errors"
	"strings"
	"work-backend/interno/entidade"
	"work-backend/interno/regras"
)

type EmpresaCasoDeUso struct {
	repositorio entidade.EmpresaRepositorio
}

func NovoEmpresaCasoDeUso(repo entidade.EmpresaRepositorio) *EmpresaCasoDeUso {
	return &EmpresaCasoDeUso{repositorio: repo}
}

func (uc *EmpresaCasoDeUso) Cadastrar(empresa *entidade.Empresa) error {
	if err := regras.ValidarEmpresa(empresa); err != nil {
		return err
	}

	// Remove pontos e traços antes de persistir no MySQL (Boa prática)
	empresa.Cnpj = strings.NewReplacer(".", "", "-", "", "/", "").Replace(empresa.Cnpj)

	existente, _ := uc.repositorio.FindByCnpj(empresa.Cnpj)
	if existente != nil {
		return errors.New("já existe uma empresa cadastrada com este CNPJ no sistema")
	}

	return uc.repositorio.Create(empresa)
}

func (uc *EmpresaCasoDeUso) Listar() ([]*entidade.Empresa, error) {
	return uc.repositorio.FindAll()
}

func (uc *EmpresaCasoDeUso) Buscar(id int) (*entidade.Empresa, error) {
	if id <= 0 {
		return nil, errors.New("o ID da empresa deve ser um número válido")
	}
	return uc.repositorio.FindByID(id)
}

func (uc *EmpresaCasoDeUso) ObterDRE(id int) (float64, error) {
	// Carrega a empresa com todo o histórico de faturamento e despesas preenchidos
	empresa, err := uc.repositorio.BuscarCompletoPorID(id)
	if err != nil || empresa == nil {
		return 0, errors.New("não foi possível encontrar os dados financeiros desta empresa")
	}

	// Executa a função contábil isolada da pasta de regras
	return regras.CalcularLucroLiquidoEmpresa(empresa), nil
}

func (uc *EmpresaCasoDeUso) Atualizar(id int, novoNome, novoCnpj string) error {
	empresa, err := uc.Buscar(id)
	if err != nil {
		return err
	}

	if novoNome != "" { empresa.Nome = novoNome }
	if novoCnpj != "" { empresa.Cnpj = novoCnpj }

	if err := regras.ValidarEmpresa(empresa); err != nil {
		return err
	}

	return uc.repositorio.Update(empresa)
}

func (uc *EmpresaCasoDeUso) Deletar(id int) error {
	if _, err := uc.Buscar(id); err != nil {
		return err
	}
	return uc.repositorio.Delete(id)
}