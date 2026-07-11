package entidade

type PatrimonioEmpresa struct {
	ID                  int                     `json:"id" gorm:"primaryKey;autoIncrement"`
	EmpresaID           int                     `json:"empresa_id" gorm:"not null;unique"`
	ObjetosPatrimoniais []ObjetoPatrimonialCorp `json:"objetos_patrimoniais" gorm:"foreignKey:PatrimonioEmpresaID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"`
}

type PatrimonioEmpresaRepositorio interface {
	Create(pat *PatrimonioEmpresa) error
	FindByID(id int) (*PatrimonioEmpresa, error) // Adicionado para consistência arquitetural
	FindByEmpresaID(empresaID int) (*PatrimonioEmpresa, error)
	Update(pat *PatrimonioEmpresa) error
	Delete(id int) error
}

func (PatrimonioEmpresa) TableName() string {
	return "patrimonios_empresa"
}