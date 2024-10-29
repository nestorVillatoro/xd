package analyzer

import (
	commands "P1_202200252/commands" // Importa el paquete "commands" desde el directorio "PRUEBA01/commands"
	"errors"                         // Importa el paquete "errors" para manejar errores
	"fmt"                            // Importa el paquete "fmt" para formatear e imprimir texto
	"strings"                        // Importa el paquete "strings" para manipulación de cadenas
)

// Analyzer analiza el comando de entrada y ejecuta la acción correspondiente
func Analyzer(input string) (string, error) {
	// Divide la entrada en tokens usando espacios en blanco como delimitadores
	tokens := strings.Fields(input)

	// Si no se proporcionó ningún comando, devuelve un error
	if len(tokens) == 0 {
		return "", errors.New("no se proporcionó ningún comando o es un salto de línea en la entrada")
	}
	tokens[0] = strings.ToLower(tokens[0])
	// Switch para manejar diferentes comandos
	switch tokens[0] {
	case "mkdisk":
		// Llama a la función ParseMkdisk
		return commands.ParserMkdisk(tokens[1:])
	case "rmdisk":
		// Llama a la función CommandRmdisk
		return commands.ParserRmdisk(tokens[1:])
	case "fdisk":
		// Llama a la función CommandFdisk
		return commands.ParserFdisk(tokens[1:])
	case "mount":
		// Llama a la función CommandMount
		return commands.ParserMount(tokens[1:])
	case "mkfs":
		// Llama a la función CommandMkfs
		return commands.ParserMkfs(tokens[1:])
	case "mkdir":
		// Llama a la función CommandMkdir
		return commands.ParserMkdir(tokens[1:])
	case "mkfile":
		// Llama a la función CommandMkfile
		return commands.ParserMkfile(tokens[1:])
	case "rep":
		// Llama a la función CommandRep
		return commands.ParserRep(tokens[1:])

	case "#":
		result := strings.Join(tokens[1:], " ")
		return result, nil
	default:
		// Si el comando no es reconocido, devuelve un error
		return "", fmt.Errorf("comando desconocido: %s", tokens[0])
	}
}
