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
	return r.db.Create(admin).Error
}

func (r *GORMAdministradorRepositorio) FindByID(id int) (*entidade.Administrador, error) {
	var admin entidade.Administrador
	
	// .Find não dispara o log "record not found" se o ID não existir
	err := r.db.Limit(1).Find(&admin, id).Error
	if err != nil {
		return nil, err
	}
	
	// Como o ID no banco começa em 1, se voltar 0 significa que não encontrou nada
	if admin.ID == 0 {
		return nil, errors.New("Administrador não encontrado")
	}
	
	return &admin, nil
}

func (r *GORMAdministradorRepositorio) FindAll() ([]*entidade.Administrador, error) {
	var lista []*entidade.Administrador
	err := r.db.Find(&lista).Error
	return lista, err
}

func (r *GORMAdministradorRepositorio) FindByEmail(email string) (*entidade.Administrador, error) {
	var admins []entidade.Administrador
	
	// .Find preenche o slice e não gera erros ou logs caso não encontre registros
	err := r.db.Where("email = ?", email).Limit(1).Find(&admins).Error
	if err != nil {
		return nil, err
	}
	
	// Se o slice veio vazio, o e-mail não existe no banco. Retorna nil, nil
	if len(admins) == 0 {
		return nil, nil
	}
	
	// Retorna o ponteiro do único elemento encontrado
	return &admins[0], nil
}

func (r *GORMAdministradorRepositorio) Update(admin *entidade.Administrador) error {
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