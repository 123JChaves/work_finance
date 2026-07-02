package caso_de_uso

import (
	"errors"
	"time"
	"work-backend/interno/entidade"
)

type Administrador_casoDeUso struct {
	repositorio entidade.AdministradorRepositorio
	hasher      entidade.PasswordHasher
}

func NovoAdministradorCasoDeUso(
	repositorio entidade.AdministradorRepositorio, 
	hasher entidade.PasswordHasher) *Administrador_casoDeUso {
		return &Administrador_casoDeUso{
			repositorio: repositorio,
			hasher:      hasher,
	}
}

func (uc *Administrador_casoDeUso) Cadastrar(administrador *entidade.Administrador) error {
	if administrador.Email == "" {
		return errors.New("O email é obrigatório")
	}
	if len(administrador.Senha) < 6 {
		return errors.New("A senha deve ter no mínimo 6 caracteres")
	}

	existente, _ := uc.repositorio.FindByEmail(administrador.Email)
	if existente != nil {
		return errors.New("Este email já está cadastrado no sistema")
	}

	senhaCriptografada, erro := uc.hasher.Hash(administrador.Senha)
	if erro != nil {
		return errors.New("Erro ao processar a segurança da senha") // Corrigido: errors com um 'r'
	}
	administrador.Senha = senhaCriptografada

	administrador.DataCriacao = time.Now()
	administrador.DataEdicao = time.Now()

	return uc.repositorio.Create(administrador)
}

func (uc *Administrador_casoDeUso) Listar() ([]*entidade.Administrador, error) {
	// Acessa o repositório e retorna a lista de administradores
	return uc.repositorio.FindAll()
}

