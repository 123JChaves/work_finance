package entidade

import "time"

type Servico struct {
	ID          int       `json:"id" gorm:"primaryKey;autoIncrement"`
	Nome        string    `json:"nome" gorm:"type:varchar(255);not null"`
	CategoriaID int       `json:"categoria_id" gorm:"not null"`
	Categoria   Categoria `json:"categoria" gorm:"foreignKey:CategoriaID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"`
	Valor       float64   `json:"valor" gorm:"type:decimal(10,2);not null"`
	Data        time.Time `json:"data" gorm:"type:datetime;not null"`
	
	// ADICIONADO: Campo obrigatório que ancora o serviço à Empresa correspondente
	EmpresaID   int       `json:"empresa_id" gorm:"not null"` 
}

type ServicoRepositorio interface {
	Create(servico *Servico) error
	FindByID(id int) (*Servico, error)
	FindAll() ([]*Servico, error)
	Update(servico *Servico) error
	Delete(id int) error
}

func (Servico) TableName() string {
	return "servicos"
}