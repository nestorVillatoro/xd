package structures

var listaEBR []EBR

func AgregarElemento(elemento EBR) {

	listaEBR = append(listaEBR, elemento)
}

func ObtenerLista() []EBR {
	return listaEBR
}
