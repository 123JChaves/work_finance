package caso_de_uso

import (
	"errors"
	"strings"
	"time"
	"work-backend/interno/entidade"
	"work-backend/interno/validador"
)

type Administrador_casoDeUso struct {
	repositorio entidade.AdministradorRepositorio
	hasher      entidade.PasswordHasher
}

func NovoAdministradorCasoDeUso(repo entidade.AdministradorRepositorio, hasher entidade.PasswordHasher) *Administrador_casoDeUso {
	return &Administrador_casoDeUso{repositorio: repo, hasher: hasher}
}

func (uc *Administrador_casoDeUso) Cadastrar(admin *entidade.Administrador) error {
	nomeLimpo := strings.TrimSpace(admin.Nome)
	if nomeLimpo == "" || !validador.NomeSeguro(nomeLimpo) {
		return errors.New("O nome é obrigatório e não pode conter caracteres especiais")
	}
	admin.Nome = nomeLimpo

	if admin.Email == "" || !validador.EmailValido(admin.Email) { // CORRIGIDO: de EmailValid para EmailValido
		return errors.New("Formato de e-mail inválido")
	}

	if len(admin.Senha) < 6 {
		return errors.New("A senha deve ter no mínimo 6 caracteres")
	}

	existente, _ := uc.repositorio.FindByEmail(admin.Email)
	if existente != nil {
		return errors.New("Este email já está cadastrado no sistema")
	}

	senhaCriptografada, err := uc.hasher.Hash(admin.Senha)
	if err != nil {
		return errors.New("Erro ao processar a segurança da senha")
	}
	admin.Senha = senhaCriptografada
	admin.DataCriacao = time.Now()
	admin.DataEdicao = time.Now()

	return uc.repositorio.Create(admin)
}

func (uc *Administrador_casoDeUso) Listar() ([]*entidade.Administrador, error) {
	return uc.repositorio.FindAll()
}

func (uc *Administrador_casoDeUso) Buscar(id int) (*entidade.Administrador, error) {
	if id <= 0 {
		return nil, errors.New("O ID do administrador deve ser um número válido")
	}
	return uc.repositorio.FindByID(id)
}

func (uc *Administrador_casoDeUso) Atualizar(id int, nome, email, novaSenha string) error {
	admin, err := uc.Buscar(id)
	if err != nil {
		return err
	}

	if nome != "" {
		if !validador.NomeSeguro(nome) {
			return errors.New("O nome possui caracteres inválidos")
		}
		admin.Nome = nome
	}

	if email != "" && email != admin.Email {
		if !validador.EmailValido(email) {
			return errors.New("Formato de e-mail inválido")
		}
		existente, _ := uc.repositorio.FindByEmail(email)
		if existente != nil {
			return errors.New("Este email já está cadastrado em outra conta")
		}
		admin.Email = email
	}

	if novaSenha != "" {
		if len(novaSenha) < 6 {
			return errors.New("A nova senha deve ter no mínimo 6 caracteres")
		}
		senhaHash, err := uc.hasher.Hash(novaSenha)
		if err != nil {
			return errors.New("Erro ao processar segurança da nova senha")
		}
		admin.Senha = senhaHash
	}

	admin.DataEdicao = time.Now()
	return uc.repositorio.Update(admin)
}

func (uc *Administrador_casoDeUso) Deletar(id int) error {
	if _, err := uc.Buscar(id); err != nil {
		return err
	}
	return uc.repositorio.Delete(id)
}