package entidade

type Usuario struct {
	ID           int              `json:"id" gorm:"primaryKey;autoIncrement"`
	Nome         string           `json:"nome" gorm:"type:varchar(255);not null"`
	Cpf          string           `json:"cpf" gorm:"type:varchar(20);unique;not null"`
	Email        string           `json:"email" gorm:"type:varchar(255);unique;not null"`
	Senha        string           `json:"senha" gorm:"type:varchar(255);not null"`
	SalarioBase  float64          `json:"salario_base" gorm:"type:decimal(10,2);default:0.00;not null"`
	Despesas     []DespesaUsuario `json:"despesas" gorm:"foreignKey:UsuarioID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"`
	Empresas     []Empresa        `json:"empresas" gorm:"many2many:usuarios_empresas;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"`
	Patrimonio   Patrimonio       `json:"patrimonio" gorm:"foreignKey:UsuarioID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"`
	TotalLiquido float64          `json:"total_liquido" gorm:"type:decimal(10,2);default:0.00;not null"`
}

type UsuarioRepositorio interface {
	Create(usuario *Usuario) error
	FindByID(id int) (*Usuario, error)
	BuscarCompletoPorID(id int) (*Usuario, error) // Fundamental para o fechamento patrimonial profundo
	FindAll() ([]*Usuario, error)
	FindByEmail(email string) (*Usuario, error)
	Update(usuario *Usuario) error
	UpdateTotalLiquido(id int, valor float64) error
	Delete(id int) error
}

func (Usuario) TableName() string {
	return "usuarios"
}