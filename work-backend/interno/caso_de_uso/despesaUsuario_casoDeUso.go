package caso_de_uso

import (
	"errors"
	"time"
	"work-backend/interno/entidade"
	"work-backend/interno/regras"
)

type DespesaUsuarioCasoDeUso struct {
	repositorio   entidade.DespesaUsuarioRepositorio
	categoriaRepo entidade.CategoriaDespesaUsuarioRepositorio
	usuarioRepo   entidade.UsuarioRepositorio
}

func NovoDespesaUsuarioCasoDeUso(
	repo entidade.DespesaUsuarioRepositorio,
	catRepo entidade.CategoriaDespesaUsuarioRepositorio,
	userRepo entidade.UsuarioRepositorio,
) *DespesaUsuarioCasoDeUso {
	return &DespesaUsuarioCasoDeUso{
		repositorio:   repo,
		categoriaRepo: catRepo,
		usuarioRepo:   userRepo,
	}
}

func (uc *DespesaUsuarioCasoDeUso) Cadastrar(despesa *entidade.DespesaUsuario) error {
	if err := regras.ValidarDespesaUsuario(despesa); err != nil {
		return err
	}

	// Garante a existência do Usuário dono do registro
	usuario, err := uc.usuarioRepo.FindByID(despesa.UsuarioID)
	if err != nil || usuario == nil {
		return errors.New("o usuário informado para esta despesa não existe")
	}

	// Garante a existência e acopla a Categoria correspondente
	categoria, err := uc.categoriaRepo.FindByID(despesa.CategoriaDespesaUsuarioID)
	if err != nil || categoria == nil {
		return errors.New("a categoria informada para esta despesa não existe")
	}
	despesa.CategoriaDespesaUsuario = *categoria

	return uc.repositorio.Create(despesa)
}

func (uc *DespesaUsuarioCasoDeUso) ListarPorUsuario(usuarioID int) ([]*entidade.DespesaUsuario, error) {
	if usuarioID <= 0 {
		return nil, errors.New("ID do usuário inválido")
	}
	return uc.repositorio.FindAllByUsuarioID(usuarioID)
}

func (uc *DespesaUsuarioCasoDeUso) Buscar(id int) (*entidade.DespesaUsuario, error) {
	if id <= 0 {
		return nil, errors.New("o ID da despesa deve ser válido")
	}
	return uc.repositorio.FindByID(id)
}

func (uc *DespesaUsuarioCasoDeUso) Atualizar(id int, descricao string, valor float64, catID int, data time.Time) error {
	despesa, err := uc.Buscar(id)
	if err != nil {
		return err
	}

	if descricao != "" { despesa.Descricao = descricao }
	if valor > 0 { despesa.Valor = valor }
	if catID > 0 { despesa.CategoriaDespesaUsuarioID = catID }
	if !data.IsZero() { despesa.Data = data }

	if err := regras.ValidarDespesaUsuario(despesa); err != nil {
		return err
	}

	if catID > 0 {
		categoria, err := uc.categoriaRepo.FindByID(despesa.CategoriaDespesaUsuarioID)
		if err != nil || categoria == nil {
			return errors.New("a nova categoria informada não existe")
		}
		despesa.CategoriaDespesaUsuario = *categoria
	}

	return uc.repositorio.Update(despesa)
}

func (uc *DespesaUsuarioCasoDeUso) Deletar(id int) error {
	if _, err := uc.Buscar(id); err != nil {
		return err
	}
	return uc.repositorio.Delete(id)
}