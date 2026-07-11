package repositorio

import (
	"errors"
	"work-backend/interno/entidade"
	"gorm.io/gorm"
)

var _ entidade.ObjetoPatrimonialPessRepositorio = (*GORMObjetoPatrimonialPessRepo)(nil)

type GORMObjetoPatrimonialPessRepo struct {
	db *gorm.DB
}

func NovoGORMObjetoPatrimonialPessRepo(db *gorm.DB) *GORMObjetoPatrimonialPessRepo {
	return &GORMObjetoPatrimonialPessRepo{db: db}
}

func (r *GORMObjetoPatrimonialPessRepo) Create(objeto *entidade.ObjetoPatrimonialPess) error {
	return r.db.Create(objeto).Error
}

func (r *GORMObjetoPatrimonialPessRepo) FindByID(id int) (*entidade.ObjetoPatrimonialPess, error) {
	var objeto entidade.ObjetoPatrimonialPess
	err := r.db.Preload("CategoriaObjetoPatrimonial").Limit(1).Find(&objeto, id).Error
	if err != nil {
		return nil, err
	}
	if objeto.ID == 0 {
		return nil, errors.New("objeto patrimonial pessoal não encontrado")
	}
	return &objeto, nil
}

func (r *GORMObjetoPatrimonialPessRepo) FindAllByPatrimonioID(patrimonioID int) ([]*entidade.ObjetoPatrimonialPess, error) {
	var lista []*entidade.ObjetoPatrimonialPess
	err := r.db.Preload("CategoriaObjetoPatrimonial").Where("patrimonio_id = ?", patrimonioID).Find(&lista).Error
	return lista, err
}

func (r *GORMObjetoPatrimonialPessRepo) Update(objeto *entidade.ObjetoPatrimonialPess) error {
	return r.db.Save(objeto).Error
}

func (r *GORMObjetoPatrimonialPessRepo) Delete(id int) error {
	resultado := r.db.Delete(&entidade.ObjetoPatrimonialPess{}, id)
	if resultado.Error != nil {
		return resultado.Error
	}
	if resultado.RowsAffected == 0 {
		return errors.New("objeto patrimonial pessoal não encontrado para exclusão")
	}
	return nil
}