package entidade

type Empresa struct {
	ID         int               `json:"id" gorm:"primaryKey;autoIncrement"`
	Nome       string            `json:"nome" gorm:"type:varchar(255);not null"`
	Cnpj       string            `json:"cnpj" gorm:"type:varchar(20);unique;not null"`
	Usuarios   []Usuario         `json:"usuarios" gorm:"many2many:usuarios_empresas;"`
	Servicos   []Servico         `json:"servicos" gorm:"foreignKey:EmpresaID;"`
	Despesas   []DespesaEmpresa  `json:"despesas" gorm:"foreignKey:EmpresaID;"`
	Clientes   []Cliente         `json:"clientes" gorm:"foreignKey:EmpresaID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"`
	Patrimonio PatrimonioEmpresa `json:"patrimonio" gorm:"foreignKey:EmpresaID;"`
}

type EmpresaRepositorio interface {
	Create(empresa *Empresa) error
	FindByID(id int) (*Empresa, error)
	BuscarCompletoPorID(id int) (*Empresa, error) // Novo: Indispensável para o Caso de Uso Financeiro
	FindAll() ([]*Empresa, error)
	FindByCnpj(cnpj string) (*Empresa, error)
	Update(empresa *Empresa) error
	Delete(id int) error
}

func (Empresa) TableName() string {
	return "empresas"
}