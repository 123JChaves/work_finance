package entidade

import "time"

type ObjetoPatrimonialPess struct {
	ID                           int                        `json:"id" gorm:"primaryKey;autoIncrement"`
	Nome                         string                     `json:"nome" gorm:"type:varchar(255);not null"`
	ValorPatrimonial             float64                    `json:"valor_patrimonial" gorm:"type:decimal(10,2);not null"`
	DataAquisicao                time.Time                  `json:"data_aquisicao" gorm:"type:datetime;not null"`
	PatrimonioID                 int                        `json:"patrimonio_id" gorm:"not null"`
	CategoriaObjetoPatrimonialID int                        `json:"categoria_objeto_id" gorm:"not null"`
	CategoriaObjetoPatrimonial   CategoriaObjetoPatrimonial `json:"categoria_objeto" gorm:"foreignKey:CategoriaObjetoPatrimonialID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"`
}

type ObjetoPatrimonialPessRepositorio interface {
	Create(objeto *ObjetoPatrimonialPess) error
	FindByID(id int) (*ObjetoPatrimonialPess, error)
	FindAllByPatrimonioID(patrimonioID int) ([]*ObjetoPatrimonialPess, error)
	Update(objeto *ObjetoPatrimonialPess) error
	Delete(id int) error
}

func (ObjetoPatrimonialPess) TableName() string {
	return "objetos_patrimoniais_pess"
}