package entidade

type CategoriaDespesaEmpresa struct {
	ID   int    `json:"id" gorm:"primaryKey;autoIncrement"`
	Nome string `json:"nome" gorm:"type:varchar(255);not null"`
}

type CategoriaDespesaEmpresaRepositorio interface {
	Create(cat *CategoriaDespesaEmpresa) error
	FindByID(id int) (*CategoriaDespesaEmpresa, error)
	FindAll() ([]*CategoriaDespesaEmpresa, error)
	FindByNome(nome string) (*CategoriaDespesaEmpresa, error)
	Update(cat *CategoriaDespesaEmpresa) error
	Delete(id int) error
}

func (CategoriaDespesaEmpresa) TableName() string {
	return "categorias_despesas_empresa"
}