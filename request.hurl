# Ciclo de prueba de la API. Ejecutar con el comando de la línea siguiente.
# La base debe estar levantada y el servidor escuchando en localhost:8080.
# Cada ejecución genera un email distinto para que el test se pueda repetir.
# hurl --test --variable email="ana.$(date +%s%N)@example.com" request.hurl

POST http://localhost:8080/libros
Content-Type: application/json
{
	"titulo": "Ficciones",
	"autor": "Jorge Luis Borges",
	"fecha_publicacion": "1944-01-01",
	"genero": "Ficción"
}
HTTP/1.1 201
[Captures]
libro_id: jsonpath "$.id"
[Asserts]
jsonpath "$.titulo" == "Ficciones"

GET http://localhost:8080/libros
HTTP/1.1 200

GET http://localhost:8080/libros/{{libro_id}}
HTTP/1.1 200
[Asserts]
jsonpath "$.titulo" == "Ficciones"

PUT http://localhost:8080/libros/{{libro_id}}
Content-Type: application/json
{
	"titulo": "Ficciones (edición de prueba)",
	"autor": "Jorge Luis Borges",
	"fecha_publicacion": "1944-01-01",
	"genero": "Ficción"
}
HTTP/1.1 200
[Asserts]
jsonpath "$.mensaje" == "Libro actualizado correctamente"

GET http://localhost:8080/libros/{{libro_id}}
HTTP/1.1 200
[Asserts]
jsonpath "$.titulo" == "Ficciones (edición de prueba)"

GET http://localhost:8080/libros/no-es-un-id
HTTP/1.1 400

GET http://localhost:8080/libros/2147483647
HTTP/1.1 404

POST http://localhost:8080/libros
Content-Type: application/json
{
	"titulo": "JSON inválido"
}
HTTP/1.1 400

DELETE http://localhost:8080/libros/{{libro_id}}
HTTP/1.1 204

GET http://localhost:8080/libros/{{libro_id}}
HTTP/1.1 404

POST http://localhost:8080/libros
Content-Type: application/json
{
	"titulo": "Libro del Futuro",
	"autor": "Autor Futuro",
	"fecha_publicacion": "2099-01-01",
	"genero": "Ciencia Ficción"
}
HTTP/1.1 400
[Asserts]
jsonpath "$.error" == "La fecha de publicación no puede ser futura"
