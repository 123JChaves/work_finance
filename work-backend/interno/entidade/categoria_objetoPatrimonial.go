package entidade

type CategoriaObjetoPatrimonial struct {
	ID   int    `json:"id" gorm:"primaryKey;autoIncrement"`
	Nome string `json:"nome" gorm:"type:varchar(255);not null"`
}

type CategoriaObjetoPatrimonialRepositorio interface {
	Create(cat *CategoriaObjetoPatrimonial) error
	FindByID(id int) (*CategoriaObjetoPatrimonial, error)
	FindAll() ([]*CategoriaObjetoPatrimonial, error)
	FindByNome(nome string) (*CategoriaObjetoPatrimonial, error) // Requisito de consistência
	Update(cat *CategoriaObjetoPatrimonial) error
	Delete(id int) error
}

func (CategoriaObjetoPatrimonial) TableName() string {
	return "categorias_objetos_patrimoniais"
}