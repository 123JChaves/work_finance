package caso_de_uso

import (
	"errors"
	"time"
	"work-backend/interno/entidade"
	"work-backend/interno/regras"
)

type ObjetoPatrimonialPessCasoDeUso struct {
	repositorio           entidade.ObjetoPatrimonialPessRepositorio
	patrimonioRepo        entidade.PatrimonioRepositorio
	categoriaObjetoRepo   entidade.CategoriaObjetoPatrimonialRepositorio
}

func NovoObjetoPatrimonialPessCasoDeUso(
	repo entidade.ObjetoPatrimonialPessRepositorio,
	patRepo entidade.PatrimonioRepositorio,
	catObjRepo entidade.CategoriaObjetoPatrimonialRepositorio,
) *ObjetoPatrimonialPessCasoDeUso {
	return &ObjetoPatrimonialPessCasoDeUso{
		repositorio:           repo,
		patrimonioRepo:        patRepo,
		categoriaObjetoRepo:   catObjRepo,
	}
}

func (uc *ObjetoPatrimonialPessCasoDeUso) Cadastrar(objeto *entidade.ObjetoPatrimonialPess, usuarioID int) error {
	// Busca o inventário de patrimônio ativo do usuário
	patrimonio, err := uc.patrimonioRepo.FindByUsuarioID(usuarioID)
	if err != nil || patrimonio == nil {
		return errors.New("o patrimônio pessoal deste usuário não está inicializado")
	}
	objeto.PatrimonioID = patrimonio.ID

	if err := regras.ValidarObjetoPatrimonialPess(objeto); err != nil {
		return err
	}

	// Garante a integridade da categoria vinculada
	categoria, err := uc.categoriaObjetoRepo.FindByID(objeto.CategoriaObjetoPatrimonialID)
	if err != nil || categoria == nil {
		return errors.New("a categoria informada para este bem pessoal não existe")
	}
	objeto.CategoriaObjetoPatrimonial = *categoria

	return uc.repositorio.Create(objeto)
}

func (uc *ObjetoPatrimonialPessCasoDeUso) ListarPorPatrimonio(patrimonioID int) ([]*entidade.ObjetoPatrimonialPess, error) {
	if patrimonioID <= 0 {
		return nil, errors.New("ID do patrimônio inválido")
	}
	return uc.repositorio.FindAllByPatrimonioID(patrimonioID)
}

func (uc *ObjetoPatrimonialPessCasoDeUso) Buscar(id int) (*entidade.ObjetoPatrimonialPess, error) {
	if id <= 0 {
		return nil, errors.New("o ID do objeto pessoal deve ser válido")
	}
	return uc.repositorio.FindByID(id)
}

func (uc *ObjetoPatrimonialPessCasoDeUso) ObterProjecaoValorizacao(id int, taxaAnual float64) (float64, error) {
	objeto, err := uc.Buscar(id)
	if err != nil || objeto == nil {
		return 0, errors.New("bem patrimonial pessoal não localizado")
	}
	if taxaAnual < 0 {
		return 0, errors.New("a taxa de valorização anual projetada não pode ser um número negativo")
	}
	return regras.CalcularValorizacaoImobiliaria(objeto, taxaAnual), nil
}

func (uc *ObjetoPatrimonialPessCasoDeUso) Atualizar(id int, nome string, valor float64, data time.Time, catID int) error {
	objeto, err := uc.Buscar(id)
	if err != nil {
		return err
	}

	if nome != "" { objeto.Nome = nome }
	if valor > 0 { objeto.ValorPatrimonial = valor }
	if !data.IsZero() { objeto.DataAquisicao = data }
	if catID > 0 { objeto.CategoriaObjetoPatrimonialID = catID }

	if err := regras.ValidarObjetoPatrimonialPess(objeto); err != nil {
		return err
	}

	if catID > 0 {
		categoria, err := uc.categoriaObjetoRepo.FindByID(objeto.CategoriaObjetoPatrimonialID)
		if err != nil || categoria == nil {
			return errors.New("a nova categoria patrimonial informada não existe")
		}
		objeto.CategoriaObjetoPatrimonial = *categoria
	}

	return uc.repositorio.Update(objeto)
}

func (uc *ObjetoPatrimonialPessCasoDeUso) Deletar(id int) error {
	if _, err := uc.Buscar(id); err != nil {
		return err
	}
	return uc.repositorio.Delete(id)
}