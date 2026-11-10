// ==========================================================================
// Configuración y Endpoints de la API
// ==========================================================================
// Si se accede directamente desde el servidor Go (:8080) se usa '/libros',
// si se abre desde otro servidor local (ej. Live Server) o file://, apunta a localhost:8080
const API_URL = (window.location.protocol === 'file:' || (window.location.port !== '8080' && window.location.port !== ''))
    ? 'http://localhost:8080/libros'
    : '/libros';

// Elementos del DOM
const formularioLibro = document.getElementById('formulario-libro');
const inputTitulo = document.getElementById('titulo');
const inputAutor = document.getElementById('autor');
const inputFecha = document.getElementById('fecha_publicacion');
const inputGenero = document.getElementById('genero');
const contenedorLibros = document.getElementById('contenedor-libros');
const mensajeEstado = document.getElementById('mensaje-estado');
const btnGuardar = document.getElementById('btn-guardar');

// ==========================================================================
// Funciones Auxiliares
// ==========================================================================

/**
 * Escapa caracteres especiales para evitar inyecciones HTML en el DOM.
 */
function escapeHTML(str) {
    if (!str) return '';
    return str
        .replace(/&/g, '&amp;')
        .replace(/</g, '&lt;')
        .replace(/>/g, '&gt;')
        .replace(/"/g, '&quot;')
        .replace(/'/g, '&#039;');
}

/**
 * Formatea una fecha ISO (o YYYY-MM-DD) a formato legible YYYY-MM-DD.
 */
function formatearFecha(fechaStr) {
    if (!fechaStr) return 'Sin fecha';
    return fechaStr.split('T')[0];
}

/**
 * Muestra un mensaje temporal de éxito o error en la interfaz.
 */
let temporizadorMensaje = null;
function mostrarMensaje(texto, tipo = 'exito') {
    if (temporizadorMensaje) {
        clearTimeout(temporizadorMensaje);
    }
    mensajeEstado.textContent = texto;
    mensajeEstado.className = `mensaje-estado ${tipo}`;

    temporizadorMensaje = setTimeout(() => {
        mensajeEstado.textContent = '';
        mensajeEstado.className = 'mensaje-estado';
    }, 4000);
}

// ==========================================================================
// Paso 2: Obtener y Renderizar Entidades (GET /libros)
// ==========================================================================

/**
 * Realiza la petición GET /libros y renderiza los elementos en el DOM.
 */
async function obtenerLibros() {
    try {
        const respuesta = await fetch(API_URL);

        if (!respuesta.ok) {
            throw new Error(`Error en el servidor: ${respuesta.status} ${respuesta.statusText}`);
        }

        const libros = await respuesta.json();
        renderizarLibros(libros);
    } catch (error) {
        console.error('Error al obtener libros:', error);
        contenedorLibros.innerHTML = `
            <div class="mensaje-vacio">
                <p> No se pudo conectar con el servidor.</p>
                <small>Asegúrate de que el backend Go esté corriendo en el puerto 8080.</small>
            </div>
        `;
    }
}

/**
 * Genera dinámicamente el HTML para la lista de libros.
 */
function renderizarLibros(libros) {
    // Si no hay libros, mostrar mensaje descriptivo
    if (!Array.isArray(libros) || libros.length === 0) {
        contenedorLibros.innerHTML = `
            <div class="mensaje-vacio">
                <p>No hay libros registrados aún.</p>
                <small>¡Completa el formulario de la izquierda para registrar el primero!</small>
            </div>
        `;
        return;
    }

    // Construcción dinámica de las tarjetas
    const fragmento = document.createDocumentFragment();

    libros.forEach((libro) => {
        const tarjeta = document.createElement('article');
        tarjeta.className = 'tarjeta-libro';
        tarjeta.dataset.id = libro.id;

        tarjeta.innerHTML = `
            <div class="info-libro">
                <h3 class="titulo-libro">${escapeHTML(libro.titulo)}</h3>
                <p class="autor-libro">Por: ${escapeHTML(libro.autor)}</p>
                <div class="detalles-libro">
                    <span class="badge-genero">${escapeHTML(libro.genero)}</span>
                    <span> Publicado: ${escapeHTML(formatearFecha(libro.fecha_publicacion))}</span>
                </div>
            </div>
            <div class="acciones-libro">
                <!-- Botón con atributo data-id para identificar la entidad a borrar -->
                <button type="button" class="btn btn-peligro btn-eliminar" data-id="${libro.id}">
                    Eliminar
                </button>
            </div>
        `;

        fragmento.appendChild(tarjeta);
    });

    contenedorLibros.innerHTML = '';
    contenedorLibros.appendChild(fragmento);
}

// ==========================================================================
// Paso 2: Creación de Entidad (POST /libros)
// ==========================================================================

formularioLibro.addEventListener('submit', async (evento) => {
    evento.preventDefault();

    const titulo = inputTitulo.value.trim();
    const autor = inputAutor.value.trim();
    const fechaPublicacion = inputFecha.value.trim();
    const genero = inputGenero.value.trim();

    // Validar campos obligatorios en el cliente
    if (!titulo || !autor || !fechaPublicacion || !genero) {
        mostrarMensaje('Por favor, completa todos los campos del formulario.', 'error');
        return;
    }

    const nuevoLibro = {
        titulo: titulo,
        autor: autor,
        fecha_publicacion: fechaPublicacion,
        genero: genero
    };

    try {
        btnGuardar.disabled = true;
        btnGuardar.textContent = 'Guardando...';

        const respuesta = await fetch(API_URL, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify(nuevoLibro)
        });

        if (!respuesta.ok) {
            let errorMensaje = 'No se pudo crear el libro.';
            try {
                const errorData = await respuesta.json();
                if (errorData.error) errorMensaje = errorData.error;
            } catch (_) {}
            throw new Error(errorMensaje);
        }

        // Éxito: limpiar formulario, notificar y refrescar lista
        formularioLibro.reset();
        mostrarMensaje(`¡"${nuevoLibro.titulo}" agregado con éxito!`, 'exito');
        await obtenerLibros();

    } catch (error) {
        console.error('Error al crear libro:', error);
        mostrarMensaje(error.message, 'error');
    } finally {
        btnGuardar.disabled = false;
        btnGuardar.textContent = 'Agregar Libro';
    }
});

// ==========================================================================
// Paso 3: Borrado con Delegación de Eventos (DELETE /libros/{id})
// ==========================================================================

contenedorLibros.addEventListener('click', async (evento) => {
    // Delegación de eventos: detectar si el clic fue en un botón de eliminar
    const botonEliminar = evento.target.closest('.btn-eliminar');
    if (!botonEliminar) return;

    const id = botonEliminar.dataset.id;
    if (!id) return;

    const confirmacion = confirm(`¿Estás seguro de que deseas eliminar este libro?`);
    if (!confirmacion) return;

    try {
        botonEliminar.disabled = true;
        botonEliminar.textContent = 'Borrando...';

        const respuesta = await fetch(`${API_URL}/${id}`, {
            method: 'DELETE'
        });

        if (!respuesta.ok && respuesta.status !== 204) {
            let errorMensaje = 'No se pudo eliminar el libro.';
            try {
                const errorData = await respuesta.json();
                if (errorData.error) errorMensaje = errorData.error;
            } catch (_) {}
            throw new Error(errorMensaje);
        }

        mostrarMensaje('Libro eliminado correctamente.', 'exito');

        // Refrescar la vista con la lista actualizada
        await obtenerLibros();

    } catch (error) {
        console.error('Error al eliminar libro:', error);
        mostrarMensaje(error.message, 'error');
        botonEliminar.disabled = false;
        botonEliminar.textContent = 'Eliminar';
    }
});

// ==========================================================================
// Inicialización al cargar el documento
// ==========================================================================
document.addEventListener('DOMContentLoaded', () => {
    obtenerLibros();
});
