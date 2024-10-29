package reports

import (
	structures "P1_202200252/structures"
	utils "P1_202200252/utils"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"
)

// ReportMBR genera un reporte del MBR y lo guarda en la ruta especificada
func ReportMBR(mbr *structures.MBR, path string, mountedDiskPath string) error {
	// Crear las carpetas padre si no existen
	err := utils.CreateParentDirs(path)
	if err != nil {
		return err
	}

	// Obtener el nombre base del archivo sin la extensión
	dotFileName, outputImage := utils.GetFileNames(path)

	// Iniciar el contenido DOT
	dotContent := `digraph G {
        node [shape=plaintext]
    `

	// Definir el contenido DOT con una tabla de MBR con estilo
	dotContent += fmt.Sprintf(`
        mbr [label=<
            <table border="1" cellborder="2" cellspacing="0" cellpadding="5" bgcolor="lightgray">
                <tr>
                    <td colspan="2" bgcolor="black"><b><font color="white">REPORTE MBR</font></b></td>
                </tr>
                <tr><td bgcolor="gold"><b>mbr_tamano</b></td><td>%d</td></tr>
                <tr><td bgcolor="gold"><b>mrb_fecha_creacion</b></td><td>%s</td></tr>
                <tr><td bgcolor="gold"><b>mbr_disk_signature</b></td><td>%d</td></tr>
            `, mbr.Mbr_size, time.Unix(int64(mbr.Mbr_creation_date), 0).Format(time.RFC3339), mbr.Mbr_disk_signature)

	// Agregar las particiones a la tabla
	for i, part := range mbr.Mbr_partitions {
		partName := strings.TrimRight(string(part.Part_name[:]), "\x00")
		partStatus := rune(part.Part_status[0])
		partType := rune(part.Part_type[0])
		partFit := rune(part.Part_fit[0])

		// Agregar la partición a la tabla con estilo
		dotContent += fmt.Sprintf(`
                <tr>
                    <td colspan="2" bgcolor="brown"><b>PARTICIÓN %d</b></td>
                </tr>
                <tr><td bgcolor="red"><b>part_status</b></td><td>%c</td></tr>
                <tr><td bgcolor="red"><b>part_type</b></td><td>%c</td></tr>
                <tr><td bgcolor="red"><b>part_fit</b></td><td>%c</td></tr>
                <tr><td bgcolor="red"><b>part_start</b></td><td>%d</td></tr>
                <tr><td bgcolor="red"><b>part_size</b></td><td>%d</td></tr>
                <tr><td bgcolor="red"><b>part_name</b></td><td>%s</td></tr>
            `, i+1, partStatus, partType, partFit, part.Part_start, part.Part_size, partName)

		// Agregar particiones lógicas si la partición es extendida
		if string(part.Part_type[0]) == "E" {
			offset := int64(part.Part_start)
			for j := 0; j < 100; j++ {
				var ebr structures.EBR
				err := ebr.Deserialize(mountedDiskPath, offset)
				if err != nil {
					fmt.Println("Error deserializando el EBR:", err)
					return err
				}

				logicPartName := strings.TrimRight(string(ebr.Ebr_name[:]), "\x00")
				ebrStatus := rune(ebr.Ebr_mount[0])
				ebrFit := rune(part.Part_fit[0])

				// Agregar la partición lógica con estilo
				dotContent += fmt.Sprintf(`
                    <tr>
                        <td colspan="2" bgcolor="brown"><b>PARTICIÓN LÓGICA %d</b></td>
                    </tr>
                    <tr><td bgcolor="orange"><b>EBR_mount</b></td><td>%c</td></tr>
                    <tr><td bgcolor="orange"><b>EBR_fit</b></td><td>%c</td></tr>
                    <tr><td bgcolor="orange"><b>EBR_start</b></td><td>%d</td></tr>
                    <tr><td bgcolor="orange"><b>EBR_size</b></td><td>%d</td></tr>
                    <tr><td bgcolor="orange"><b>EBR_next</b></td><td>%d</td></tr>
                    <tr><td bgcolor="orange"><b>EBR_name</b></td><td>%s</td></tr>
                `, j+1, ebrStatus, ebrFit, ebr.Ebr_start, ebr.Ebr_size, ebr.Ebr_next, logicPartName)

				if ebr.Ebr_next == -1 {
					break
				}
				offset = offset + 30 + int64(ebr.Ebr_size)
			}
		}
	}

	// Cerrar la tabla y el contenido DOT
	dotContent += "</table>>]; }"

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

	fmt.Println("Imagen del MBR generada:", outputImage)
	return nil
}
