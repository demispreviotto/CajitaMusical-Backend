# CajitaMusical.

## PostgreSQL.

1. Iniciar el servicio de PostgreSQL:
   Si no está corriendo automáticamente, puedes iniciarlo con Homebrew:
   `brew services start postgresql`

Si ya está corriendo, te dirá Service 'postgresql' already started.

2. Verificar el estado:
   `brew services list`

Busca postgresql en la lista y asegúrate de que su estado sea started.

3. Acceder a la línea de comandos de PostgreSQL (psql - opcional):
   Si quieres conectarte a la base de datos para verificar tablas o datos:
   `psql postgres`

Esto te conectará a la base de datos por defecto postgres. Luego, puedes conectarte a tu base de datos de la aplicación con \c your_database_name;.

## Go Backend.

`go run cmd/server/main.go`

`go mod tidy`
`go clean -cache`
`go run cmd/server/main.go`

## Svelte Frontend

`npm run dev -- --open`
