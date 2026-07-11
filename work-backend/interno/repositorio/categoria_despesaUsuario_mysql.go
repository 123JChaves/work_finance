package repositorio

import (
	"errors"
	"work-backend/interno/entidade"
	"gorm.io/gorm"
)

// SEGURANÇA DE COMPILAÇÃO: Esta linha garante em tempo de compilação que o repositório
// implementa 100% da interface definida na camada de domínio. Se faltar algo, o Go avisa aqui.
var _ entidade.CategoriaDespesaUsuarioRepositorio = (*GORMCategoriaDespesaUsuarioRepo)(nil)

type GORMCategoriaDespesaUsuarioRepo struct {
	db *gorm.DB
}

func NovoGORMCategoriaDespesaUsuarioRepo(db *gorm.DB) *GORMCategoriaDespesaUsuarioRepo {
	return &GORMCategoriaDespesaUsuarioRepo{db: db}
}

func (r *GORMCategoriaDespesaUsuarioRepo) Create(cat *entidade.CategoriaDespesaUsuario) error {
	return r.db.Create(cat).Error
}

func (r *GORMCategoriaDespesaUsuarioRepo) FindByID(id int) (*entidade.CategoriaDespesaUsuario, error) {
	var cat entidade.CategoriaDespesaUsuario
	
	// .Find não dispara logs invasivos de erro no console caso o registro não exista
	err := r.db.Limit(1).Find(&cat, id).Error
	if err != nil {
		return nil, err
	}
	
	// Como chaves primárias começam em 1, se voltar 0 significa que não localizou nada
	if cat.ID == 0 {
		return nil, errors.New("categoria de despesa pessoal não encontrada")
	}
	return &cat, nil
}

func (r *GORMCategoriaDespesaUsuarioRepo) FindAll() ([]*entidade.CategoriaDespesaUsuario, error) {
	var lista []*entidade.CategoriaDespesaUsuario
	err := r.db.Find(&lista).Error
	return lista, err
}

func (r *GORMCategoriaDespesaUsuarioRepo) FindByNome(nome string) (*entidade.CategoriaDespesaUsuario, error) {
	var categorias []entidade.CategoriaDespesaUsuario
	
	err := r.db.Where("nome = ?", nome).Limit(1).Find(&categorias).Error
	if err != nil {
		return nil, err
	}
	
	// PROTEÇÃO: Retorna explicitamente nil caso o slice venha vazio, evitando panics
	if len(categorias) == 0 {
		return nil, nil
	}
	return &categorias[0], nil
}

func (r *GORMCategoriaDespesaUsuarioRepo) Update(cat *entidade.CategoriaDespesaUsuario) error {
	// O método .Save do GORM atualiza todos os campos do modelo com base na chave primária (ID)
	return r.db.Save(cat).Error
}

func (r *GORMCategoriaDespesaUsuarioRepo) Delete(id int) error {
	resultado := r.db.Delete(&entidade.CategoriaDespesaUsuario{}, id)
	if resultado.Error != nil {
		return resultado.Error
	}
	
	// PROTEÇÃO ADICIONADA: Dispara erro explícito caso tente deletar um ID inexistente
	if resultado.RowsAffected == 0 {
		return errors.New("categoria de despesa pessoal não encontrada para exclusão")
	}
	return nil
}
