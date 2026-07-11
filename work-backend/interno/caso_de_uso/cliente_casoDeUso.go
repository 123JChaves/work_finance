package caso_de_uso

import (
	"errors"
	"time"
	"work-backend/interno/entidade"
	"work-backend/interno/regras"
)

type ClienteCasoDeUso struct {
	repositorio entidade.ClienteRepositorio
	empresaRepo entidade.EmpresaRepositorio
}

func NovoClienteCasoDeUso(repo entidade.ClienteRepositorio, empRepo entidade.EmpresaRepositorio) *ClienteCasoDeUso {
	return &ClienteCasoDeUso{
		repositorio: repo,
		empresaRepo: empRepo,
	}
}

func (uc *ClienteCasoDeUso) Cadastrar(cliente *entidade.Cliente) error {
	if err := regras.ValidarCliente(cliente); err != nil {
		return err
	}

	// Garante a existência da empresa controladora do cliente
	empresa, err := uc.empresaRepo.FindByID(cliente.EmpresaID)
	if err != nil || empresa == nil {
		return errors.New("a empresa informada para este cliente não existe no sistema")
	}

	cliente.DataCriacao = time.Now()
	return uc.repositorio.Create(cliente)
}

func (uc *ClienteCasoDeUso) ListarPorEmpresa(empresaID int) ([]*entidade.Cliente, error) {
	if empresaID <= 0 {
		return nil, errors.New("ID da empresa inválido")
	}
	return uc.repositorio.FindAllByEmpresaID(empresaID)
}

func (uc *ClienteCasoDeUso) Buscar(id int) (*entidade.Cliente, error) {
	if id <= 0 {
		return nil, errors.New("o ID do cliente deve ser um número válido")
	}
	return uc.repositorio.FindByID(id)
}

func (uc *ClienteCasoDeUso) Atualizar(id int, nome, email, cpfCnpj string, contato *string) error {
	cliente, err := uc.Buscar(id)
	if err != nil {
		return err
	}

	if nome != "" { cliente.Nome = nome }
	if email != "" { cliente.Email = email }
	if cpfCnpj != "" { cliente.CpfCnpj = cpfCnpj }
	if contato != nil { cliente.Contato = contato }

	if err := regras.ValidarCliente(cliente); err != nil {
		return err
	}

	return uc.repositorio.Update(cliente)
}

func (uc *ClienteCasoDeUso) Deletar(id int) error {
	if _, err := uc.Buscar(id); err != nil {
		return err
	}
	return uc.repositorio.Delete(id)
}