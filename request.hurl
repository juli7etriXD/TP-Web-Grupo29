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

POST http://localhost:8080/usuarios
Content-Type: application/json
{
	"nombre": "Ana",
	"apellido": "Pérez",
	"email": "{{email}}",
	"password": "clave123",
	"fecha_nacimiento": "1998-04-15"
}
HTTP/1.1 201
[Captures]
usuario_id: jsonpath "$.id"
[Asserts]
jsonpath "$.email" == "{{email}}"

POST http://localhost:8080/usuarios
Content-Type: application/json
{
	"nombre": "Ana Copia",
	"apellido": "Pérez",
	"email": "{{email}}",
	"password": "clave123",
	"fecha_nacimiento": "1998-04-15"
}
HTTP/1.1 409
[Asserts]
jsonpath "$.error" == "El email ya se encuentra registrado"

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

GET http://localhost:8080/usuarios
HTTP/1.1 200

GET http://localhost:8080/usuarios/{{usuario_id}}
HTTP/1.1 200
[Asserts]
jsonpath "$.email" == "{{email}}"

PUT http://localhost:8080/usuarios/{{usuario_id}}
Content-Type: application/json
{
	"nombre": "Ana María",
	"apellido": "Pérez",
	"email": "{{email}}",
	"password": "clave456"
}
HTTP/1.1 200
[Asserts]
jsonpath "$.mensaje" == "Usuario actualizado correctamente"

GET http://localhost:8080/usuarios/{{usuario_id}}
HTTP/1.1 200
[Asserts]
jsonpath "$.nombre" == "Ana María"
jsonpath "$.email" == "{{email}}"

POST http://localhost:8080/usuarios/{{usuario_id}}/libros
Content-Type: application/json
{
	"libro_id": {{libro_id}},
	"leido": false
}
HTTP/1.1 201
[Asserts]
jsonpath "$.leido" == false

POST http://localhost:8080/usuarios/{{usuario_id}}/libros
Content-Type: application/json
{
	"libro_id": {{libro_id}},
	"leido": false
}
HTTP/1.1 409
[Asserts]
jsonpath "$.error" == "El libro ya está asociado al usuario"

POST http://localhost:8080/usuarios/2147483647/libros
Content-Type: application/json
{
	"libro_id": {{libro_id}},
	"leido": false
}
HTTP/1.1 404
[Asserts]
jsonpath "$.error" == "El usuario o el libro especificado no existe"

PATCH http://localhost:8080/usuarios/{{usuario_id}}/libros/{{libro_id}}
Content-Type: application/json
{
	"leido": true
}
HTTP/1.1 200
[Asserts]
jsonpath "$.mensaje" == "Estado de lectura actualizado correctamente"

GET http://localhost:8080/usuarios/{{usuario_id}}/libros/leidos
HTTP/1.1 200
[Asserts]
jsonpath "$[0].titulo" == "Ficciones (edición de prueba)"

GET http://localhost:8080/libros/{{libro_id}}/usuarios
HTTP/1.1 200
[Asserts]
jsonpath "$[0].email" == "{{email}}"

DELETE http://localhost:8080/usuarios/{{usuario_id}}/libros/{{libro_id}}
HTTP/1.1 204

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

DELETE http://localhost:8080/usuarios/{{usuario_id}}
HTTP/1.1 204

DELETE http://localhost:8080/libros/{{libro_id}}
HTTP/1.1 204

GET http://localhost:8080/libros/{{libro_id}}
HTTP/1.1 404

GET http://localhost:8080/usuarios/{{usuario_id}}
HTTP/1.1 404

GET http://localhost:8080/usuarios/{{usuario_id}}/libros/leidos
HTTP/1.1 404
[Asserts]
jsonpath "$.error" == "Usuario no encontrado"

GET http://localhost:8080/libros/{{libro_id}}/usuarios
HTTP/1.1 404
[Asserts]
jsonpath "$.error" == "Libro no encontrado"

POST http://localhost:8080/usuarios
Content-Type: application/json
{
	"nombre": "Carlos",
	"apellido": "Gómez",
	"email": "correo-invalido",
	"password": "clave123",
	"fecha_nacimiento": "1995-05-10"
}
HTTP/1.1 400
[Asserts]
jsonpath "$.error" == "Formato de email inválido"

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
