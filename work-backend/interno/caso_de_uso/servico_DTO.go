package caso_de_uso

type TotalPorCategoriaDTO struct {
	CategoriaID   int     `json:"categoria_id"`
	CategoriaNome string  `json:"categoria_nome"`
	TotalValor    float64 `json:"total_valor"`
}

type RelatorioFinanceiroDTO struct {
	TotalMesPorCategoria   []TotalPorCategoriaDTO `json:"total_mes_por_categoria"`
	TotalAno               float64                `json:"total_ano"`
	TotalGeralPorCategoria []TotalPorCategoriaDTO `json:"total_geral_por_categoria"`
	TotalServicosPrestados int64                  `json:"total_servicos_prestados"`
}