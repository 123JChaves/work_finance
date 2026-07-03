package entidade

import "time"

type Administrador struct {
	ID          int       `json:"id"`
	Nome        string    `json:"nome"`
	Email       string    `json:"email"`
	Senha       string    `json:"senha"`
	DataCriacao time.Time `json:"dataCriacao"`
	DataEdicao  time.Time `json:"dataEdicao"`
}


type AdministradorRepositorio interface {
	Create(admin *Administrador) error
	FindByID(id int) (*Administrador, error)
	FindAll() ([]*Administrador, error)
	FindByEmail(email string) (*Administrador, error)
	Update(admin *Administrador) error
	Delete(id int) error
}

type PasswordHasher interface {
	Hash(password string) (string, error)
	Compare(hash, password string) bool
}

// Força o GORM a usar o nome correto da tabela no plural em português
func (Administrador) TableName() string {
	return "administradores"
}