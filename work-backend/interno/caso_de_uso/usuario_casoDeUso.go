package caso_de_uso

import (
	"errors"
	"strings"
	"work-backend/interno/entidade"
	"work-backend/interno/regras"
)

type UsuarioCasoDeUso struct {
	repositorio entidade.UsuarioRepositorio
	hasher      entidade.PasswordHasher
}

func NovoUsuarioCasoDeUso(repo entidade.UsuarioRepositorio, hasher entidade.PasswordHasher) *UsuarioCasoDeUso {
	return &UsuarioCasoDeUso{repositorio: repo, hasher: hasher}
}

type ResumoPatrimonialDTO struct {
	UsuarioNome            string  `json:"usuario_nome"`
	TotalLiquidoDisponivel float64 `json:"total_liquido_disponivel"`
	TotalBensPessoais      float64 `json:"total_bens_pessoais"`
	TotalPatrimonioPJ      float64 `json:"total_patrimonio_empresas"`
	RiquezaTotalAcumulada  float64 `json:"riqueza_total_acumulada"`
}

func (uc *UsuarioCasoDeUso) Cadastrar(usuario *entidade.Usuario) error {
	if err := regras.ValidarUsuario(usuario); err != nil {
		return err
	}

	usuario.Cpf = strings.NewReplacer(".", "", "-", "").Replace(usuario.Cpf)
	
	if len(usuario.Senha) < 6 {
		return errors.New("a senha de acesso deve ter no mínimo 6 caracteres")
	}

	existente, _ := uc.repositorio.FindByEmail(usuario.Email)
	if existente != nil {
		return errors.New("este endereço de e-mail já está em uso por outro usuário")
	}

	senhaHash, err := uc.hasher.Hash(usuario.Senha)
	if err != nil {
		return errors.New("falha interna ao processar segurança da senha")
	}
	usuario.Senha = senhaHash

	return uc.repositorio.Create(usuario)
}

func (uc *UsuarioCasoDeUso) Listar() ([]*entidade.Usuario, error) {
	usuarios, err := uc.repositorio.FindAll()
	if err != nil {
		return nil, err
	}

	// Oculta as hashes de segurança das senhas na resposta da listagem
	for _, u := range usuarios {
		u.Senha = "[PROTEGIDO]"
	}

	return usuarios, nil
}

func (uc *UsuarioCasoDeUso) Buscar(id int) (*entidade.Usuario, error) {
	if id <= 0 {
		return nil, errors.New("o ID deve ser um número inteiro positivo")
	}
	usuario, err := uc.repositorio.FindByID(id)
	if err != nil {
		return nil, err
	}
	return usuario, nil
}

func (uc *UsuarioCasoDeUso) Atualizar(id int, nome, email, cpf string, salarioBase float64) error {
	usuario, err := uc.Buscar(id)
	if err != nil {
		return err
	}

	if nome != "" { 
		usuario.Nome = nome 
	}
	if email != "" && email != usuario.Email {
		existente, _ := uc.repositorio.FindByEmail(email)
		if existente != nil {
			return errors.New("este endereço de e-mail já está cadastrado em outra conta")
		}
		usuario.Email = email
	}
	if cpf != "" {
		usuario.Cpf = strings.NewReplacer(".", "", "-", "").Replace(cpf)
	}
	if salarioBase >= 0 {
		usuario.SalarioBase = salarioBase
	}

	// Valida os novos dados consolidados usando as regras de negócio
	if err := regras.ValidarUsuario(usuario); err != nil {
		return err
	}

	return uc.repositorio.Update(usuario)
}

func (uc *UsuarioCasoDeUso) Deletar(id int) error {
	if _, err := uc.Buscar(id); err != nil {
		return err
	}
	return uc.repositorio.Delete(id)
}

func (uc *UsuarioCasoDeUso) ProcessarFechamentoMensal(id int) (float64, error) {
	usuario, err := uc.repositorio.BuscarCompletoPorID(id)
	if err != nil || usuario == nil {
		return 0, errors.New("dados do usuário não encontrados")
	}

	novaLiquidez := regras.CalcularFechamentoFluxoCaixa(usuario)
	
	err = uc.repositorio.UpdateTotalLiquido(usuario.ID, novaLiquidez)
	if err != nil {
		return 0, errors.New("erro ao persistir consolidação de saldo em conta")
	}

	return novaLiquidez, nil
}

func (uc *UsuarioCasoDeUso) ObterBalançoPatrimonial(id int) (*ResumoPatrimonialDTO, error) {
	usuario, err := uc.repositorio.BuscarCompletoPorID(id)
	if err != nil || usuario == nil {
		return nil, errors.New("usuário não localizado")
	}

	saldoCorrente := regras.CalcularFechamentoFluxoCaixa(usuario)
	usuario.TotalLiquido = saldoCorrente

	var totalPJ float64
	for _, emp := range usuario.Empresas {
		totalPJ += regras.CalcularTotalPatrimonioEmpresa(&emp.Patrimonio)
	}

	return &ResumoPatrimonialDTO{
		UsuarioNome:            usuario.Nome,
		TotalLiquidoDisponivel: usuario.TotalLiquido,
		TotalBensPessoais:      regras.CalcularTotalPatrimonioPessoal(&usuario.Patrimonio),
		TotalPatrimonioPJ:      totalPJ,
		RiquezaTotalAcumulada:  regras.CalcularTotalPatrimonioUsuario(usuario),
	}, nil
}