package validador

import (
	"regexp"
	"strings" // Adicionado para manipulação de strings/NewReplacer/TrimSpace
	"unicode" // Adicionado para checagem de IsUpper/ToUpper em runes
)

var RegexNomeSeguro = regexp.MustCompile(`^[a-zA-Z0-9À-ÿ\s]+$`)
var RegexEmailValido = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

func NomeSeguro(nome string) bool {
	return RegexNomeSeguro.MatchString(nome)
}

func EmailValido(email string) bool {
	return RegexEmailValido.MatchString(email)
}

// ComecaComMaiuscula valida se o primeiro caractere de uma string (como Categorias)
// é uma letra maiúscula. Trata corretamente caracteres acentuados Unicode (ex: "Água", "Óleo").
func ComecaComMaiuscula(texto string) bool {
	textoLimpo := strings.TrimSpace(texto)
	if textoLimpo == "" {
		return false
	}
	caracteres := []rune(textoLimpo)
	return unicode.IsUpper(caracteres[0])
}

// CapitalizarPrimeiraLetra corrige automaticamente strings transformando o primeiro caractere
// em maiúsculo (útil caso queira corrigir o dado em vez de retornar um erro ao usuário).
func CapitalizarPrimeiraLetra(texto string) string {
	textoLimpo := strings.TrimSpace(texto)
	if textoLimpo == "" {
		return textoLimpo
	}
	caracteres := []rune(textoLimpo)
	caracteres[0] = unicode.ToUpper(caracteres[0]) // CORRIGIDO: Atribuição correta na posição [0]
	return string(caracteres)
}

// EhPessoaJuridica limpa máscaras de documentos (como pontos, traços e barras)
// e avalia se o documento se trata de um CNPJ (maior que 11 dígitos).
func EhPessoaJuridica(documento string) bool {
	documentoLimpo := strings.NewReplacer(".", "", "-", "", "/", "").Replace(documento)
	return len(documentoLimpo) > 11
}

// DocumentoPossuiTamanhoValido verifica se o CPF ou CNPJ limpo possui a quantidade exata
// de dígitos exigidos pela Receita Federal (11 para CPF ou 14 para CNPJ).
func DocumentoPossuiTamanhoValido(documento string) bool {
	documentoLimpo := strings.NewReplacer(".", "", "-", "", "/", "").Replace(documento)
	tamanho := len(documentoLimpo)
	return tamanho == 11 || tamanho == 14
}