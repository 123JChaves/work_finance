package repositorio

import (
	"errors"
	"work-backend/interno/entidade"
	"gorm.io/gorm"
)

var _ entidade.DespesaUsuarioRepositorio = (*GORMDespesaUsuarioRepo)(nil)

type GORMDespesaUsuarioRepo struct {
	db *gorm.DB
}

func NovoGORMDespesaUsuarioRepo(db *gorm.DB) *GORMDespesaUsuarioRepo {
	return &GORMDespesaUsuarioRepo{db: db}
}

func (r *GORMDespesaUsuarioRepo) Create(despesa *entidade.DespesaUsuario) error {
	return r.db.Create(despesa).Error
}

func (r *GORMDespesaUsuarioRepo) FindByID(id int) (*entidade.DespesaUsuario, error) {
	var despesa entidade.DespesaUsuario
	err := r.db.Preload("CategoriaDespesaUsuario").Limit(1).Find(&despesa, id).Error
	if err != nil {
		return nil, err
	}
	if despesa.ID == 0 {
		return nil, errors.New("despesa pessoal não encontrada")
	}
	return &despesa, nil
}

func (r *GORMDespesaUsuarioRepo) FindAllByUsuarioID(usuarioID int) ([]*entidade.DespesaUsuario, error) {
	var lista []*entidade.DespesaUsuario
	err := r.db.Preload("CategoriaDespesaUsuario").Where("usuario_id = ?", usuarioID).Find(&lista).Error
	return lista, err
}

func (r *GORMDespesaUsuarioRepo) Update(despesa *entidade.DespesaUsuario) error {
	return r.db.Save(despesa).Error
}

func (r *GORMDespesaUsuarioRepo) Delete(id int) error {
	resultado := r.db.Delete(&entidade.DespesaUsuario{}, id)
	if resultado.Error != nil {
		return resultado.Error
	}
	if resultado.RowsAffected == 0 {
		return errors.New("despesa pessoal não encontrada para exclusão")
	}
	return nil
}