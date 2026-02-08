package main

// =============================================================================
// GO + HTMX TODO APP
// =============================================================================
// A modern server-side rendered web application using HTMX for interactivity.
// Templates are embedded into the binary using go:embed for portability.
//
// HTMX Core Concepts Used:
// - hx-post/hx-delete: Issue AJAX requests
// - hx-target: Specify where to swap response
// - hx-swap: Control how content is swapped (innerHTML, outerHTML)
// - hx-indicator: Show loading state during requests
// - hx-confirm: Confirm before destructive actions
// =============================================================================

import (
	"embed"              // Embed templates into binary
	"fmt"                // Formatting
	"html/template"      // HTML templating
	"log/slog"           // Structured logging
	"net/http"           // HTTP server
	"os"                 // OS operations
	"strconv"            // String conversion
	"sync"               // Mutex for thread safety
	"time"               // Time operations
)

// =============================================================================
// EMBEDDED TEMPLATES
// =============================================================================
// Templates are embedded at compile time, making the binary self-contained.

//go:embed templates/*.html
var templateFS embed.FS

// =============================================================================
// DATA MODEL
// =============================================================================

// Todo represents a single todo item.
type Todo struct {
	ID        int    // Unique identifier
	Title     string // Todo title/description
	Completed bool   // Whether the todo is completed
}

// =============================================================================
// APPLICATION STATE
// =============================================================================

// App holds the application state and dependencies.
type App struct {
	tmpl   *template.Template // Parsed HTML templates
	logger *slog.Logger       // Structured logger
	mu     sync.Mutex         // Mutex for thread-safe todo access
	todos  []Todo             // In-memory todo storage
	nextID int                // Next available ID
}

// NewApp creates and initializes the application.
func NewApp(logger *slog.Logger) (*App, error) {
	// Parse embedded templates
	tmpl, err := template.ParseFS(templateFS, "templates/*.html")
	if err != nil {
		return nil, fmt.Errorf("failed to parse templates: %w", err)
	}

	return &App{
		tmpl:   tmpl,
		logger: logger,
		todos: []Todo{
			{ID: 1, Title: "Learn Go", Completed: true},
			{ID: 2, Title: "Learn HTMX", Completed: false},
			{ID: 3, Title: "Build something awesome", Completed: false},
		},
		nextID: 4,
	}, nil
}

// =============================================================================
// HTTP ROUTES
// =============================================================================

// Routes returns the HTTP handler with all routes configured.
func (app *App) Routes() http.Handler {
	mux := http.NewServeMux()

	// Page routes
	mux.HandleFunc("GET /", app.handleIndex)

	// API routes (return HTML fragments for HTMX)
	mux.HandleFunc("POST /todos", app.handleAddTodo)
	mux.HandleFunc("POST /todos/{id}/toggle", app.handleToggleTodo)
	mux.HandleFunc("DELETE /todos/{id}", app.handleDeleteTodo)

	// Apply middleware
	return app.loggingMiddleware(mux)
}

// =============================================================================
// HANDLERS
// =============================================================================

// handleIndex renders the full page (initial load).
func (app *App) handleIndex(w http.ResponseWriter, r *http.Request) {
	// Only handle exact "/" path
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	app.mu.Lock()
	todos := make([]Todo, len(app.todos))
	copy(todos, app.todos)
	app.mu.Unlock()

	data := map[string]any{
		"Todos":      todos,
		"TotalCount": len(todos),
	}

	if err := app.tmpl.ExecuteTemplate(w, "index.html", data); err != nil {
		app.logger.Error("template error", slog.String("error", err.Error()))
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

// handleAddTodo adds a new todo and returns the updated list fragment.
// This demonstrates HTMX's partial page update pattern.
func (app *App) handleAddTodo(w http.ResponseWriter, r *http.Request) {
	title := r.FormValue("title")
	if title == "" {
		http.Error(w, "Title is required", http.StatusBadRequest)
		return
	}

	app.mu.Lock()
	newTodo := Todo{
		ID:        app.nextID,
		Title:     title,
		Completed: false,
	}
	app.nextID++
	app.todos = append(app.todos, newTodo)
	todos := make([]Todo, len(app.todos))
	copy(todos, app.todos)
	app.mu.Unlock()

	// Return only the todo-list fragment (HTMX will swap it into #todo-list)
	app.renderTodoList(w, todos)
}

// handleToggleTodo toggles the completion status of a todo.
func (app *App) handleToggleTodo(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	app.mu.Lock()
	for i := range app.todos {
		if app.todos[i].ID == id {
			app.todos[i].Completed = !app.todos[i].Completed
			break
		}
	}
	todos := make([]Todo, len(app.todos))
	copy(todos, app.todos)
	app.mu.Unlock()

	app.renderTodoList(w, todos)
}

// handleDeleteTodo removes a todo from the list.
func (app *App) handleDeleteTodo(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	app.mu.Lock()
	for i := range app.todos {
		if app.todos[i].ID == id {
			app.todos = append(app.todos[:i], app.todos[i+1:]...)
			break
		}
	}
	todos := make([]Todo, len(app.todos))
	copy(todos, app.todos)
	app.mu.Unlock()

	app.renderTodoList(w, todos)
}

// renderTodoList renders just the todo list fragment.
func (app *App) renderTodoList(w http.ResponseWriter, todos []Todo) {
	data := map[string]any{
		"Todos": todos,
	}
	if err := app.tmpl.ExecuteTemplate(w, "todo-list", data); err != nil {
		app.logger.Error("template error", slog.String("error", err.Error()))
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

// =============================================================================
// MIDDLEWARE
// =============================================================================

// loggingMiddleware logs each HTTP request.
func (app *App) loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		app.logger.Info("request",
			slog.String("method", r.Method),
			slog.String("path", r.URL.Path),
			slog.Duration("duration", time.Since(start)),
		)
	})
}

// =============================================================================
// MAIN
// =============================================================================

func main() {
	// Initialize structured logger
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))

	// Create application
	app, err := NewApp(logger)
	if err != nil {
		logger.Error("failed to create app", slog.String("error", err.Error()))
		os.Exit(1)
	}

	// Start server
	addr := ":8080"
	logger.Info("starting server", slog.String("addr", addr))
	logger.Info("open http://localhost:8080 in your browser")

	if err := http.ListenAndServe(addr, app.Routes()); err != nil {
		logger.Error("server error", slog.String("error", err.Error()))
		os.Exit(1)
	}
}
