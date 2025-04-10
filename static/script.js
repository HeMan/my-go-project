const apiBase = "/todos";

document.addEventListener("DOMContentLoaded", () => {
    const todoList = document.getElementById("todo-list");
    const todoSubjectInput = document.getElementById("todo-subject");
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
            const todos = await response.json();
            todoList.innerHTML = "";
            todos.forEach(todo => {
                const li = document.createElement("li");
                li.className = todo.completed ? "completed" : "";
                li.innerHTML = `
                    <span>${todo.subject}</span>
                    <div>
                        <button onclick="toggleTodo(${todo.ID}, ${todo.completed})">
                            ${todo.completed ? "Unmark" : "Complete"}
                        </button>
                        <button onclick="deleteTodo(${todo.ID})">Delete</button>
                    </div>
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
        if (!subject) return alert("Please enter a todo subject.");
        await fetch(apiBase, {
            method: "POST",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify({ subject, completed: false }),
        });
        todoSubjectInput.value = "";
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

    // Initial fetch
    fetchTodos();
});
