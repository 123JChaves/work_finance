package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"work-backend/interno/caso_de_uso"
	"work-backend/interno/conectores"
	"work-backend/interno/entidade"
	"work-backend/interno/repositorio"
	"work-backend/interno/seguranca"

	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// HabilitarCORS injeta os headers necessários para comunicação com React / React Native
func HabilitarCORS(proximo http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*") 
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS, PATCH")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		proximo.ServeHTTP(w, r)
	})
}

func main() {
	if erro := godotenv.Load(); erro != nil {
		fmt.Println("Aviso: Nenhum arquivo .env encontrado, usando variáveis de ambiente do sistema.")
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
		os.Getenv("DB_USUARIO"),
		os.Getenv("BD_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NOME"),
	)

	bancoDeDados, erro := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if erro != nil {
		log.Fatalf("Erro ao conectar ao MySQL via GORM: %v", erro)
	}

	// EXECUÇÃO DAS MIGRATIONS:
	// A ordem garante que tabelas independentes ou pais criem suas colunas 
	// antes que as filhas tentem amarrar Constraints de Foreign Key.
	erro = bancoDeDados.AutoMigrate(
		&entidade.Administrador{},
		&entidade.Categoria{},
		&entidade.CategoriaDespesaEmpresa{},
		&entidade.CategoriaDespesaUsuario{},
		&entidade.CategoriaObjetoPatrimonial{},
		&entidade.Usuario{},
		&entidade.Empresa{},
		&entidade.Servico{},
		&entidade.Cliente{},
		&entidade.DespesaEmpresa{},
		&entidade.DespesaUsuario{},
		&entidade.Patrimonio{},
		&entidade.PatrimonioEmpresa{},
		&entidade.ObjetoPatrimonialCorp{},
		&entidade.ObjetoPatrimonialPess{},
	)
	if erro != nil {
		log.Fatalf("Erro ao executar migration do banco: %v", erro)
	}
	fmt.Println("Conexão com o MySQL e Migrations realizadas com sucesso via GORM!")

	//=============================================================================
	// INJEÇÃO DE DEPENDÊNCIAS
	//=============================================================================

	hasher := seguranca.NovoBcryptHasher()

	// 1. Administrador
	repoAdmin := repositorio.NovoGORMAdministradorRepositorio(bancoDeDados)
	casoUsoAdmin := caso_de_uso.NovoAdministradorCasoDeUso(repoAdmin, hasher)
	conectorAdmin := conectores.NovoAdministradorConector(casoUsoAdmin)

	// 2. Categoria (Serviços)
	repoCategoria := repositorio.NovoGORMCategoriaRepositorio(bancoDeDados)
	casoUsoCategoria := caso_de_uso.NovoCategoriaCasoDeUso(repoCategoria)
	conectorCategoria := conectores.NovoCategoriaConector(casoUsoCategoria)

	// 3. Categoria Despesa Empresa
	repoCatDespesaEmp := repositorio.NovoGORMCategoriaDespesaEmpresaRepo(bancoDeDados)
	casoUsoCatDespesaEmp := caso_de_uso.NovoCategoriaDespesaEmpresaCasoDeUso(repoCatDespesaEmp)
	conectorCatDespesaEmp := conectores.NovoCategoriaDespesaEmpresaConector(casoUsoCatDespesaEmp)

	// 4. Categoria Despesa Usuário
	repoCatDespesaUser := repositorio.NovoGORMCategoriaDespesaUsuarioRepo(bancoDeDados)
	casoUsoCatDespesaUser := caso_de_uso.NovoCategoriaDespesaUsuarioCasoDeUso(repoCatDespesaUser)
	conectorCatDespesaUser := conectores.NovoCategoriaDespesaUsuarioConector(casoUsoCatDespesaUser)

	// 5. Categoria Objeto Patrimonial
	repoCatObjPatrimonial := repositorio.NovoGORMCategoriaObjetoPatrimonialRepo(bancoDeDados)
	casoUsoCatObjPatrimonial := caso_de_uso.NovoCategoriaObjetoPatrimonialCasoDeUso(repoCatObjPatrimonial)
	conectorCatObjPatrimonial := conectores.NovoCategoriaObjetoPatrimonialConector(casoUsoCatObjPatrimonial)

	// 6. Empresa
	repoEmpresa := repositorio.NovoGORMEmpresaRepo(bancoDeDados)
	casoUsoEmpresa := caso_de_uso.NovoEmpresaCasoDeUso(repoEmpresa)
	conectorEmpresa := conectores.NovoEmpresaConector(casoUsoEmpresa)

	// 7. Usuário
	repoUsuario := repositorio.NovoGORMUsuarioRepository(bancoDeDados)
	casoUsoUsuario := caso_de_uso.NovoUsuarioCasoDeUso(repoUsuario, hasher)
	conectorUsuario := conectores.NovoUsuarioConector(casoUsoUsuario)

	// 8. Serviço (Faturamento Bruto da Empresa)
	repoServico := repositorio.NovoGORMServicoRepositorio(bancoDeDados)
	casoUsoServico := caso_de_uso.NovoServicoCasoDeUso(repoServico, repoCategoria)
	conectorServico := conectores.NovoServicoConector(casoUsoServico)

	// 9. Cliente
	repoCliente := repositorio.NovoGORMClienteRepo(bancoDeDados)
	casoUsoCliente := caso_de_uso.NovoClienteCasoDeUso(repoCliente, repoEmpresa)
	conectorCliente := conectores.NovoClienteConector(casoUsoCliente)

	// 10. Despesa Empresa
	repoDespesaEmp := repositorio.NovoGORMDespesaEmpresaRepo(bancoDeDados)
	casoUsoDespesaEmp := caso_de_uso.NovoDespesaEmpresaCasoDeUso(repoDespesaEmp, repoCatDespesaEmp, repoEmpresa)
	conectorDespesaEmp := conectores.NovoDespesaEmpresaConector(casoUsoDespesaEmp)

	// 11. Despesa Usuário
	repoDespesaUser := repositorio.NovoGORMDespesaUsuarioRepo(bancoDeDados)
	casoUsoDespesaUser := caso_de_uso.NovoDespesaUsuarioCasoDeUso(repoDespesaUser, repoCatDespesaUser, repoUsuario)
	conectorDespesaUser := conectores.NovoDespesaUsuarioConector(casoUsoDespesaUser)

	// 12. Patrimônio Pessoal
	repoPatrimonio := repositorio.NovoGORMPatrimonioRepo(bancoDeDados)
	casoUsoPatrimonio := caso_de_uso.NovoPatrimonioCasoDeUso(repoPatrimonio, repoUsuario)
	conectorPatrimonio := conectores.NovoPatrimonioConector(casoUsoPatrimonio)

	// 13. Patrimônio Empresa
	repoPatrimonioEmp := repositorio.NovoGORMPatrimonioEmpresaRepo(bancoDeDados)
	casoUsoPatrimonioEmp := caso_de_uso.NovoPatrimonioEmpresaCasoDeUso(repoPatrimonioEmp, repoEmpresa)
	conectorPatrimonioEmp := conectores.NovoPatrimonioEmpresaConector(casoUsoPatrimonioEmp)

	// 14. Objeto Patrimonial Corp
	repoObjCorp := repositorio.NovoGORMObjetoPatrimonialCorpRepo(bancoDeDados)
	casoUsoObjCorp := caso_de_uso.NovoObjetoPatrimonialCorpCasoDeUso(repoObjCorp, repoEmpresa, repoPatrimonioEmp)
	conectorObjCorp := conectores.NovoObjetoPatrimonialCorpConector(casoUsoObjCorp)

	// 15. Objeto Patrimonial Pess
	repoObjPess := repositorio.NovoGORMObjetoPatrimonialPessRepo(bancoDeDados)
	casoUsoObjPess := caso_de_uso.NovoObjetoPatrimonialPessCasoDeUso(repoObjPess, repoPatrimonio, repoCatObjPatrimonial)
	conectorObjPess := conectores.NovoObjetoPatrimonialPessConector(casoUsoObjPess)

	//=============================================================================
	// MAPEAMENTO DE ROTAS HTTP (REST)
	//=============================================================================

	// Rota raiz "/"
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		usuarioEnv := os.Getenv("USUARIO")
		if usuarioEnv == "" {
			usuarioEnv = "Visitante"
		}
		mensagem := fmt.Sprintf("Bem-vindo, %s", usuarioEnv)
		json.NewEncoder(w).Encode(mensagem)
	})

	// Associação de Handlers encapsulados pelo CORS Middleware
	http.Handle("/administradores/", HabilitarCORS(conectorAdmin))
	http.Handle("/administradores", HabilitarCORS(conectorAdmin))

	http.Handle("/categorias/", HabilitarCORS(conectorCategoria))
	http.Handle("/categorias", HabilitarCORS(conectorCategoria))

	http.Handle("/categorias-despesas-empresa/", HabilitarCORS(conectorCatDespesaEmp))
	http.Handle("/categorias-despesas-empresa", HabilitarCORS(conectorCatDespesaEmp))

	http.Handle("/categorias-despesas-usuario/", HabilitarCORS(conectorCatDespesaUser))
	http.Handle("/categorias-despesas-usuario", HabilitarCORS(conectorCatDespesaUser))

	http.Handle("/categorias-objetos-patrimoniais/", HabilitarCORS(conectorCatObjPatrimonial))
	http.Handle("/categorias-objetos-patrimoniais", HabilitarCORS(conectorCatObjPatrimonial))

	http.Handle("/empresas/", HabilitarCORS(conectorEmpresa))
	http.Handle("/empresas", HabilitarCORS(conectorEmpresa))

	http.Handle("/usuarios/", HabilitarCORS(conectorUsuario))
	http.Handle("/usuarios", HabilitarCORS(conectorUsuario))

	http.Handle("/servicos/", HabilitarCORS(conectorServico))
	http.Handle("/servicos", HabilitarCORS(conectorServico))

	http.Handle("/clientes/", HabilitarCORS(conectorCliente))
	http.Handle("/clientes", HabilitarCORS(conectorCliente))

	http.Handle("/despesas-empresa/", HabilitarCORS(conectorDespesaEmp))
	http.Handle("/despesas-empresa", HabilitarCORS(conectorDespesaEmp))

	http.Handle("/despesas-usuario/", HabilitarCORS(conectorDespesaUser))
	http.Handle("/despesas-usuario", HabilitarCORS(conectorDespesaUser))

	http.Handle("/patrimonios/", HabilitarCORS(conectorPatrimonio))
	http.Handle("/patrimonios", HabilitarCORS(conectorPatrimonio))

	http.Handle("/patrimonios-empresa/", HabilitarCORS(conectorPatrimonioEmp))
	http.Handle("/patrimonios-empresa", HabilitarCORS(conectorPatrimonioEmp))

	http.Handle("/objetos-patrimoniais-corp/", HabilitarCORS(conectorObjCorp))
	http.Handle("/objetos-patrimoniais-corp", HabilitarCORS(conectorObjCorp))

	http.Handle("/objetos-patrimoniais-pess/", HabilitarCORS(conectorObjPess))
	http.Handle("/objetos-patrimoniais-pess", HabilitarCORS(conectorObjPess))

	//=============================================================================
	// INICIALIZAÇÃO DO SERVIDOR
	//=============================================================================
	porta := os.Getenv("PORT")
	if porta == "" {
		porta = "8080"
	}
	enderecoServidor := fmt.Sprintf(":%s", porta)

	fmt.Printf("Servidor Go rodando perfeitamente em http://localhost:%s\n", porta)
	if erro := http.ListenAndServe(enderecoServidor, nil); erro != nil {
		fmt.Printf("Erro crítico ao iniciar o servidor: %v\n", erro)
	}
}
