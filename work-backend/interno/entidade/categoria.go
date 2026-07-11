package entidade

type Categoria struct {
	ID   int    `json:"id" gorm:"primaryKey;autoIncrement"`
	Nome string `json:"nome" gorm:"type:varchar(255);not null"`
}

type CategoriaRepositorio interface {
	Create(categ *Categoria) error
	FindByID(id int) (*Categoria, error)
	FindAll() ([]*Categoria, error)
	FindByNome(nome string) (*Categoria, error)
	Update(categ *Categoria) error
	Delete(id int) error
}

func (Categoria) TableName() string {
	return "categorias"
}