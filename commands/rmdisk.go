package commands

import (
	"errors"
	"fmt"
	"os" // Paquete para trabajar con el sistema de archivos, necesario para eliminar archivos
	"regexp"
	"strings"
)

// RMDISK estructura que representa el comando rmdisk con su parámetro
type RMDISK struct {
	path string // Ruta del archivo del disco
}

// ParserRmdisk parsea el comando rmdisk y devuelve una instancia de RMDISK
func ParserRmdisk(tokens []string) (string, error) {
	cmd := &RMDISK{} // Crea una nueva instancia de RMDISK

	// Unir tokens en una sola cadena y luego dividir por espacios, respetando las comillas
	args := strings.Join(tokens, " ")
	// Expresión regular para encontrar el parámetro del comando rmdisk
	re := regexp.MustCompile(`-path="[^"]+"|-path=[^\s]+`)
	// Encuentra todas las coincidencias de la expresión regular en la cadena de argumentos
	matches := re.FindAllString(args, -1)

	// Verificar que todos los tokens fueron reconocidos por la expresión regular
	if len(matches) != len(tokens) {
		// Identificar el parámetro inválido
		for _, token := range tokens {
			if !re.MatchString(token) {
				return "", fmt.Errorf("parámetro inválido: %s", token)
			}
		}
	}

	// Itera sobre cada coincidencia encontrada
	for _, match := range matches {
		// Divide cada parte en clave y valor usando "=" como delimitador
		kv := strings.SplitN(match, "=", 2)
		if len(kv) != 2 {
			return "", fmt.Errorf("formato de parámetro inválido: %s", match)
		}
		key, value := strings.ToLower(kv[0]), kv[1]

		// Remover comillas del valor si están presentes
		if strings.HasPrefix(value, "\"") && strings.HasSuffix(value, "\"") {
			value = strings.Trim(value, "\"")
		}

		// Switch para manejar el parámetro -path
		switch key {
		case "-path":
			// Verifica que el path no esté vacío
			if value == "" {
				return "", errors.New("el path no puede estar vacío")
			}
			cmd.path = value
		default:
			// Si el parámetro no es reconocido, devuelve un error
			return "", fmt.Errorf("parámetro desconocido: %s", key)
		}
	}

	// Verifica que el parámetro -path haya sido proporcionado
	if cmd.path == "" {
		return "", errors.New("faltan parámetros requeridos: -path")
	}

	// Aquí va la lógica para eliminar el archivo
	err := os.Remove(cmd.path)
	if err != nil {
		// Si ocurre un error al intentar eliminar el archivo, devuelve el error
		return "", fmt.Errorf("no se pudo eliminar el archivo: %v", err)
	}

	// Devuelve el mensaje de éxito si el archivo fue eliminado correctamente
	return fmt.Sprintf("RMDISK: Disco %s eliminado correctamente.", cmd.path), nil
}
