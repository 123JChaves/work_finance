package repositorio

import (
	"errors"
	"work-backend/interno/entidade"
	"gorm.io/gorm"
)

var _ entidade.ServicoRepositorio = (*GORMServicoRepositorio)(nil)

type GORMServicoRepositorio struct {
	db *gorm.DB
}

func NovoGORMServicoRepositorio(db *gorm.DB) *GORMServicoRepositorio {
	return &GORMServicoRepositorio{db: db}
}

func (r *GORMServicoRepositorio) Create(servico *entidade.Servico) error {
	return r.db.Create(servico).Error
}

func (r *GORMServicoRepositorio) FindByID(id int) (*entidade.Servico, error) {
	var servico entidade.Servico
	err := r.db.Preload("Categoria").Limit(1).Find(&servico, id).Error
	if err != nil {
		return nil, err
	}
	if servico.ID == 0 {
		return nil, errors.New("Serviço não encontrado")
	}
	return &servico, nil
}

func (r *GORMServicoRepositorio) FindAll() ([]*entidade.Servico, error) {
	var lista []*entidade.Servico
	err := r.db.Preload("Categoria").Find(&lista).Error
	return lista, err
}

func (r *GORMServicoRepositorio) Update(servico *entidade.Servico) error {
	resultado := r.db.Save(servico)
	if resultado.Error != nil {
		return resultado.Error
	}
	if resultado.RowsAffected == 0 {
		return errors.New("Serviço não encontrado para atualização")
	}
	return nil
}

func (r *GORMServicoRepositorio) Delete(id int) error {
	resultado := r.db.Delete(&entidade.Servico{}, id)
	if resultado.Error != nil {
		return resultado.Error
	}
	if resultado.RowsAffected == 0 {
		return errors.New("Serviço não encontrado para exclusão")
	}
	return nil
}

// ADICIONADO: Expõe o DB de forma limpa para os casos de uso calcularem relatórios sem violar a interface de domínio
func (r *GORMServicoRepositorio) DB() *gorm.DB {
	return r.db
}