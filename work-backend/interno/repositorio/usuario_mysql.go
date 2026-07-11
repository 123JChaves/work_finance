package repositorio

import (
	"errors"
	"work-backend/interno/entidade"
	"gorm.io/gorm"
)

// SEGURANÇA DE COMPILAÇÃO: Garante o contrato da interface de domínio sem furos
var _ entidade.UsuarioRepositorio = (*GORMUsuarioRepositorio)(nil)

type GORMUsuarioRepositorio struct {
	db *gorm.DB
}

func NovoGORMUsuarioRepository(db *gorm.DB) *GORMUsuarioRepositorio {
	return &GORMUsuarioRepositorio{db: db}
}

func (r *GORMUsuarioRepositorio) Create(usuario *entidade.Usuario) error {
	return r.db.Create(usuario).Error
}

func (r *GORMUsuarioRepositorio) FindByID(id int) (*entidade.Usuario, error) {
	var usuario entidade.Usuario
	err := r.db.Limit(1).Find(&usuario, id).Error
	if err != nil {
		return nil, err
	}
	if usuario.ID == 0 {
		return nil, errors.New("usuário não encontrado")
	}
	return &usuario, nil
}

// BuscarCompletoPorID realiza os Preloads de forma profunda para abastecer o Caso de Uso Financeiro
func (r *GORMUsuarioRepositorio) BuscarCompletoPorID(id int) (*entidade.Usuario, error) {
	var usuario entidade.Usuario
	err := r.db.
		Preload("Despesas").
		Preload("Patrimonio").
		Preload("Patrimonio.ObjetosPatrimoniais").
		Preload("Empresas").
		Preload("Empresas.Servicos").
		Preload("Empresas.Despesas").
		Preload("Empresas.Patrimonio").
		Preload("Empresas.Patrimonio.ObjetosPatrimoniais").
		Limit(1).Find(&usuario, id).Error

	if err != nil {
		return nil, err
	}
	if usuario.ID == 0 {
		return nil, errors.New("usuário não encontrado")
	}
	return &usuario, nil
}

func (r *GORMUsuarioRepositorio) FindAll() ([]*entidade.Usuario, error) {
	var lista []*entidade.Usuario
	err := r.db.Find(&lista).Error
	return lista, err
}

func (r *GORMUsuarioRepositorio) FindByEmail(email string) (*entidade.Usuario, error) {
	var usuario entidade.Usuario
	
	// CORRIGIDO: Busca direto em uma struct única para evitar panics de índice e ponteiros inválidos
	err := r.db.Where("email = ?", email).Limit(1).Find(&usuario).Error
	if err != nil {
		return nil, err
	}
	if usuario.ID == 0 {
		return nil, nil // Retorna nil, nil para o caso de uso saber que o e-mail está livre para cadastro
	}
	return &usuario, nil
}

func (r *GORMUsuarioRepositorio) Update(usuario *entidade.Usuario) error {
	return r.db.Save(usuario).Error
}

func (r *GORMUsuarioRepositorio) UpdateTotalLiquido(id int, valor float64) error {
	return r.db.Model(&entidade.Usuario{}).Where("id = ?", id).Update("total_liquido", valor).Error
}

func (r *GORMUsuarioRepositorio) Delete(id int) error {
	resultado := r.db.Delete(&entidade.Usuario{}, id)
	if resultado.Error != nil {
		return resultado.Error
	}
	if resultado.RowsAffected == 0 {
		return errors.New("usuário não localizado para exclusão")
	}
	return nil
}