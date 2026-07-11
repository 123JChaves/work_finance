package caso_de_uso

import (
	"errors"
	"work-backend/interno/entidade"
	"work-backend/interno/regras"
)

type PatrimonioCasoDeUso struct {
	repositorio entidade.PatrimonioRepositorio
	usuarioRepo entidade.UsuarioRepositorio
}

func NovoPatrimonioCasoDeUso(repo entidade.PatrimonioRepositorio, userRepo entidade.UsuarioRepositorio) *PatrimonioCasoDeUso {
	return &PatrimonioCasoDeUso{
		repositorio: repo,
		usuarioRepo: userRepo,
	}
}

func (uc *PatrimonioCasoDeUso) Inicializar(usuarioID int) error {
	if usuarioID <= 0 {
		return errors.New("ID de usuário inválido")
	}

	// Garante a existência do usuário dono
	usuario, err := uc.usuarioRepo.FindByID(usuarioID)
	if err != nil || usuario == nil {
		return errors.New("o usuário informado não existe")
	}

	// Impede a duplicação (1 para 1 exclusivo)
	existente, _ := uc.repositorio.FindByUsuarioID(usuarioID)
	if existente != nil {
		return errors.New("este usuário já possui uma carteira de patrimônio inicializada")
	}

	patrimonio := &entidade.Patrimonio{UsuarioID: usuarioID}
	return uc.repositorio.Create(patrimonio)
}

func (uc *PatrimonioCasoDeUso) ObterValorConsolidado(usuarioID int) (float64, error) {
	patrimonio, err := uc.repositorio.FindByUsuarioID(usuarioID)
	if err != nil || patrimonio == nil {
		return 0, errors.New("patrimônio pessoal não localizado para este usuário")
	}

	return regras.CalcularTotalPatrimonioPessoal(patrimonio), nil
}

func (uc *PatrimonioCasoDeUso) Deletar(id int) error {
	var pat *entidade.Patrimonio
	var err error
	if id <= 0 {
		return errors.New("ID inválido")
	}
	pat, err = uc.repositorio.FindByID(id)
	if err != nil || pat == nil {
		return errors.New("inventário patrimonial não localizado")
	}
	return uc.repositorio.Delete(id)
}