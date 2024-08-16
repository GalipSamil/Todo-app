document.addEventListener("DOMContentLoaded", () => {
  const form = document.getElementById("todo-form");

  form.addEventListener("submit", async (e) => {
    e.preventDefault(); // Formun varsayılan submit olayını engeller

    const title = document.getElementById("todo-title").value;
    if (title.trim() === "") {
      alert("Please enter a todo title.");
      return;
    }

    const response = await fetch("/todos", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({ title: title }), // Görev başlığını JSON formatında gönderir
    });

    if (response.ok) {
      document.getElementById("todo-title").value = ""; // Input'u temizle
      fetchTodos(); // Listeyi güncelle
    } else {
      console.error("Failed to add the todo", response.statusText);
    }
  });

  const fetchTodos = async () => {
    const response = await fetch("/todos");
    const todos = await response.json();
    renderTodos(todos); // Görevleri frontend'de göster
  };

  const renderTodos = (todos) => {
    const todoList = document.getElementById("todo-list");
    todoList.innerHTML = ""; // Listeyi temizler

    todos.forEach((todo) => {
      const li = document.createElement("li");
      li.textContent = todo.title;
      if (todo.completed) {
        li.classList.add("completed");
      }

      const completeBtn = document.createElement("button");
      completeBtn.textContent = "Complete";
      completeBtn.onclick = () => completeTodo(todo.id);

      const deleteBtn = document.createElement("button");
      deleteBtn.textContent = "Delete";
      deleteBtn.onclick = () => deleteTodo(todo.id);

      li.appendChild(completeBtn);
      li.appendChild(deleteBtn);
      todoList.appendChild(li); // Görevleri listeye ekler
    });
  };

  const completeTodo = async (id) => {
    const response = await fetch(`/todos/${id}/complete`, { method: "POST" });
    if (response.ok) {
      fetchTodos(); // Görev tamamlandıktan sonra listeyi güncelle
    } else {
      console.error("Failed to complete the todo", response.statusText);
    }
  };

  const deleteTodo = async (id) => {
    const response = await fetch(`/todos/${id}`, { method: "DELETE" });
    if (response.ok) {
      fetchTodos(); // Görev silindikten sonra listeyi güncelle
    } else {
      console.error("Failed to delete the todo", response.statusText);
    }
  };

  fetchTodos(); // Sayfa yüklendiğinde mevcut todo listesini çek
});
