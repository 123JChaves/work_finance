package entidade

import "time"

type DespesaEmpresa struct {
	ID                        int                     `json:"id" gorm:"primaryKey;autoIncrement"`
	Descricao                 string                  `json:"descricao" gorm:"type:varchar(255);not null"`
	Valor                     float64                 `json:"valor" gorm:"type:decimal(10,2);not null"`
	Data                      time.Time               `json:"data" gorm:"type:datetime;not null"`
	CategoriaDespesaEmpresaID int                     `json:"categoria_despesa_empresa_id" gorm:"not null"`
	CategoriaDespesaEmpresa   CategoriaDespesaEmpresa `json:"categoria_despesa_empresa" gorm:"foreignKey:CategoriaDespesaEmpresaID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"`
	EmpresaID                 int                     `json:"empresa_id" gorm:"not null"`
}

type DespesaEmpresaRepositorio interface {
	Create(despesa *DespesaEmpresa) error
	FindByID(id int) (*DespesaEmpresa, error)
	FindAllByEmpresaID(empresaID int) ([]*DespesaEmpresa, error)
	Update(despesa *DespesaEmpresa) error
	Delete(id int) error
}

func (DespesaEmpresa) TableName() string {
	return "despesas_empresa"
}