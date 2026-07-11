package repositorio

import (
	"errors"
	"work-backend/interno/entidade"
	"gorm.io/gorm"
)

// SEGURANÇA DE COMPILAÇÃO: Garante que o repositório implementa 100% da interface do domínio
var _ entidade.CategoriaRepositorio = (*GORMCategoriaRepositorio)(nil)

type GORMCategoriaRepositorio struct {
	db *gorm.DB
}

func NovoGORMCategoriaRepositorio(db *gorm.DB) *GORMCategoriaRepositorio {
	return &GORMCategoriaRepositorio{db: db}
}

func (r *GORMCategoriaRepositorio) Create(categ *entidade.Categoria) error {
	return r.db.Create(categ).Error
}

func (r *GORMCategoriaRepositorio) FindByID(id int) (*entidade.Categoria, error) {
	var categ entidade.Categoria
	err := r.db.Limit(1).Find(&categ, id).Error
	if err != nil {
		return nil, err
	}
	
	if categ.ID == 0 {
		return nil, errors.New("categoria não encontrada")
	}
	return &categ, nil
}

func (r *GORMCategoriaRepositorio) FindAll() ([]*entidade.Categoria, error) {
	var lista []*entidade.Categoria
	err := r.db.Find(&lista).Error
	return lista, err
}

func (r *GORMCategoriaRepositorio) FindByNome(nome string) (*entidade.Categoria, error) {
	var categorias []entidade.Categoria
	err := r.db.Where("nome = ?", nome).Limit(1).Find(&categorias).Error
	if err != nil {
		return nil, err
	}
	
	// CORRIGIDO: Retorna explicitamente nil se o slice vier vazio (evita panics de index 0)
	if len(categorias) == 0 {
		return nil, nil
	}
	return &categorias[0], nil
}

func (r *GORMCategoriaRepositorio) Update(categ *entidade.Categoria) error {
	return r.db.Save(categ).Error
}

func (r *GORMCategoriaRepositorio) Delete(id int) error {
	resultado := r.db.Delete(&entidade.Categoria{}, id)
	if resultado.Error != nil {
		return resultado.Error
	}
	
	// PROTEÇÃO ADICIONADA: Retorna erro explícito caso o ID não exista no banco
	if resultado.RowsAffected == 0 {
		return errors.New("categoria não encontrada para exclusão")
	}
	return nil
}