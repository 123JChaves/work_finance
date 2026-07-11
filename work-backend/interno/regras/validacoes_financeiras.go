package regras

import (
	"errors"
	"time"
	"work-backend/interno/entidade"
	"work-backend/interno/validador"
)

// ValidarServico isola todas as travas regulatórias de faturamento do serviço
func ValidarServico(s *entidade.Servico) error {
	if !validador.NomeSeguro(s.Nome) {
		return errors.New("o nome do serviço possui caracteres inválidos")
	}
	if s.Valor <= 0 {
		return errors.New("o valor do serviço deve ser maior do que zero")
	}
	if s.CategoriaID <= 0 {
		return errors.New("o ID da categoria vinculada é inválido")
	}
	if s.Data.IsZero() {
		return errors.New("a data em que o serviço foi prestado é obrigatória")
	}
	if s.Data.After(time.Now()) {
		return errors.New("a data da prestação do serviço não pode ser uma data futura")
	}
	return nil
}

// ValidarCliente isola o teste de consistência cadastral de clientes da empresa
func ValidarCliente(c *entidade.Cliente) error {
	if !validador.NomeSeguro(c.Nome) {
		return errors.New("o nome do cliente possui caracteres inválidos")
	}
	if !validador.EmailValido(c.Email) {
		return errors.New("o formato do e-mail do cliente é inválido")
	}
	if !validador.DocumentoPossuiTamanhoValido(c.CpfCnpj) {
		return errors.New("o documento CPF/CNPJ deve possuir 11 ou 14 dígitos numéricos")
	}
	if c.EmpresaID <= 0 {
		return errors.New("o cliente precisa estar associado a uma empresa válida")
	}
	return nil
}

// ValidarCategoriaDespesaEmpresa garante que a categoria PJ cumpre os requisitos textuais
func ValidarCategoriaDespesaEmpresa(c *entidade.CategoriaDespesaEmpresa) error {
	if !validador.NomeSeguro(c.Nome) {
		return errors.New("o nome da categoria possui caracteres inválidos")
	}
	if !validador.ComecaComMaiuscula(c.Nome) {
		return errors.New("o nome da categoria corporativa deve começar com uma letra maiúscula")
	}
	return nil
}

// ValidarCategoriaDespesaUsuario garante que a categoria PF cumpre os requisitos textuais
func ValidarCategoriaDespesaUsuario(c *entidade.CategoriaDespesaUsuario) error {
	if !validador.NomeSeguro(c.Nome) {
		return errors.New("o nome da categoria possui caracteres inválidos")
	}
	if !validador.ComecaComMaiuscula(c.Nome) {
		return errors.New("o nome da categoria deve obrigatoriamente começar com uma letra maiúscula")
	}
	return nil
}

// ValidarCategoriaObjetoPatrimonial valida categorias de ativos tangíveis
func ValidarCategoriaObjetoPatrimonial(c *entidade.CategoriaObjetoPatrimonial) error {
	if !validador.NomeSeguro(c.Nome) {
		return errors.New("o nome da categoria patrimonial possui caracteres inválidos")
	}
	if !validador.ComecaComMaiuscula(c.Nome) {
		return errors.New("o nome da categoria patrimonial deve começar com uma letra maiúscula")
	}
	return nil
}

// ValidarDespesaEmpresa centraliza as restrições fiscais das saídas de caixa da PJ
func ValidarDespesaEmpresa(d *entidade.DespesaEmpresa) error {
	if !validador.NomeSeguro(d.Descricao) {
		return errors.New("a descrição da despesa corporativa possui caracteres inválidos")
	}
	if d.Valor <= 0 {
		return errors.New("o valor da despesa deve ser estritamente maior do que zero")
	}
	if d.CategoriaDespesaEmpresaID <= 0 {
		return errors.New("o ID da categoria corporativa informado é inválido")
	}
	if d.EmpresaID <= 0 {
		return errors.New("a despesa deve estar associada a uma empresa válida")
	}
	if d.Data.IsZero() {
		return errors.New("a data de lançamento da despesa é obrigatória")
	}
	if d.Data.After(time.Now()) {
		return errors.New("não é permitido lançar despesas corporativas em datas futuras")
	}
	return nil
}

