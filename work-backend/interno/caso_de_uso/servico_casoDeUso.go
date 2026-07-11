package caso_de_uso

import (
	"errors"
	"fmt"
	"work-backend/interno/entidade"
	"work-backend/interno/regras"
	"gorm.io/gorm"
)

// Interface customizada local estendida para ler o DB sem quebrar o contrato do domínio
type ServicoRepositorioGORM interface {
	entidade.ServicoRepositorio
	DB() *gorm.DB
}

type Servico_casoDeUso struct {
	repositorio   ServicoRepositorioGORM
	categoriaRepo entidade.CategoriaRepositorio
}

func NovoServicoCasoDeUso(repo ServicoRepositorioGORM, catRepo entidade.CategoriaRepositorio) *Servico_casoDeUso {
	return &Servico_casoDeUso{repositorio: repo, categoriaRepo: catRepo}
}

func (uc *Servico_casoDeUso) Cadastrar(servico *entidade.Servico) error {
	if err := regras.ValidarServico(servico); err != nil {
		return err
	}
	categoria, err := uc.categoriaRepo.FindByID(servico.CategoriaID)
	if err != nil || categoria == nil {
		return errors.New("a categoria informada para este serviço não existe")
	}
	servico.Categoria = *categoria
	return uc.repositorio.Create(servico)
}

func (uc *Servico_casoDeUso) Listar() ([]*entidade.Servico, error) {
	return uc.repositorio.FindAll()
}

func (uc *Servico_casoDeUso) Buscar(id int) (*entidade.Servico, error) {
	if id <= 0 {
		return nil, errors.New("o ID do serviço deve ser válido")
	}
	return uc.repositorio.FindByID(id)
}

func (uc *Servico_casoDeUso) Atualizar(id int, nome string, valor float64, categoriaID int) error {
	servico, err := uc.Buscar(id)
	if err != nil {
		return err
	}

	// Atribui os campos novos caso tenham sido enviados na requisição
	if nome != "" { 
		servico.Nome = nome 
	}
	if valor > 0 { 
		servico.Valor = valor 
	}
	if categoriaID > 0 { 
		servico.CategoriaID = categoriaID 
	}

	// Invoca a função pura do pacote de regras passando os novos dados numéricos
	if err := regras.ValidarServico(servico); err != nil {
		return err
	}

	return uc.repositorio.Update(servico)
}

func (uc *Servico_casoDeUso) Deletar(id int) error {
	if _, err := uc.Buscar(id); err != nil {
		return err
	}
	return uc.repositorio.Delete(id)
}

func (uc *Servico_casoDeUso) ObterDadosRelatorio(ano int, mes int) (*RelatorioFinanceiroDTO, error) {
	var relatorio RelatorioFinanceiroDTO
	db := uc.repositorio.DB()

	err := db.Table("servicos").
		Select("servicos.categoria_id, categorias.nome as categoria_nome, SUM(servicos.valor) as total_valor").
		Joins("JOIN categorias ON categorias.id = servicos.categoria_id").
		Where("YEAR(servicos.data) = ? AND MONTH(servicos.data) = ?", ano, mes).
		Group("servicos.categoria_id, categorias.nome").
		Scan(&relatorio.TotalMesPorCategoria).Error
	if err != nil {
		return nil, fmt.Errorf("erro ao calcular total mensal por categoria: %w", err)
	}

	err = db.Table("servicos").
		Select("COALESCE(SUM(valor), 0.00)").
		Where("YEAR(data) = ?", ano).
		Scan(&relatorio.TotalAno).Error
	if err != nil {
		return nil, fmt.Errorf("erro ao calcular total anual: %w", err)
	}

	err = db.Table("servicos").
		Select("servicos.categoria_id, categorias.nome as categoria_nome, SUM(servicos.valor) as total_valor").
		Joins("JOIN categorias ON categorias.id = servicos.categoria_id").
		Group("servicos.categoria_id, categorias.nome").
		Scan(&relatorio.TotalGeralPorCategoria).Error
	if err != nil {
		return nil, fmt.Errorf("erro ao calcular total histórico por categoria: %w", err)
	}

	err = db.Table("servicos").Count(&relatorio.TotalServicosPrestados).Error
	if err != nil {
		return nil, fmt.Errorf("erro ao contar total de serviços: %w", err)
	}

	return &relatorio, nil
}