package regras

import (
	"time"
	"work-backend/interno/entidade"
)

// CalcularTotalPatrimonioUsuario (Clean Code): Centraliza de forma declarativa a soma de toda a riqueza.
// Fórmula: Saldo Operacional Líquido + Bens Pessoais (PF) + Patrimônio total das Empresas societárias (PJ).
func CalcularTotalPatrimonioUsuario(u *entidade.Usuario) float64 {
	total := u.TotalLiquido
	
	// Adiciona a soma de todos os ativos tangíveis pessoais (PF)
	total += CalcularTotalPatrimonioPessoal(&u.Patrimonio)

	// Adiciona a soma de todos os ativos imobilizados de suas empresas (PJ)
	for _, empresa := range u.Empresas {
		total += CalcularTotalPatrimonioEmpresa(&empresa.Patrimonio)
	}

	return total
}

// CalcularTotalPatrimonioPessoal isola a agregação matemática de bens da pessoa física
func CalcularTotalPatrimonioPessoal(p *entidade.Patrimonio) float64 {
	var soma float64
	for _, obj := range p.ObjetosPatrimoniais {
		soma += obj.ValorPatrimonial
	}
	return soma
}

// CalcularTotalPatrimonioEmpresa isola a agregação matemática de imobilizados (CAPEX) da pessoa jurídica
func CalcularTotalPatrimonioEmpresa(p *entidade.PatrimonioEmpresa) float64 {
	var soma float64
	for _, obj := range p.ObjetosPatrimoniais {
		soma += obj.ValorPatrimonial
	}
	return soma
}

// CalcularLucroLiquidoEmpresa calcula o faturamento operacional líquido corrente da empresa (OPEX)
func CalcularLucroLiquidoEmpresa(e *entidade.Empresa) float64 {
	var faturamento float64
	var despesas float64

	for _, servico := range e.Servicos {
		faturamento += servico.Valor
	}
	for _, despesa := range e.Despesas {
		despesas += despesa.Valor
	}

	return faturamento - despesas
}

// CalcularDepreciacaoSimplesCorp computa a perda de valor de ativos corporativos ao longo do tempo
func CalcularDepreciacaoSimplesCorp(o *entidade.ObjetoPatrimonialCorp, taxaAnual float64) float64 {
	anos := time.Since(o.DataAquisicao).Hours() / 24 / 365
	if anos <= 0 {
		return o.ValorPatrimonial
	}
	valorDepreciado := o.ValorPatrimonial * (1 - (taxaAnual * anos))
	if valorDepreciado < 0 {
		return 0
	}
	return valorDepreciado
}

// CalcularValorizacaoImobiliaria computa o ganho de capital composto para ativos pessoais (PF)
func CalcularValorizacaoImobiliaria(o *entidade.ObjetoPatrimonialPess, taxaAnual float64) float64 {
	anos := time.Since(o.DataAquisicao).Hours() / 24 / 365
	if anos <= 0 {
		return o.ValorPatrimonial
	}
	
	valorAtualizado := o.ValorPatrimonial
	for i := 0; i < int(anos); i++ {
		valueAcumulado := valorAtualizado * taxaAnual
		valorAtualizado += valueAcumulado
	}
	return valorAtualizado
}

// CalcularFechamentoFluxoCaixa apura a liquidez corrente consolidada do indivíduo
func CalcularFechamentoFluxoCaixa(u *entidade.Usuario) float64 {
	var lucroEmpresas float64
	for _, empresa := range u.Empresas {
		var receitaPJ float64
		var despesaPJ float64
		
		for _, s := range empresa.Servicos { receitaPJ += s.Valor }
		for _, d := range empresa.Despesas { despesaPJ += d.Valor }
		
		lucroEmpresas += (receitaPJ - despesaPJ)
	}

	var despesasPessoais float64
	for _, d := range u.Despesas {
		despesasPessoais += d.Valor
	}

	return u.SalarioBase + lucroEmpresas - despesasPessoais
}