// ValidarDespesaUsuario centraliza as restrições de gastos do indivíduo (PF)
func ValidarDespesaUsuario(d *entidade.DespesaUsuario) error {
	if !validador.NomeSeguro(d.Descricao) {
		return errors.New("a descrição da despesa pessoal possui caracteres inválidos")
	}
	if d.Valor <= 0 {
		return errors.New("o valor da despesa deve ser maior do que zero")
	}
	if d.CategoriaDespesaUsuarioID <= 0 {
		return errors.New("o ID da categoria informado é inválido")
	}
	if d.UsuarioID <= 0 {
		return errors.New("a despesa deve pertencer a um usuário válido")
	}
	if d.Data.IsZero() {
		return errors.New("a data da despesa pessoal é obrigatória")
	}
	if d.Data.After(time.Now()) {
		return errors.New("não é permitido registrar despesas em datas futuras")
	}
	return nil
}

// ValidarEmpresa valida as regras cadastrais de uma nova Pessoa Jurídica
func ValidarEmpresa(e *entidade.Empresa) error {
	if !validador.NomeSeguro(e.Nome) {
		return errors.New("o nome da empresa possui caracteres inválidos ou inseguros")
	}
	if !validador.DocumentoPossuiTamanhoValido(e.Cnpj) || !validador.EhPessoaJuridica(e.Cnpj) {
		return errors.New("o CNPJ corporativo fornecido é inválido. Deve possuir 14 dígitos numéricos")
	}
	return nil
}

// ValidarObjetoPatrimonialCorp centraliza as travas cadastrais do bem de capital da PJ
func ValidarObjetoPatrimonialCorp(o *entidade.ObjetoPatrimonialCorp) error {
	if !validador.NomeSeguro(o.Nome) {
		return errors.New("o nome do objeto patrimonial possui caracteres inválidos")
	}
	if o.ValorPatrimonial <= 0 {
		return errors.New("o valor patrimonial do bem deve ser maior do que zero")
	}
	if o.PatrimonioEmpresaID <= 0 {
		return errors.New("o objeto deve estar vinculado a um patrimônio corporativo válido")
	}
	if o.EmpresaID <= 0 {
		return errors.New("o objeto deve estar associado a uma empresa válida")
	}
	if o.DataAquisicao.IsZero() {
		return errors.New("a data de aquisição do bem é obrigatória")
	}
	if o.DataAquisicao.After(time.Now()) {
		return errors.New("a data de aquisição do patrimônio não pode ser uma data futura")
	}
	return nil
}

// ValidarObjetoPatrimonialPess centraliza as restrições cadastrais do bem particular da PF
func ValidarObjetoPatrimonialPess(o *entidade.ObjetoPatrimonialPess) error {
	if !validador.NomeSeguro(o.Nome) {
		return errors.New("o nome do objeto patrimonial pessoal possui caracteres inválidos")
	}
	if o.ValorPatrimonial <= 0 {
		return errors.New("o valor patrimonial do bem deve ser estritamente maior do que zero")
	}
	if o.PatrimonioID <= 0 {
		return errors.New("o objeto deve estar acoplado a um patrimônio pessoal válido")
	}
	if o.CategoriaObjetoPatrimonialID <= 0 {
		return errors.New("o ID da categoria do objeto informado é inválido")
	}
	if o.DataAquisicao.IsZero() {
		return errors.New("a data de aquisição do bem é obrigatória")
	}
	if o.DataAquisicao.After(time.Now()) {
		return errors.New("a data de aquisição do patrimônio pessoal não pode ser uma data futura")
	}
	return nil
}

func ValidarUsuario(u *entidade.Usuario) error {
	if !validador.NomeSeguro(u.Nome) {
		return errors.New("o nome do usuário possui caracteres inválidos")
	}
	if !validador.EmailValido(u.Email) {
		return errors.New("o formato do e-mail do usuário é inválido")
	}
	if !validador.DocumentoPossuiTamanhoValido(u.Cpf) || validador.EhPessoaJuridica(u.Cpf) {
		return errors.New("o documento CPF fornecido é inválido. Deve possuir 11 dígitos numéricos")
	}
	if u.SalarioBase < 0 {
		return errors.New("o salário base do usuário não pode ser um valor negativo")
	}
	return nil
}

func ValidarCategoria(c *entidade.Categoria) error {
	if !validador.NomeSeguro(c.Nome) {
		return errors.New("o nome da categoria possui caracteres inválidos")
	}
	if !validador.ComecaComMaiuscula(c.Nome) {
		return errors.New("o nome da categoria deve começar com uma letra maiúscula")
	}
	return nil
}
