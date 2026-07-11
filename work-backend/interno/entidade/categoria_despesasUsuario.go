package entidade

type CategoriaDespesaUsuario struct {
	ID   int    `json:"id" gorm:"primaryKey;autoIncrement"`
	Nome string `json:"nome" gorm:"type:varchar(255);not null"`
}

type CategoriaDespesaUsuarioRepositorio interface {
	Create(cat *CategoriaDespesaUsuario) error
	FindByID(id int) (*CategoriaDespesaUsuario, error)
	FindAll() ([]*CategoriaDespesaUsuario, error)
	FindByNome(nome string) (*CategoriaDespesaUsuario, error)
	Update(cat *CategoriaDespesaUsuario) error
	Delete(id int) error
}

func (CategoriaDespesaUsuario) TableName() string {
	return "categorias_despesas_usuario"
}