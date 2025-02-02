package commands

import (
	structures "P1_202200252/structures"
	utils "P1_202200252/utils"
	"errors"  // Paquete para manejar errores y crear nuevos errores con mensajes personalizados
	"fmt"     // Paquete para formatear cadenas y realizar operaciones de entrada/salida
	"regexp"  // Paquete para trabajar con expresiones regulares, útil para encontrar y manipular patrones en cadenas
	"strconv" // Paquete para convertir cadenas a otros tipos de datos, como enteros
	"strings" // Paquete para manipular cadenas, como unir, dividir, y modificar contenido de cadenas
)

// FDISK estructura que representa el comando fdisk con sus parámetros
type FDISK struct {
	size int    // Tamaño de la partición
	unit string // Unidad de medida del tamaño (K o M)
	fit  string // Tipo de ajuste (BF, FF, WF)
	path string // Ruta del archivo del disco
	typ  string // Tipo de partición (P, E, L)
	name string // Nombre de la partición
}

// CommandFdisk parsea el comando fdisk y devuelve una instancia de FDISK
func ParserFdisk(tokens []string) (string, error) {
	cmd := &FDISK{} // Crea una nueva instancia de FDISK

	// Unir tokens en una sola cadena y luego dividir por espacios, respetando las comillas
	args := strings.Join(tokens, " ")
	// Expresión regular para encontrar los parámetros del comando fdisk
	re := regexp.MustCompile(`-size=\d+|-unit=[kKmMbB]|-fit=[bBfF]{2}|-path="[^"]+"|-path=[^\s]+|-type=[pPeElL]|-name="[^"]+"|-name=[^\s]+`)
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

		// Remove quotes from value if present
		if strings.HasPrefix(value, "\"") && strings.HasSuffix(value, "\"") {
			value = strings.Trim(value, "\"")
		}

		// Switch para manejar diferentes parámetros
		switch key {
		case "-size":
			// Convierte el valor del tamaño a un entero
			size, err := strconv.Atoi(value)
			if err != nil || size <= 0 {
				return "", errors.New("el tamaño debe ser un número entero positivo")
			}
			cmd.size = size
		case "-unit":
			// Verifica que la unidad sea "K" o "M"
			if value != "K" && value != "M" && value != "m" && value != "k" && value != "b" && value != "B" {
				return "", errors.New("la unidad debe ser K , M o B")
			}
			cmd.unit = strings.ToUpper(value)
		case "-fit":
			// Verifica que el ajuste sea "BF", "FF" o "WF"
			value = strings.ToUpper(value)
			if value != "BF" && value != "FF" && value != "WF" {
				return "", errors.New("el ajuste debe ser BF, FF o WF")
			}
			cmd.fit = value
		case "-path":
			// Verifica que el path no esté vacío
			if value == "" {
				return "", errors.New("el path no puede estar vacío")
			}
			cmd.path = value
		case "-type":
			// Verifica que el tipo sea "P", "E" o "L"
			value = strings.ToUpper(value)
			if value != "P" && value != "E" && value != "L" {
				return "", errors.New("el tipo debe ser P, E o L")
			}
			cmd.typ = value
		case "-name":
			// Verifica que el nombre no esté vacío
			if value == "" {
				return "", errors.New("el nombre no puede estar vacío")
			}
			cmd.name = value
		default:
			// Si el parámetro no es reconocido, devuelve un error
			return "", fmt.Errorf("parámetro desconocido: %s", key)
		}
	}

	// Verifica que los parámetros -size, -path y -name hayan sido proporcionados
	if cmd.size == 0 {
		return "", errors.New("faltan parámetros requeridos: -size")
	}
	if cmd.path == "" {
		return "", errors.New("faltan parámetros requeridos: -path")
	}
	if cmd.name == "" {
		return "", errors.New("faltan parámetros requeridos: -name")
	}

	// Si no se proporcionó la unidad, se establece por defecto a "M"
	if cmd.unit == "" {
		cmd.unit = "M"
	}

	// Si no se proporcionó el ajuste, se establece por defecto a "FF"
	if cmd.fit == "" {
		cmd.fit = "WF"
	}

	// Si no se proporcionó el tipo, se establece por defecto a "P"
	if cmd.typ == "" {
		cmd.typ = "P"
	}

	// Crear la partición con los parámetros proporcionados
	err := commandFdisk(cmd)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("FDISK: Partición %s creada correctamente de %d%s de tipo %s con ajuste %s.", cmd.name, cmd.size, cmd.unit, cmd.typ, cmd.fit), nil // Devuelve el comando FDISK creado
}

func commandFdisk(fdisk *FDISK) error {
	// Convertir el tamaño a bytes
	sizeBytes, err := utils.ConvertToBytes(fdisk.size, fdisk.unit)
	if err != nil {
		fmt.Println("Error converting size:", err)
		return err
	}

	if fdisk.typ == "P" {
		// Crear partición primaria
		err = createPrimaryPartition(fdisk, sizeBytes)
		if err != nil {
			fmt.Println("Error creando partición primaria:", err)
			return err
		}
	} else if fdisk.typ == "E" {
		// Crear partición extendida
		err = createExtendedPartition(fdisk, sizeBytes)
		if err != nil {
			fmt.Println("Error creando partición extendida:", err)
			return err
		}
	} else if fdisk.typ == "L" {
		// Crear partición logica
		err = createLogicPartition(fdisk, sizeBytes)
		if err != nil {
			fmt.Println("Error creando partición logica:", err)
			return err
		}
	}

	return nil
}

