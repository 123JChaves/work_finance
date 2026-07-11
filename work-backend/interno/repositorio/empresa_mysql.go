package repositorio

import (
	"errors"
	"work-backend/interno/entidade"
	"gorm.io/gorm"
)

var _ entidade.EmpresaRepositorio = (*GORMEmpresaRepo)(nil)

type GORMEmpresaRepo struct {
	db *gorm.DB
}

func NovoGORMEmpresaRepo(db *gorm.DB) *GORMEmpresaRepo {
	return &GORMEmpresaRepo{db: db}
}

func (r *GORMEmpresaRepo) Create(empresa *entidade.Empresa) error {
	return r.db.Create(empresa).Error
}

func (r *GORMEmpresaRepo) FindByID(id int) (*entidade.Empresa, error) {
	var empresa entidade.Empresa
	err := r.db.Limit(1).Find(&empresa, id).Error
	if err != nil {
		return nil, err
	}
	if empresa.ID == 0 {
		return nil, errors.New("empresa não encontrada")
	}
	return &empresa, nil
}

// BuscarCompletoPorID (Clean Arch): Resolve a árvore de dependências sem expor o GORM para o Caso de Uso
func (r *GORMEmpresaRepo) BuscarCompletoPorID(id int) (*entidade.Empresa, error) {
	var empresa entidade.Empresa
	err := r.db.
		Preload("Servicos").
		Preload("Despesas").
		Preload("Patrimonio").
		Preload("Patrimonio.ObjetosPatrimoniais").
		Limit(1).Find(&empresa, id).Error
		
	if err != nil {
		return nil, err
	}
	return &empresa, nil
}

func (r *GORMEmpresaRepo) FindAll() ([]*entidade.Empresa, error) {
	var lista []*entidade.Empresa
	err := r.db.Find(&lista).Error
	return lista, err
}

func (r *GORMEmpresaRepo) FindByCnpj(cnpj string) (*entidade.Empresa, error) {
	var empresas []entidade.Empresa
	err := r.db.Where("cnpj = ?", cnpj).Limit(1).Find(&empresas).Error
	if err != nil || len(empresas) == 0 {
		return nil, err
	}
	return &empresas[0], nil
}

func (r *GORMEmpresaRepo) Update(empresa *entidade.Empresa) error {
	return r.db.Save(empresa).Error
}

func (r *GORMEmpresaRepo) Delete(id int) error {
	resultado := r.db.Delete(&entidade.Empresa{}, id)
	if resultado.Error != nil {
		return resultado.Error
	}
	if resultado.RowsAffected == 0 {
		return errors.New("empresa não encontrada para exclusão")
	}
	return nil
}