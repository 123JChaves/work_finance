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

func main() {
	// Carrega o arquivo .env
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

	// Alterado: Substituído o sql.Open pelo gorm.Open
	bancoDeDados, erro := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if erro != nil {
		log.Fatalf("Erro ao conectar ao MySQL via GORM: %v", erro)
	}

	// Executa o AutoMigrate para criar/atualizar a tabela 'administradores' no MySQL
	erro = bancoDeDados.AutoMigrate(&entidade.Administrador{})
	if erro != nil {
		log.Fatalf("Erro ao executar migration do banco: %v", erro)
	}
	fmt.Println("Conexão com o MySQL realizada com sucesso via GORM!")

	//==== Instâncias:

	// 1. Instancia o mecanismo de segurança (Bcrypt):
	hasher := seguranca.NovoBcryptHasher()

	// 2. Alterado: Agora instancia o repositório GORM injetando o banco de dados
	repo := repositorio.NovoGORMAdministradorRepositorio(bancoDeDados)

	// 3. Instancia o caso de uso injetando o repositório e o componente de segurança:
	casoUso := caso_de_uso.NovoAdministradorCasoDeUso(repo, hasher)

	// 4. Instancia o conector HTTP injetando o caso de uso:
	conectorAdmin := conectores.NovoAdministradorConector(casoUso)

	// 5. Rota raiz "/" enviando a mensagem com a variável do .env
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		usuario := os.Getenv("USUARIO")

		if usuario == "" {
			usuario = "Visitante"
		}

		mensagem := fmt.Sprintf("Bem-vindo, %s", usuario)
		json.NewEncoder(w).Encode(mensagem)
	})

	// 6. Define a rota da API que o React/React-Native vão chamar:
	http.Handle("/administradores", conectorAdmin)

	// 7. Instancia o servidor local obtendo a porta do arquivo .env
	porta := os.Getenv("PORT")

	// Formata a porta adicionando os dois pontos (":") exigidos pelo ListenAndServe
	enderecoServidor := fmt.Sprintf(":%s", porta)

	fmt.Printf("Servidor Go rodando em http://localhost:%s\n", porta)
	if erro := http.ListenAndServe(enderecoServidor, nil); erro != nil {
		fmt.Printf("Erro ao iniciar o servidor: %v\n", erro)
	}
}