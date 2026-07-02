package main

import (
	"fmt"
	"net/http"
	"work-backend/interno/caso_de_uso"
	"work-backend/interno/conectores"
	"work-backend/interno/repositorio"
	"work-backend/interno/seguranca"
)

func main() {
	// 1. Instancia o mecanismo de segurança (Bcrypt):
	hasher := seguranca.NovoBcryptHasher()

	// 2. Instancia o repositório em memória:
	repo := repositorio.NovoMemoriaAdministradorRepositorio()

	// 3. Instancia o caso de uso injetando o repositório e o componente de segurança:
	casoUso := caso_de_uso.NovoAdministradorCasoDeUso(repo, hasher)

	// 4. Instancia o conector HTTP injetando o caso de uso:
	conectorAdmin := conectores.NovoAdministradorConector(casoUso)

	// 5. Define a rota da API que o React/React-Native vão chamar:
	// Use este caminho exato ao testar no navegador ou aplicativo
	http.Handle("/administradores", conectorAdmin)

	// 6. Instancia o servidor local:
	fmt.Println("Servidor Go rodando em http://localhost:8080")
	// Mudado de Println para Printf para exibir a variável 'erro' corretamente se algo falhar
	if erro := http.ListenAndServe(":8080", nil); erro != nil {
		fmt.Printf("Erro ao iniciar o servidor: %v\n", erro)
	}
}