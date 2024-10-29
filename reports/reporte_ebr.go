package reports

import (
	structures "P1_202200252/structures"
	utils "P1_202200252/utils"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// ReportMBR genera un reporte del MBR y lo guarda en la ruta especificada
func ReportEBR(path string) error {
	// Crear las carpetas padre si no existen
	err := utils.CreateParentDirs(path)
	if err != nil {
		return err
	}

	// Obtener el nombre base del archivo sin la extensión
	dotFileName, outputImage := utils.GetFileNames(path)

	// Definir el contenido DOT con una tabla estilizada
	dotContent := `digraph G {
        node [shape=plaintext]
        tabla [label=<
            <table border="1" cellborder="2" cellspacing="0" cellpadding="5" bgcolor="lightgray">
                <tr>
                    <td colspan="2" bgcolor="black" fontcolor="white">
                        <b><font color="white">REPORTE EBR</font></b>
                    </td>
                </tr>
`
	lista := structures.ObtenerLista()
	// Agregar las particiones a la tabla con estilo
	for i, part := range lista {
		// Convertir Part_name a string y eliminar los caracteres nulos
		partName := strings.TrimRight(string(part.Ebr_name[:]), "\x00")
		// Convertir Part_status, Part_type y Part_fit a char
		partStatus := rune(part.Ebr_mount[0])
		partFit := rune(part.Ebr_fit[0])

		// Agregar la partición a la tabla con un estilo mejorado
		dotContent += fmt.Sprintf(`
				<tr><td colspan="2" bgcolor="gray" fontcolor="white"><b><font color="white">PARTICIÓN %d</font></b></td></tr>
				<tr><td><b>Estado</b></td><td>%d</td></tr>
				<tr><td><b>Fit</b></td><td>%c</td></tr>
				<tr><td><b>Inicio</b></td><td>%d</td></tr>
				<tr><td><b>Tamaño</b></td><td>%d bytes</td></tr>
				<tr><td><b>Nombre</b></td><td>%s</td></tr>
			`, i+1, partStatus, partFit, part.Ebr_start, part.Ebr_size, partName)
	}

	// Cerrar la tabla y el contenido DOT
	dotContent += "</table>>] }"

	// Guardar el contenido DOT en un archivo
	file, err := os.Create(dotFileName)
	if err != nil {
		return fmt.Errorf("error al crear el archivo: %v", err)
	}
	defer file.Close()

	_, err = file.WriteString(dotContent)
	if err != nil {
		return fmt.Errorf("error al escribir en el archivo: %v", err)
	}

	// Ejecutar el comando Graphviz para generar la imagen
	cmd := exec.Command("dot", "-Tpng", dotFileName, "-o", outputImage)
	err = cmd.Run()
	if err != nil {
		return fmt.Errorf("error al ejecutar el comando Graphviz: %v", err)
	}

	fmt.Println("Imagen de la tabla generada:", outputImage)
	return nil
}
