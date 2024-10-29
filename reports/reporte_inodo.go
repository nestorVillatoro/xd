package reports

import (
	structures "P1_202200252/structures"
	utils "P1_202200252/utils"
	"fmt"
	"os"
	"os/exec"
	"time"
)

// ReportInode genera un reporte de un inodo y lo guarda en la ruta especificada
func ReportInode(superblock *structures.SuperBlock, diskPath string, path string) error {
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

	// Iterar sobre cada inodo
	for i := int32(0); i < superblock.S_inodes_count; i++ {
		inode := &structures.Inode{}
		// Deserializar el inodo
		err := inode.Deserialize(diskPath, int64(superblock.S_inode_start+(i*superblock.S_inode_size)))
		if err != nil {
			return err
		}

		// Convertir tiempos a string
		atime := time.Unix(int64(inode.I_atime), 0).Format(time.RFC3339)
		ctime := time.Unix(int64(inode.I_ctime), 0).Format(time.RFC3339)
		mtime := time.Unix(int64(inode.I_mtime), 0).Format(time.RFC3339)

		// Definir el contenido DOT para el inodo actual con estilo mejorado
		dotContent += fmt.Sprintf(`inode%d [label=<
            <table border="1" cellborder="2" cellspacing="0" cellpadding="5" bgcolor="lightgray">
                <tr>
                    <td colspan="2" bgcolor="black" fontcolor="white">
                        <b><font color="white">REPORTE INODO %d</font></b>
                    </td>
                </tr>
                <tr><td><b>i_uid</b></td><td>%d</td></tr>
                <tr><td><b>i_gid</b></td><td>%d</td></tr>
                <tr><td><b>i_size</b></td><td>%d bytes</td></tr>
                <tr><td><b>i_atime</b></td><td>%s</td></tr>
                <tr><td><b>i_ctime</b></td><td>%s</td></tr>
                <tr><td><b>i_mtime</b></td><td>%s</td></tr>
                <tr><td><b>i_type</b></td><td>%c</td></tr>
                <tr><td><b>i_perm</b></td><td>%s</td></tr>
                <tr><td colspan="2" bgcolor="gray" fontcolor="white"><b><font color="white">BLOQUES DIRECTOS</font></b></td></tr>
            `, i, i, inode.I_uid, inode.I_gid, inode.I_size, atime, ctime, mtime, rune(inode.I_type[0]), string(inode.I_perm[:]))

		// Agregar los bloques directos a la tabla con un estilo mejorado
		for j, block := range inode.I_block {
			if j > 11 {
				break
			}
			dotContent += fmt.Sprintf("<tr><td><b>Bloque %d</b></td><td>%d</td></tr>", j+1, block)
		}

		// Agregar los bloques indirectos con formato
		dotContent += fmt.Sprintf(`
                <tr><td colspan="2" bgcolor="gray" fontcolor="white"><b><font color="white">BLOQUE INDIRECTO</font></b></td></tr>
                <tr><td><b>Bloque 13</b></td><td>%d</td></tr>
                <tr><td colspan="2" bgcolor="gray" fontcolor="white"><b><font color="white">BLOQUE INDIRECTO DOBLE</font></b></td></tr>
                <tr><td><b>Bloque 14</b></td><td>%d</td></tr>
                <tr><td colspan="2" bgcolor="gray" fontcolor="white"><b><font color="white">BLOQUE INDIRECTO TRIPLE</font></b></td></tr>
                <tr><td><b>Bloque 15</b></td><td>%d</td></tr>
            </table>>];
        `, inode.I_block[12], inode.I_block[13], inode.I_block[14])

		// Agregar enlace al siguiente inodo si no es el último
		if i < superblock.S_inodes_count-1 {
			dotContent += fmt.Sprintf("inode%d -> inode%d;\n", i, i+1)
		}
	}

	// Cerrar el contenido DOT
	dotContent += "}"

	// Crear el archivo DOT
	dotFile, err := os.Create(dotFileName)
	if err != nil {
		return err
	}
	defer dotFile.Close()

	// Escribir el contenido DOT en el archivo
	_, err = dotFile.WriteString(dotContent)
	if err != nil {
		return err
	}

	// Generar la imagen con Graphviz
	cmd := exec.Command("dot", "-Tpng", dotFileName, "-o", outputImage)
	err = cmd.Run()
	if err != nil {
		return err
	}

	fmt.Println("Imagen de los inodos generada:", outputImage)
	return nil
}
