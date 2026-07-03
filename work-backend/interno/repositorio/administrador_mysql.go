package repositorio

import (
	"errors"
	"work-backend/interno/entidade"
	"gorm.io/gorm"
)

type GORMAdministradorRepositorio struct {
	db *gorm.DB
}

func NovoGORMAdministradorRepositorio(db *gorm.DB) *GORMAdministradorRepositorio {
	return &GORMAdministradorRepositorio{db: db}
}

func (r *GORMAdministradorRepositorio) Create(admin *entidade.Administrador) error {
	// O GORM faz o INSERT e já preenche o admin.ID automaticamente
	return r.db.Create(admin).Error
}

func (r *GORMAdministradorRepositorio) FindByID(id int) (*entidade.Administrador, error) {
	var admin entidade.Administrador
	err := r.db.First(&admin, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("Administrador não encontrado")
		}
		return nil, err
	}
	return &admin, nil
}

func (r *GORMAdministradorRepositorio) FindAll() ([]*entidade.Administrador, error) {
	var lista []*entidade.Administrador
	err := r.db.Find(&lista).Error
	return lista, err
}

func (r *GORMAdministradorRepositorio) FindByEmail(email string) (*entidade.Administrador, error) {
	var admin entidade.Administrador
	err := r.db.Where("email = ?", email).First(&admin).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &admin, nil
}

func (r *GORMAdministradorRepositorio) Update(admin *entidade.Administrador) error {
	// Atualiza todos os campos do modelo com base no ID
	resultado := r.db.Save(admin)
	if resultado.Error != nil {
		return resultado.Error
	}
	if resultado.RowsAffected == 0 {
		return errors.New("Administrador não encontrado para atualização")
	}
	return nil
}

func (r *GORMAdministradorRepositorio) Delete(id int) error {
	resultado := r.db.Delete(&entidade.Administrador{}, id)
	if resultado.Error != nil {
		return resultado.Error
	}
	if resultado.RowsAffected == 0 {
		return errors.New("Administrador não encontrado para exclusão")
	}
	return nil
}