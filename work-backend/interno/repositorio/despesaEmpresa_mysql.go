package repositorio

import (
	"errors"
	"work-backend/interno/entidade"
	"gorm.io/gorm"
)

var _ entidade.DespesaEmpresaRepositorio = (*GORMDespesaEmpresaRepo)(nil)

type GORMDespesaEmpresaRepo struct {
	db *gorm.DB
}

func NovoGORMDespesaEmpresaRepo(db *gorm.DB) *GORMDespesaEmpresaRepo {
	return &GORMDespesaEmpresaRepo{db: db}
}

func (r *GORMDespesaEmpresaRepo) Create(despesa *entidade.DespesaEmpresa) error {
	return r.db.Create(despesa).Error
}

func (r *GORMDespesaEmpresaRepo) FindByID(id int) (*entidade.DespesaEmpresa, error) {
	var despesa entidade.DespesaEmpresa
	err := r.db.Preload("CategoriaDespesaEmpresa").Limit(1).Find(&despesa, id).Error
	if err != nil {
		return nil, err
	}
	if despesa.ID == 0 {
		return nil, errors.New("despesa corporativa não encontrada")
	}
	return &despesa, nil
}

func (r *GORMDespesaEmpresaRepo) FindAllByEmpresaID(empresaID int) ([]*entidade.DespesaEmpresa, error) {
	var lista []*entidade.DespesaEmpresa
	err := r.db.Preload("CategoriaDespesaEmpresa").Where("empresa_id = ?", empresaID).Find(&lista).Error
	return lista, err
}

func (r *GORMDespesaEmpresaRepo) Update(despesa *entidade.DespesaEmpresa) error {
	return r.db.Save(despesa).Error
}

func (r *GORMDespesaEmpresaRepo) Delete(id int) error {
	resultado := r.db.Delete(&entidade.DespesaEmpresa{}, id)
	if resultado.Error != nil {
		return resultado.Error
	}
	if resultado.RowsAffected == 0 {
		return errors.New("despesa corporativa não encontrada para exclusão")
	}
	return nil
}