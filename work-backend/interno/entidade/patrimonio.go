package entidade

type Patrimonio struct {
	ID                  int                     `json:"id" gorm:"primaryKey;autoIncrement"`
	UsuarioID           int                     `json:"usuario_id" gorm:"not null;unique"`
	ObjetosPatrimoniais []ObjetoPatrimonialPess `json:"objetos_patrimoniais" gorm:"foreignKey:PatrimonioID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"`
}

type PatrimonioRepositorio interface {
	Create(pat *Patrimonio) error
	FindByID(id int) (*Patrimonio, error)
	FindByUsuarioID(usuarioID int) (*Patrimonio, error)
	Update(pat *Patrimonio) error
	Delete(id int) error
}

func (Patrimonio) TableName() string {
	return "patrimonios"
}