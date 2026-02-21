package http

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/sun1tar/MIREA-TIP-Practice-18/tech-ip-sem2/tasks/internal/client/authclient"
	"github.com/sun1tar/MIREA-TIP-Practice-18/tech-ip-sem2/tasks/internal/service"
)

type TaskHandler struct {
	taskService *service.TaskService
	authClient  *authclient.Client
}

func NewTaskHandler(ts *service.TaskService, ac *authclient.Client) *TaskHandler {
	return &TaskHandler{
		taskService: ts,
		authClient:  ac,
	}
}

func (h *TaskHandler) verifyToken(w http.ResponseWriter, r *http.Request) bool {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		http.Error(w, `{"error":"missing authorization header"}`, http.StatusUnauthorized)
		return false
	}
	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		http.Error(w, `{"error":"invalid authorization header format"}`, http.StatusUnauthorized)
		return false
	}
	token := parts[1]

	valid, _, err := h.authClient.VerifyToken(r.Context(), token)
	if err != nil {
		http.Error(w, `{"error":"authentication service unavailable"}`, http.StatusServiceUnavailable)
		return false
	}
	if !valid {
		http.Error(w, `{"error":"invalid token"}`, http.StatusUnauthorized)
		return false
	}
	return true
}

type createTaskRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	DueDate     string `json:"due_date"`
}

type updateTaskRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	DueDate     string `json:"due_date"`
	Done        bool   `json:"done"`
}

type taskResponse struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	DueDate     string `json:"due_date,omitempty"`
	Done        bool   `json:"done"`
}

func toTaskResponse(t service.Task) taskResponse {
	return taskResponse{
		ID:          t.ID,
		Title:       t.Title,
		Description: t.Description,
		DueDate:     t.DueDate,
		Done:        t.Done,
	}
}

func (h *TaskHandler) CreateTask(w http.ResponseWriter, r *http.Request) {
	if !h.verifyToken(w, r) {
		return
	}

	var req createTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"invalid request body"}`, http.StatusBadRequest)
		return
	}
	if req.Title == "" {
		http.Error(w, `{"error":"title is required"}`, http.StatusBadRequest)
		return
	}

	task := service.Task{
		Title:       req.Title,
		Description: req.Description,
		DueDate:     req.DueDate,
		Done:        false,
	}
	created := h.taskService.Create(task)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(toTaskResponse(created))
}

func (h *TaskHandler) ListTasks(w http.ResponseWriter, r *http.Request) {
	if !h.verifyToken(w, r) {
		return
	}

	tasks := h.taskService.List()
	resp := make([]taskResponse, len(tasks))
	for i, t := range tasks {
		resp[i] = toTaskResponse(t)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func (h *TaskHandler) GetTask(w http.ResponseWriter, r *http.Request) {
	if !h.verifyToken(w, r) {
		return
	}

	id := r.PathValue("id")
	task, ok := h.taskService.Get(id)
	if !ok {
		http.Error(w, `{"error":"task not found"}`, http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(toTaskResponse(task))
}

func (h *TaskHandler) UpdateTask(w http.ResponseWriter, r *http.Request) {
	if !h.verifyToken(w, r) {
		return
	}

	id := r.PathValue("id")
	var req updateTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"invalid request body"}`, http.StatusBadRequest)
		return
	}

	updatedTask := service.Task{
		Title:       req.Title,
		Description: req.Description,
		DueDate:     req.DueDate,
		Done:        req.Done,
	}
	task, ok := h.taskService.Update(id, updatedTask)
	if !ok {
		http.Error(w, `{"error":"task not found"}`, http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(toTaskResponse(task))
}

func (h *TaskHandler) DeleteTask(w http.ResponseWriter, r *http.Request) {
	if !h.verifyToken(w, r) {
		return
	}

	id := r.PathValue("id")
	ok := h.taskService.Delete(id)
	if !ok {
		http.Error(w, `{"error":"task not found"}`, http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
