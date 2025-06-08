## Usando curl desde la terminal:

<i>curl</i> es una herramienta de línea de comandos muy potente para realizar solicitudes HTTP. Es una manera rápida y sencilla de probar tus endpoints.

- Hacer una solicitud POST (por ejemplo, para /register):

```Bash
curl -X POST -H "Content-Type: application/json" -d '{"username": "testuser", "email": "test@example.com", "name": "Test User", "password": "securepassword"}' http://localhost:8080/register
```

- Hacer una solicitud POST para /login:

```Bash
curl -X POST -H "Content-Type: application/json" -d '{"username": "testuser", "password": "securepassword"}' http://localhost:8080/login
```

- Hacer una solicitud GET (por ejemplo, para /library, si requiere autenticación, necesitarás incluir el header adecuado, como una cookie):

```Bash
curl http://localhost:8080/library -H "Cookie: session_id=tu_sesion_id"
```

- Hacer una solicitud GET con un parámetro en la ruta (por ejemplo, /audio/:filename):

```Bash
curl http://localhost:8080/audio/song.mp3 -H "Cookie: session_id=tu_sesion_id" -o downloaded_song.mp3
```

La flag -o guarda la respuesta en un archivo.