func createPrimaryPartition(fdisk *FDISK, sizeBytes int) error {
	// Crear una instancia de MBR
	var mbr structures.MBR

	// Deserializar la estructura MBR desde un archivo binario
	err := mbr.Deserialize(fdisk.path)
	if err != nil {
		fmt.Println("Error deserializando el MBR:", err)
		return err
	}

	// Obtener la primera partición disponible
	availablePartition, startPartition, indexPartition := mbr.GetFirstAvailablePartition()
	if availablePartition == nil {
		return errors.New("no hay particiones disponibles")
	}

	// Crear la partición con los parámetros proporcionados
	availablePartition.CreatePartition(startPartition, sizeBytes, fdisk.typ, fdisk.fit, fdisk.name)

	// Colocar la partición en el MBR
	if availablePartition != nil {
		mbr.Mbr_partitions[indexPartition] = *availablePartition
	}

	// Serializar el MBR en el archivo binario
	err = mbr.Serialize(fdisk.path)
	if err != nil {
		fmt.Println("Error:", err)
	}

	return nil
}

func createExtendedPartition(fdisk *FDISK, sizeBytes int) error {
	// Crear una instancia de MBR
	var mbr structures.MBR

	// Deserializar la estructura MBR desde un archivo binario
	err := mbr.Deserialize(fdisk.path)
	if err != nil {
		fmt.Println("Error deserializando el MBR:", err)
		return err
	}

	// Veremos si existe una partición extendida en el disco
	extendedPartition, startExtended, indexExtended := mbr.GetExtended()
	fmt.Println("Particion" + string(startExtended) + string(indexExtended))

	if extendedPartition == nil {
		// Obtener la primera partición disponible
		availablePartition, startPartition, indexPartition := mbr.GetFirstAvailablePartition()
		if availablePartition == nil {
			return errors.New("No hay particiones disponibles")
		}

		// Crear la partición con los parámetros proporcionados
		availablePartition.CreatePartition(startPartition, sizeBytes, fdisk.typ, fdisk.fit, fdisk.name)

		// Colocar la partición en el MBR
		if availablePartition != nil {
			mbr.Mbr_partitions[indexPartition] = *availablePartition
		}

		// Serializar el MBR en el archivo binario
		err = mbr.Serialize(fdisk.path)
		if err != nil {
			fmt.Println("Error:", err)
		}

		err = createEBR(fdisk, int64(startPartition))
		if err != nil {
			fmt.Println("Error creating EBR:", err)
			return err
		}

	} else {
		return errors.New("Ya existe una partición extendida en el disco")
	}

	// Serializar el MBR en el archivo binario
	err = mbr.Serialize(fdisk.path)
	if err != nil {
		fmt.Println("Error:", err)
	}

	return nil
}

func createLogicPartition(fdisk *FDISK, sizeBytes int) error {
	// Crear una instancia de MBR
	var mbr structures.MBR

	// Deserializar la estructura MBR desde un archivo binario
	err := mbr.Deserialize(fdisk.path)
	if err != nil {
		fmt.Println("Error deserializando el MBR:", err)
		return err
	}

	// Veremos si existe una partición extendida en el disco
	extendedPartition, startExtended, indexExtended := mbr.GetExtended()
	fmt.Println("Particion" + string(startExtended) + string(indexExtended))

	offset := int64(startExtended)
	if extendedPartition != nil {
		for i := 0; i < 10; i++ {
			// Crear una instancia de MBR
			var ebr structures.EBR

			// Deserializar la estructura MBR desde un archivo binario
			err := ebr.Deserialize(fdisk.path, offset)
			if err != nil {
				fmt.Println("Error deserializando el MBR:", err)
				return err
			}

			if ebr.Ebr_next == -1 {
				inicio := ebr.Ebr_start + ebr.Ebr_size
				err = createEBR(fdisk, int64(inicio))
				if err != nil {
					fmt.Println("Error creating EBR:", err)
					return err
				}

				ebr.Ebr_next = inicio

				// Serializar el MBR en el archivo binario
				err = ebr.Serialize(fdisk.path, offset)
				if err != nil {
					fmt.Println("Error:", err)
				}

				break

			}
			offset = offset + int64(ebr.Ebr_next)

		}

	} else {
		return errors.New("No existe una partición extendida en el disco")
	}
	return nil
}

func createEBR(fdisk *FDISK, offset int64) error {
	sizeBytes, errr := utils.ConvertToBytes(fdisk.size, fdisk.unit)
	if errr != nil {
		fmt.Println("Error converting size:", errr)
		return errr
	}

	ebr := &structures.EBR{
		Ebr_mount: [1]byte{'9'},
		Ebr_fit:   [1]byte{'9'},
		Ebr_start: int32(offset) + 30,
		Ebr_size:  int32(sizeBytes),
		Ebr_next:  int32(-1),
		Ebr_name:  [16]byte{'0'},
	}

	// Serializar el MBR en el archivo
	err := ebr.Serialize(fdisk.path, offset)
	if err != nil {
		fmt.Println("Error:", err)
	}

	return nil
}
