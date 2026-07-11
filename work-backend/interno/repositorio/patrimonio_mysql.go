package repositorio

import (
	"errors"
	"work-backend/interno/entidade"
	"gorm.io/gorm"
)

var _ entidade.PatrimonioRepositorio = (*GORMPatrimonioRepo)(nil)

type GORMPatrimonioRepo struct {
	db *gorm.DB
}

func NovoGORMPatrimonioRepo(db *gorm.DB) *GORMPatrimonioRepo {
	return &GORMPatrimonioRepo{db: db}
}

func (r *GORMPatrimonioRepo) Create(pat *entidade.Patrimonio) error {
	return r.db.Create(pat).Error
}

func (r *GORMPatrimonioRepo) FindByID(id int) (*entidade.Patrimonio, error) {
	var pat entidade.Patrimonio
	err := r.db.Preload("ObjetosPatrimoniais.CategoriaObjetoPatrimonial").Limit(1).Find(&pat, id).Error
	if err != nil {
		return nil, err
	}
	if pat.ID == 0 {
		return nil, errors.New("registro de patrimônio não encontrado")
	}
	return &pat, nil
}

func (r *GORMPatrimonioRepo) FindByUsuarioID(usuarioID int) (*entidade.Patrimonio, error) {
	var pat entidade.Patrimonio
	// Preload aninhado via ponto (.) carrega o objeto físico e a categoria dele simultaneamente
	err := r.db.Preload("ObjetosPatrimoniais").Preload("ObjetosPatrimoniais.CategoriaObjetoPatrimonial").Where("usuario_id = ?", usuarioID).Limit(1).Find(&pat).Error
	if err != nil || pat.ID == 0 {
		return nil, err
	}
	return &pat, nil
}

func (r *GORMPatrimonioRepo) Update(pat *entidade.Patrimonio) error {
	return r.db.Save(pat).Error
}

func (r *GORMPatrimonioRepo) Delete(id int) error {
	resultado := r.db.Delete(&entidade.Patrimonio{}, id)
	if resultado.Error != nil {
		return resultado.Error
	}
	if resultado.RowsAffected == 0 {
		return errors.New("registro patrimonial não encontrado para exclusão")
	}
	return nil
}