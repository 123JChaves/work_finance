package repositorio

import (
	"errors"
	"work-backend/interno/entidade"
	"gorm.io/gorm"
)

var _ entidade.ClienteRepositorio = (*GORMClienteRepo)(nil)

type GORMClienteRepo struct {
	db *gorm.DB
}

func NovoGORMClienteRepo(db *gorm.DB) *GORMClienteRepo {
	return &GORMClienteRepo{db: db}
}

func (r *GORMClienteRepo) Create(cliente *entidade.Cliente) error {
	return r.db.Create(cliente).Error
}

func (r *GORMClienteRepo) FindByID(id int) (*entidade.Cliente, error) {
	var cliente entidade.Cliente
	err := r.db.Limit(1).Find(&cliente, id).Error
	if err != nil {
		return nil, err
	}
	if cliente.ID == 0 {
		return nil, errors.New("cliente não encontrado")
	}
	return &cliente, nil
}

func (r *GORMClienteRepo) FindAllByEmpresaID(empresaID int) ([]*entidade.Cliente, error) {
	var lista []*entidade.Cliente
	err := r.db.Where("empresa_id = ?", empresaID).Find(&lista).Error
	return lista, err
}

func (r *GORMClienteRepo) Update(cliente *entidade.Cliente) error {
	return r.db.Save(cliente).Error
}

func (r *GORMClienteRepo) Delete(id int) error {
	resultado := r.db.Delete(&entidade.Cliente{}, id)
	if resultado.Error != nil {
		return resultado.Error
	}
	if resultado.RowsAffected == 0 {
		return errors.New("cliente não encontrado para exclusão")
	}
	return nil
}