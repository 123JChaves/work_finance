package repositorio

import (
	"errors"
	"work-backend/interno/entidade"
	"gorm.io/gorm"
)

var _ entidade.CategoriaDespesaEmpresaRepositorio = (*GORMCategoriaDespesaEmpresaRepo)(nil)

type GORMCategoriaDespesaEmpresaRepo struct {
	db *gorm.DB
}

func NovoGORMCategoriaDespesaEmpresaRepo(db *gorm.DB) *GORMCategoriaDespesaEmpresaRepo {
	return &GORMCategoriaDespesaEmpresaRepo{db: db}
}

func (r *GORMCategoriaDespesaEmpresaRepo) Create(cat *entidade.CategoriaDespesaEmpresa) error {
	return r.db.Create(cat).Error
}

func (r *GORMCategoriaDespesaEmpresaRepo) FindByID(id int) (*entidade.CategoriaDespesaEmpresa, error) {
	var cat entidade.CategoriaDespesaEmpresa
	err := r.db.Limit(1).Find(&cat, id).Error
	if err != nil {
		return nil, err
	}
	if cat.ID == 0 {
		return nil, errors.New("categoria de despesa corporativa não encontrada")
	}
	return &cat, nil
}

func (r *GORMCategoriaDespesaEmpresaRepo) FindAll() ([]*entidade.CategoriaDespesaEmpresa, error) {
	var lista []*entidade.CategoriaDespesaEmpresa
	err := r.db.Find(&lista).Error
	return lista, err
}

func (r *GORMCategoriaDespesaEmpresaRepo) FindByNome(nome string) (*entidade.CategoriaDespesaEmpresa, error) {
	var categorias []entidade.CategoriaDespesaEmpresa
	err := r.db.Where("nome = ?", nome).Limit(1).Find(&categorias).Error
	if err != nil {
		return nil, err
	}
	if len(categorias) == 0 {
		return nil, nil
	}
	return &categorias[0], nil
}

func (r *GORMCategoriaDespesaEmpresaRepo) Update(cat *entidade.CategoriaDespesaEmpresa) error {
	return r.db.Save(cat).Error
}

func (r *GORMCategoriaDespesaEmpresaRepo) Delete(id int) error {
	resultado := r.db.Delete(&entidade.CategoriaDespesaEmpresa{}, id)
	if resultado.Error != nil {
		return resultado.Error
	}
	if resultado.RowsAffected == 0 {
		return errors.New("categoria de despesa corporativa não encontrada para exclusão")
	}
	return nil
}