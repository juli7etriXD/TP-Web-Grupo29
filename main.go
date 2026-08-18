package main

import (
	"fmt"
	"net/http"
)

func main() {
	// 1. Define el directorio que contiene los archivos estáticos.
	staticDir := "./static"

	// 2. Crea un manejador (handler) de servidor de archivos.
	// http.Dir convierte la ruta del directorio en un sistema de archivos HTTP.
	// http.FileServer crea un manejador que sirve archivos desde ese sistema.
	// ¡Automáticamente sirve index.html para directorios!
	fileServer := http.FileServer(http.Dir(staticDir))
	// 3. Registra el manejador para que atienda todas las peticiones ("/").
	// Usamos http.Handle porque fileServer es un http.Handler.
	http.Handle("/", fileServer)

	// 4. Define el puerto y muestra un mensaje.
	port := ":8080"
	fmt.Printf("Servidor ESTÁTICO escuchando en http://localhost%s\n", port)
	fmt.Printf("Sirviendo archivos desde: %s\n", staticDir)

	// 5. Inicia el servidor.
	err := http.ListenAndServe(port, nil)
	if err != nil {
 		fmt.Printf("Error al iniciar el servidor: %s\n", err)
  }
}
