package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"work-backend/interno/caso_de_uso"
	"work-backend/interno/conectores"
	"work-backend/interno/repositorio"
	"work-backend/interno/seguranca"
	_"github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
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

	bancoDeDados, erro := sql.Open("mysql", dsn)
	if erro != nil {
		log.Fatalf("Erro ao configurar o cliente MySQL: %v", erro)
	}
	defer bancoDeDados.Close()

	if erro := bancoDeDados.Ping(); erro != nil {
		log.Fatalf("Não foi possível conectar ao MySQL: %v", erro)
	}
	fmt.Println("Conexão com o MySQL realizada com sucesso!")

	//==== Instâncias:

	// 1. Instancia o mecanismo de segurança (Bcrypt):
	hasher := seguranca.NovoBcryptHasher()

	// 2. Instancia o repositório em memória:
	// Opcional: Quando tiver o repositório do banco pronto, substitua por:
	// repo := repositorio.NovoMySQLAdministradorRepositorio(bancoDeDados)
	repo := repositorio.NovoMemoriaAdministradorRepositorio()

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

	// 7. Instancia o servidor local:
	fmt.Println("Servidor Go rodando em http://localhost:8080")
	if erro := http.ListenAndServe(":8080", nil); erro != nil {
		fmt.Printf("Erro ao iniciar o servidor: %v\n", erro)
	}
}