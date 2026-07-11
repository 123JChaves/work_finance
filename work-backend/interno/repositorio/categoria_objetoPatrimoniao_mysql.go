package repositorio

import (
	"errors"
	"work-backend/interno/entidade"
	"gorm.io/gorm"
)

var _ entidade.CategoriaObjetoPatrimonialRepositorio = (*GORMCategoriaObjetoPatrimonialRepo)(nil)

type GORMCategoriaObjetoPatrimonialRepo struct {
	db *gorm.DB
}

func NovoGORMCategoriaObjetoPatrimonialRepo(db *gorm.DB) *GORMCategoriaObjetoPatrimonialRepo {
	return &GORMCategoriaObjetoPatrimonialRepo{db: db}
}

func (r *GORMCategoriaObjetoPatrimonialRepo) Create(cat *entidade.CategoriaObjetoPatrimonial) error {
	return r.db.Create(cat).Error
}

func (r *GORMCategoriaObjetoPatrimonialRepo) FindByID(id int) (*entidade.CategoriaObjetoPatrimonial, error) {
	var cat entidade.CategoriaObjetoPatrimonial
	err := r.db.Limit(1).Find(&cat, id).Error
	if err != nil {
		return nil, err
	}
	if cat.ID == 0 {
		return nil, errors.New("categoria patrimonial não encontrada")
	}
	return &cat, nil
}

func (r *GORMCategoriaObjetoPatrimonialRepo) FindAll() ([]*entidade.CategoriaObjetoPatrimonial, error) {
	var lista []*entidade.CategoriaObjetoPatrimonial
	err := r.db.Find(&lista).Error
	return lista, err
}

func (r *GORMCategoriaObjetoPatrimonialRepo) FindByNome(nome string) (*entidade.CategoriaObjetoPatrimonial, error) {
	var categorias []entidade.CategoriaObjetoPatrimonial
	err := r.db.Where("nome = ?", nome).Limit(1).Find(&categorias).Error
	if err != nil {
		return nil, err
	}
	if len(categorias) == 0 {
		return nil, nil
	}
	return &categorias[0], nil
}

func (r *GORMCategoriaObjetoPatrimonialRepo) Update(cat *entidade.CategoriaObjetoPatrimonial) error {
	return r.db.Save(cat).Error
}

func (r *GORMCategoriaObjetoPatrimonialRepo) Delete(id int) error {
	resultado := r.db.Delete(&entidade.CategoriaObjetoPatrimonial{}, id)
	if resultado.Error != nil {
		return resultado.Error
	}
	if resultado.RowsAffected == 0 {
		return errors.New("categoria patrimonial não encontrada para exclusão")
	}
	return nil
}