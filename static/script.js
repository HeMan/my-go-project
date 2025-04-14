const apiBase = "/todos";

document.addEventListener("DOMContentLoaded", () => {
    const todoList = document.getElementById("todo-list");
    const todoSubjectInput = document.getElementById("todo-subject");
    const todoNotesInput = document.getElementById("todo-notes");
    const addTodoButton = document.getElementById("add-todo");

    // Fetch and display todos
    const fetchTodos = async () => {
        console.log("Fetching todos...");
        try {
            const response = await fetch(apiBase);
            if (!response.ok) {
                console.error("Failed to fetch todos:", response.statusText);
                return;
            }
            let todos = await response.json();

            // Sort todos: by due_date (closest first), then by completion status
            todos.sort((a, b) => {
                if (a.completed && !b.completed) return 1; // Move completed todos to the end
                if (!a.completed && b.completed) return -1; // Keep uncompleted todos first
                if (!a.due_date) return 1; // a has no due_date, move it after todos with due_date
                if (!b.due_date) return -1; // b has no due_date, move it after todos with due_date
                return new Date(a.due_date) - new Date(b.due_date); // Sort by closest due_date
            });

            todoList.innerHTML = "";
            todos.forEach(todo => {
                const li = document.createElement("li");
                li.className = todo.completed ? "completed" : "";
                li.dataset.id = todo.ID; // Set the data-id attribute

                // Add debug logging to check the notes structure
                console.log('Todo notes:', todo.notes);

                const notesSection = Array.isArray(todo.notes) && todo.notes.length > 0
                    ? `<div class="notes-section">
                        <small>Notes:</small>
                        <ul class="notes-list">
                            ${todo.notes.map(note => `<li>${note?.note || 'No content'} <button onclick="removeNote(${todo.ID}, '${note.ID}')">Remove</button></li>`).join('')}
                        </ul>
                      </div>`
                    : '';

                const dueDateSection = todo.due_date
                    ? `<div class="todo-due-date">
                        <small>Due: ${new Date(todo.due_date).toLocaleDateString()}</small>
                      </div>`
                    : '';

                li.innerHTML = `
                    <div class="todo-main">
                        <span class="todo-subject">${todo.subject}</span>
                        ${dueDateSection}
                        <div class="todo-actions">
                            <button onclick="toggleTodo(${todo.ID}, ${todo.completed})">
                                ${todo.completed ? "Unmark" : "Complete"}
                            </button>
                            <button onclick="deleteTodo(${todo.ID})">Delete</button>
                            <button onclick="editTodo(${todo.ID})">Edit</button>
                        </div>
                    </div>
                    ${notesSection}
                `;
                todoList.appendChild(li);
            });
        } catch (error) {
            console.error("Error fetching todos:", error);
        }
    };

    // Add a new todo
    addTodoButton.addEventListener("click", async () => {
        const subject = todoSubjectInput.value.trim();
        const noteText = todoNotesInput ? todoNotesInput.value.trim() : '';
        const dueDateInput = document.getElementById("todo-due-date"); // Get the due date input
        const dueDate = dueDateInput ? new Date(dueDateInput.value.trim()).toISOString() : null; // Format as ISO string

        if (!subject) return alert("Please enter a todo subject.");

        const todoData = {
            subject,
            completed: false,
            due_date: dueDate, // Include due_date in ISO format
            notes: noteText ? [{ content: noteText }] : []
        };

        await fetch(apiBase, {
            method: "POST",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify(todoData),
        });

        todoSubjectInput.value = "";
        if (todoNotesInput) todoNotesInput.value = "";
        if (dueDateInput) dueDateInput.value = ""; // Clear the due date input
        fetchTodos();
    });

    // Toggle todo completion
    window.toggleTodo = async (id, completed) => {
        await fetch(`${apiBase}/${id}`, {
            method: "PATCH",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify({ completed: !completed }),
        });
        fetchTodos();
    };

    // Delete a todo
    window.deleteTodo = async (id) => {
        try {
            const response = await fetch(`${apiBase}/${id}`, { method: "DELETE" });
            if (!response.ok) {
                console.error(`Failed to delete todo with id ${id}:`, response.statusText);
                return;
            }
            fetchTodos();
        } catch (error) {
            console.error(`Error deleting todo with id ${id}:`, error);
        }
    };

    // Edit a todo
    window.editTodo = (id) => {
        console.log("Editing todo with id:", id);
        const editModal = document.getElementById("edit-modal");
        const editSubjectInput = document.getElementById("edit-subject");
        const editDueDateInput = document.getElementById("edit-due-date");
        const editNoteInput = document.getElementById("edit-note");
        const saveEditButton = document.getElementById("save-edit");

        // Find the todo element using the correct dataset attribute
        const todo = Array.from(todoList.children).find(li => li.dataset.id == id);
        if (!todo) {
            console.error(`Todo with id ${id} not found.`);
            return;
        }

        // Populate modal with current todo data
        editSubjectInput.value = todo.querySelector(".todo-subject").textContent;
        editDueDateInput.value = todo.dataset.dueDate || ""; // Ensure dataset.dueDate exists
        editNoteInput.value = ""; // Clear the note input field for new notes
        editModal.dataset.id = id;

        // Display the modal
        editModal.style.display = "block";

        saveEditButton.onclick = async () => {
            const updatedTodo = {
                subject: editSubjectInput.value.trim(),
                due_date: editDueDateInput.value.trim() ? new Date(editDueDateInput.value.trim()).toISOString() : null, // Format as ISO string
                notes: editNoteInput.value.trim() ? [{ note: editNoteInput.value.trim() }] : []
            };

            await fetch(`${apiBase}/${editModal.dataset.id}`, {
                method: "PATCH",
                headers: { "Content-Type": "application/json" },
                body: JSON.stringify(updatedTodo),
            });

            editModal.style.display = "none";
            fetchTodos();
        };
    };

    // Remove a note
    window.removeNote = async (id, noteId) => {
        const response = await fetch(`${apiBase}/${id}/notes/${noteId}`, {
            method: "DELETE",
            headers: { "Content-Type": "application/json" },
        });

        if (response.ok) {
            fetchTodos();
        } else {
            console.error("Failed to remove note:", response.statusText);
        }
    };

    // Initial fetch
    fetchTodos();
});
