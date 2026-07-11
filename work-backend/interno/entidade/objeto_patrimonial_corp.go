package entidade

import "time"

type ObjetoPatrimonialCorp struct {
	ID                  int       `json:"id" gorm:"primaryKey;autoIncrement"`
	Nome                string    `json:"nome" gorm:"type:varchar(255);not null"`
	ValorPatrimonial    float64   `json:"valor_patrimonial" gorm:"type:decimal(10,2);not null"`
	DataAquisicao       time.Time `json:"data_aquisicao" gorm:"type:datetime;not null"`
	PatrimonioEmpresaID int       `json:"patrimonio_empresa_id" gorm:"not null"`
	EmpresaID           int       `json:"empresa_id" gorm:"not null"`
}

type ObjetoPatrimonialCorpRepositorio interface {
	Create(objeto *ObjetoPatrimonialCorp) error
	FindByID(id int) (*ObjetoPatrimonialCorp, error)
	FindAllByEmpresaID(empresaID int) ([]*ObjetoPatrimonialCorp, error)
	Update(objeto *ObjetoPatrimonialCorp) error
	Delete(id int) error
}

func (ObjetoPatrimonialCorp) TableName() string {
	return "objetos_patrimoniais_corp"
}