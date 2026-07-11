package entidade

import "time"

type Cliente struct {
	ID          int       `json:"id" gorm:"primaryKey;autoIncrement"`
	Nome        string    `json:"nome" gorm:"type:varchar(255);not null"`
	Email       string    `json:"email" gorm:"type:varchar(255);not null"`
	CpfCnpj     string    `json:"cpf_cnpj" gorm:"type:varchar(20);not null"`
	Contato     *string   `json:"contato" gorm:"type:varchar(50);default:null"` 
	DataCriacao time.Time `json:"data_criacao" gorm:"type:datetime;not null"`
	EmpresaID   int       `json:"empresa_id" gorm:"not null"`
}

type ClienteRepositorio interface {
	Create(cliente *Cliente) error
	FindByID(id int) (*Cliente, error)
	FindAllByEmpresaID(empresaID int) ([]*Cliente, error)
	Update(cliente *Cliente) error
	Delete(id int) error
}

func (Cliente) TableName() string {
	return "clientes"
}