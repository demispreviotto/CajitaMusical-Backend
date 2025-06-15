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

`go run cmd/api/main.go`

`go mod tidy`
`go clean -cache`
`go run cmd/api/main.go`

## Svelte Frontend

`npm run dev -- --open`

## Understanding the Layers

Think of your application in layers, each with a specific responsibility:

Models (Entities): These represent the core data structures of your application, directly mapping to your database tables. They should be "pure" and contain no business logic or database access logic.

Example: models.User, models.Song, models.Session, models.Authentication.
Location: internal/models/
Database (Repository/Persistence Layer): This layer is responsible for interacting with the database. It knows how to save, retrieve, update, and delete data. It should expose interfaces that the service layer can depend on, hiding the underlying database technology (PostgreSQL, MySQL, etc.).

Example: db.SongDBer (interface), db.songDB (implementation), db.Connect(), db.DB (the GORM instance).
Location: internal/db/
Services (Business Logic Layer): This is where your application's core business rules and logic reside. It orchestrates interactions between the database layer, external services, and potentially other business concerns. It operates on your models and uses the database interfaces.

Example: services.SongServicer (interface), services.songService (implementation), services.ScanMusicLibrary(), services.GetSongFilePath().
Location: internal/services/
Handlers (API/Presentation Layer): This layer handles incoming HTTP requests, validates input, calls the appropriate service methods, and formats the response for the client. It should be thin and primarily concern itself with HTTP.

Example: handlers.songHandler, handlers.GetLibrary(), handlers.ServeAudio().
Location: internal/handlers/
DTOs (Data Transfer Objects): These are structures specifically designed for transferring data between different layers, especially between the API layer (handlers) and the client. They might represent subsets of your models, or combined data from multiple models, formatted for a specific API response or request.

Example: song.CreateSongInput, song.SongResponse, song.ListSongsResponse.
Location: internal/dtos/ (or internal/api/ if you prefer grouping API-related structures)
