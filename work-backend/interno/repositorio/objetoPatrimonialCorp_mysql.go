package repositorio

import (
	"errors"
	"work-backend/interno/entidade"
	"gorm.io/gorm"
)

var _ entidade.ObjetoPatrimonialCorpRepositorio = (*GORMObjetoPatrimonialCorpRepo)(nil)

type GORMObjetoPatrimonialCorpRepo struct {
	db *gorm.DB
}

func NovoGORMObjetoPatrimonialCorpRepo(db *gorm.DB) *GORMObjetoPatrimonialCorpRepo {
	return &GORMObjetoPatrimonialCorpRepo{db: db}
}

func (r *GORMObjetoPatrimonialCorpRepo) Create(objeto *entidade.ObjetoPatrimonialCorp) error {
	return r.db.Create(objeto).Error
}

func (r *GORMObjetoPatrimonialCorpRepo) FindByID(id int) (*entidade.ObjetoPatrimonialCorp, error) {
	var objeto entidade.ObjetoPatrimonialCorp
	err := r.db.Limit(1).Find(&objeto, id).Error
	if err != nil {
		return nil, err
	}
	if objeto.ID == 0 {
		return nil, errors.New("objeto patrimonial corporativo não encontrado")
	}
	return &objeto, nil
}

func (r *GORMObjetoPatrimonialCorpRepo) FindAllByEmpresaID(empresaID int) ([]*entidade.ObjetoPatrimonialCorp, error) {
	var lista []*entidade.ObjetoPatrimonialCorp
	err := r.db.Where("empresa_id = ?", empresaID).Find(&lista).Error
	return lista, err
}

func (r *GORMObjetoPatrimonialCorpRepo) Update(objeto *entidade.ObjetoPatrimonialCorp) error {
	return r.db.Save(objeto).Error
}

func (r *GORMObjetoPatrimonialCorpRepo) Delete(id int) error {

	resultado := r.db.Delete(&entidade.ObjetoPatrimonialCorp{}, id)
	if resultado.Error != nil {
		return resultado.Error
	}
	if resultado.RowsAffected == 0 {
		return errors.New("objeto patrimonial corporativo não encontrado para exclusão")
	}
	return nil
}