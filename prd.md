##**Goal**: A purely local, keyboard-driven (Vim + Arrows) task management board built with Elm-architecture UI and a robust SQLite data layer.

##**Constraints**: 

* Fixed columns layout (Todo, Doing, Done).

* Maximum board width enforced (centered on screen) to prevent layout tearing on ultrawide displays.

* Data layer must be fully abstracted via interfaces.

##SQLite Schema (ERD)

```SQL
CREATE TABLE tasks (
    id TEXT PRIMARY KEY,           -- UUID
    title TEXT NOT NULL,           -- e.g., "Migrate RabbitMQ events to Redis stream"
    description TEXT,              -- e.g., "Ensure OpenTelemetry traces are preserved"
    status TEXT NOT NULL,          -- 'TODO', 'DOING', 'DONE'
    created_at DATETIME NOT NULL,
    updated_at DATETIME NOT NULL
);
```