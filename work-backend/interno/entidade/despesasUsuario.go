package entidade

import "time"

type DespesaUsuario struct {
	ID                        int                     `json:"id" gorm:"primaryKey;autoIncrement"`
	Descricao                 string                  `json:"descricao" gorm:"type:varchar(255);not null"`
	Valor                     float64                 `json:"valor" gorm:"type:decimal(10,2);not null"`
	CategoriaDespesaUsuarioID int                     `json:"categoria_despesa_usuario_id" gorm:"not null"`
	CategoriaDespesaUsuario   CategoriaDespesaUsuario `json:"categoria_despesa_usuario" gorm:"foreignKey:CategoriaDespesaUsuarioID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"`
	Data                      time.Time               `json:"data" gorm:"type:datetime;not null"`
	UsuarioID                 int                     `json:"usuario_id" gorm:"not null"`
}

type DespesaUsuarioRepositorio interface {
	Create(despesa *DespesaUsuario) error
	FindByID(id int) (*DespesaUsuario, error)
	FindAllByUsuarioID(usuarioID int) ([]*DespesaUsuario, error)
	Update(despesa *DespesaUsuario) error
	Delete(id int) error
}

func (DespesaUsuario) TableName() string {
	return "despesas_usuarios"
}