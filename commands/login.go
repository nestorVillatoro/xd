package commands

import (
	global "P1_202200252/global"
	estructuras "P1_202200252/structures"
	"bytes"
	"errors"
	"fmt"
	"strings"
)

type LOGIN struct {
	user string
	pass string
	id   string
}

// login -user=root -pass=123 -id=062A
// login -user="mi usuario" -pass="mi pwd" -id=062A
func ParserLogin(tokens []string) (*LOGIN, string, error) {
	cmd := &LOGIN{} //se crea instancia de clase
	// Unir tokens en una sola cadena y luego dividir por espacios, respetando las comillas

	for _, match := range tokens {
		// Divide cada parte en clave y valor usando "=" como delimitador
		kv := strings.SplitN(match, "=", 2)
		if len(kv) != 2 {
			return nil, "formato de parámetro inválido: " + match, fmt.Errorf("formato de parámetro inválido: %s", match)
		}
		key, value := strings.ToLower(kv[0]), kv[1]
		// Remove quotes from value if present
		if strings.HasPrefix(value, "\"") && strings.HasSuffix(value, "\"") {
			value = strings.Trim(value, "\"")
		}

		// Switch para manejar diferentes parámetros
		switch key {
		case "-user":
			// Verificar que el user no este vacio
			if value == "" {
				return nil, "el user no puede estar vacio", errors.New("el user no puede estar vacío")
			}
			cmd.user = value

		case "-pass":
			// Verifica que el path no esté vacío
			if value == "" {
				return nil, "el pass no puede estar vacio", errors.New("el pass no puede estar vacío")
			}
			cmd.pass = value
		case "-id":
			if value == "" {

				return nil, "el id no puede estar vacio", errors.New("el id no puede estar vacio")

			}
			cmd.id = value
		default:
			// Si el parámetro no es reconocido, devuelve un error
			return nil, "parametro desconocido: " + key, fmt.Errorf("parámetro desconocido: %s", key)
		}
	}

	if cmd.user == "" {
		return nil, "falta parametros requeridos: -user", errors.New("falta parametros requeridos: -user")

	}
	if cmd.pass == "" {
		return nil, "falta parametros requeridos: -pass", errors.New("falta parametros requeridos: -pass")

	}
	if cmd.id == "" {
		return nil, "falta parametros requeridos: -id", errors.New("falta parametros requeridos: -id")

	}

	posible_error, err := commandLogin(cmd)
	if err != nil {
		fmt.Println("Error:", err)
		return nil, posible_error, err
	}

	return cmd, posible_error, nil
}

func commandLogin(login *LOGIN) (string, error) {

	//PRIMERO VERIFICAMOS SI EL USUARIO ESTA LOGEADO
	if global.IsUserLogged(login.user) {
		return fmt.Sprintf("Error: La sesión ya está iniciada para el usuario %s", login.user), errors.New("sesión ya iniciada")
	}

	// Obtener la partición montada usando el ID proporcionado
	mountedPartition, partitionPath, err := global.GetMountedPartition(login.id)
	if err != nil {
		return "error al obtener la particion montada", fmt.Errorf("error al obtener la partición montada: %v", err)
	}

	// Leer el superbloque
	superBlock := &estructuras.SuperBlock{}
	err = superBlock.Deserialize(partitionPath, int64(mountedPartition.Part_start))
	if err != nil {
		return "", fmt.Errorf("error al leer el superbloque: %v", err)
	}

	// Buscar el inodo de users.txt
	rootInode := &estructuras.Inode{}
	err = rootInode.Deserialize(partitionPath, int64(superBlock.S_inode_start))
	if err != nil {
		return "", fmt.Errorf("error al leer el inodo raíz: %v", err)
	}

	var usersInode *estructuras.Inode
	//var usersInodeIndex int32
	for _, blockIndex := range rootInode.I_block {
		if blockIndex == -1 {
			break
		}
		folderBlock := &estructuras.FolderBlock{}
		err = folderBlock.Deserialize(partitionPath, int64(superBlock.S_block_start+blockIndex*superBlock.S_block_size))
		if err != nil {
			return "", fmt.Errorf("error al leer el bloque de carpeta: %v", err)
		}
		for _, content := range folderBlock.B_content {
			if string(bytes.Trim(content.B_name[:], "\x00")) == "users.txt" {
				usersInode = &estructuras.Inode{}
				err = usersInode.Deserialize(partitionPath, int64(superBlock.S_inode_start+content.B_inodo*superBlock.S_inode_size))
				if err != nil {
					return "", fmt.Errorf("error al leer el inodo de users.txt: %v", err)
				}
				//usersInodeIndex = content.B_inodo
				break
			}
		}
		if usersInode != nil {
			break
		}
	}

	if usersInode == nil {
		return "", fmt.Errorf("no se encontró el archivo users.txt")
	}

	// Leer el contenido de users.txt
	fileBlock := &estructuras.FileBlock{}
	err = fileBlock.Deserialize(partitionPath, int64(superBlock.S_block_start+usersInode.I_block[0]*superBlock.S_block_size))
	if err != nil {
		return "", fmt.Errorf("error al leer el bloque de archivo: %v", err)
	}

	content := string(bytes.Trim(fileBlock.B_content[:], "\x00"))
	lines := strings.Split(content, "\n")

	// Verificar si el usuario y la contraseña son correctos
	for _, line := range lines {
		parts := strings.Split(line, ",")
		if len(parts) == 5 && parts[1] == "U" {
			if strings.TrimSpace(parts[3]) == login.user {
				if strings.TrimSpace(parts[4]) == login.pass {
					global.RegisterSession(login.user, login.pass, login.id, login.user == "root")
					return fmt.Sprintf("Inicio de sesión exitoso para el usuario %s", login.user), nil
				}
				return "contraseña incorrecta para el usuario: " + login.user, fmt.Errorf("contraseña incorrecta para el usuario %s", login.user)
			}
		}
	}

	return "el usuario no existe: " + login.user, fmt.Errorf("el usuario %s no existe", login.user)
}
