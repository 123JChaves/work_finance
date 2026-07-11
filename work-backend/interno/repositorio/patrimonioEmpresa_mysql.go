package repositorio

import (
	"errors"
	"work-backend/interno/entidade"
	"gorm.io/gorm"
)

var _ entidade.PatrimonioEmpresaRepositorio = (*GORMPatrimonioEmpresaRepo)(nil)

type GORMPatrimonioEmpresaRepo struct {
	db *gorm.DB
}

func NovoGORMPatrimonioEmpresaRepo(db *gorm.DB) *GORMPatrimonioEmpresaRepo {
	return &GORMPatrimonioEmpresaRepo{db: db}
}

func (r *GORMPatrimonioEmpresaRepo) Create(pat *entidade.PatrimonioEmpresa) error {
	return r.db.Create(pat).Error
}

func (r *GORMPatrimonioEmpresaRepo) FindByID(id int) (*entidade.PatrimonioEmpresa, error) {
	var pat entidade.PatrimonioEmpresa
	err := r.db.Preload("ObjetosPatrimoniais").Limit(1).Find(&pat, id).Error
	if err != nil {
		return nil, err
	}
	if pat.ID == 0 {
		return nil, errors.New("registro de patrimônio corporativo não encontrado")
	}
	return &pat, nil
}

func (r *GORMPatrimonioEmpresaRepo) FindByEmpresaID(empresaID int) (*entidade.PatrimonioEmpresa, error) {
	var pat entidade.PatrimonioEmpresa
	err := r.db.Preload("ObjetosPatrimoniais").Where("empresa_id = ?", empresaID).Limit(1).Find(&pat).Error
	if err != nil || pat.ID == 0 {
		return nil, err
	}
	return &pat, nil
}

func (r *GORMPatrimonioEmpresaRepo) Update(pat *entidade.PatrimonioEmpresa) error {
	return r.db.Save(pat).Error
}

func (r *GORMPatrimonioEmpresaRepo) Delete(id int) error {
	resultado := r.db.Delete(&entidade.PatrimonioEmpresa{}, id)
	if resultado.Error != nil {
		return resultado.Error
	}
	if resultado.RowsAffected == 0 {
		return errors.New("registro patrimonial corporativo não encontrado para exclusão")
	}
	return nil
}