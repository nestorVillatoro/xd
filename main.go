package main

import (
	analyzer "P1_202200252/analyzer" // Importa el paquete "bufio" para operaciones de buffer de entrada/salida
	"P1_202200252/global"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type CommandRequest struct {
	Command string `json:"command"`
}

type CommandResponse struct {
	Output string `json:"output"`
	Error  string `json:"error,omitempty"`
}

type LoginRequest struct {
	IDParticion   string `json:"ID_particion"`
	NombreUsuario string `json:"nombre_usuario"`
	Password      string `json:"password"`
}

type LoginResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

type LogoutResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

func main() {

	http.HandleFunc("/api/execute", withCORS(ejecutarComando)) // Define el endpoint
	http.HandleFunc("/api/login", withCORS(controlLogin))      // Nuevo endpoint para login
	http.HandleFunc("/api/logout", withCORS(controlLogout))    // Nuevo endpoint para logout

	fmt.Println("Servidor escuchando en http://localhost:8080")
	err := http.ListenAndServe(":8080", nil) // Inicia el servidor
	if err != nil {
		fmt.Println("Error al iniciar el servidor:", err)
	}

}

// Middleware para manejar CORS
func withCORS(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Configurar los encabezados CORS
		w.Header().Set("Access-Control-Allow-Origin", "*") // Permite todos los orígenes
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == http.MethodOptions {
			// Responder a las solicitudes preflight con un código 204 (No Content)
			w.WriteHeader(http.StatusNoContent)
			return
		}

		// Llamar al manejador principal
		next(w, r)
	}
}

func ejecutarComando(w http.ResponseWriter, r *http.Request) {

	var req CommandRequest

	// Decodifica la solicitud JSON
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	input := req.Command
	var results []string

	lines := strings.Split(input, "\n")
	for _, line := range lines {

		if strings.HasPrefix(line, "#") {
			// Separar '#' del texto con un espacio
			line = "# " + strings.TrimPrefix(line, "#")
		}

		if line == "" {
			result := ""
			results = append(results, result)
			continue
		}
		// Llamar a la función Analyzer del paquete analyzer para analizar la línea
		result, err := analyzer.Analyzer(line)
		if err != nil {
			// Si hay un error, almacenar el mensaje de error en lugar del resultado
			result = fmt.Sprintf("Error: %s", err.Error())
		}

		// Acumular los resultados
		results = append(results, result)
	}

	// Crea la respuesta
	consola := strings.Join(results, "\n")
	// Preparar la respuesta
	res := CommandResponse{
		Output: consola,
	}

	// Devuelve la respuesta como JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)

}

func controlLogin(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	fmt.Println("ENTRAMOS")
	// Decodificar la solicitud JSON
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	mensaje, verificador := global.Verificar_login(req.NombreUsuario, req.Password, req.IDParticion)

	if verificador {
		// Aquí puedes implementar tu lógica real de validación de usuario y contraseña
		// Responder con éxito si las credenciales son correctas
		res := LoginResponse{
			Success: verificador,
			Message: mensaje,
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(res)
		fmt.Println("VALIDADO")
		return
	}

	// Responder con error si las credenciales son incorrectas
	res := LoginResponse{
		Success: verificador,
		Message: mensaje,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
	fmt.Println(mensaje)
}

// Manejador para el logout
func controlLogout(w http.ResponseWriter, r *http.Request) {
	// Aquí puedes limpiar cualquier sesión o token si estás usando uno en el futuro
	salida, _ := analyzer.Analyzer("logout")

	if salida == "Sesión cerrada exitosamente" {
		res := LogoutResponse{
			Success: true,
			Message: "Sesión cerrada exitosamente",
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(res)
		fmt.Println("")
		fmt.Println("Sesión cerrada")
		return
	}

	res := LogoutResponse{
		Success: true,
		Message: "Error al cerrar sesion",
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
	fmt.Println("")
	fmt.Println("Sesión no cerrada")

}
