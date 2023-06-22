import React, { useEffect, useState } from "react";
import axios from "axios";

const api = axios.create({
  baseURL: "http://localhost:8080",
});

type Note = {
  id: number;
  title: string;
  content: string;
};

const App: React.FC = () => {
  const [notes, setNotes] = useState<Note[]>([]);

  const [title, setTitle] = useState("");
  const [content, setContent] = useState("");

  const [editNoteId, setEditNoteId] = useState<number | null>(null);
  const [editTitle, setEditTitle] = useState("");
  const [editContent, setEditContent] = useState("");

  useEffect(() => {
    fetchNotes();
  }, []);

  const fetchNotes = async () => {
    try {
      const response = await api.get("/");
      setNotes(response.data.data ? response.data.data : []);
    } catch (error) {
      console.error("Error fetching notes:", error);
    }
  };

  const createNote = async () => {
    try {
      await api.post("/", { title, content });
      fetchNotes();
      setTitle("");
      setContent("");
    } catch (error) {
      console.error("Error creating note:", error);
    }
  };

  const updateNote = async (event: React.FormEvent, id: number) => {
    event.preventDefault();
    try {
      console.log(editTitle, editContent);
      await api.put(`/${id}`, { title: editTitle, content: editContent });
      fetchNotes();
      setEditNoteId(null);
      setEditTitle("");
      setEditContent("");
    } catch (error) {
      console.error(`Error updating note with ID ${id}:`, error);
    }
  };

  const deleteNote = async (id: number) => {
    try {
      await api.delete(`/${id}`);
      fetchNotes();
    } catch (error) {
      console.error(`Error deleting note with ID ${id}:`, error);
    }
  };

  const handleEditNote = (note: Note) => {
    setEditNoteId(note.id);
    setEditTitle(note.title);
    setEditContent(note.content);
  };

  return (
    <div
      style={{
        display: "flex",
        justifyContent: "center",
        alignItems: "start",
        width: "100vw",
        height: "100vh",
      }}
    >
      <div style={{ margin: "0 auto", maxWidth: "800px" }}>
        <h1>Notes</h1>
        <div>
          <input
            type="text"
            placeholder="Title"
            value={title}
            onChange={(e) => setTitle(e.target.value)}
          />
          <input
            type="text"
            placeholder="Content"
            value={content}
            onChange={(e) => setContent(e.target.value)}
          />
          <button onClick={createNote}>Create Note</button>
        </div>
        <ul style={{ maxHeight: "70vh", overflowY: "scroll" }}>
          {notes.map((note) => (
            <li key={note.id}>
              {editNoteId === note.id ? (
                <form onSubmit={(event) => updateNote(event, note.id)}>
                  Title:
                  <input
                    type="text"
                    placeholder="Title"
                    value={editTitle}
                    onChange={(e) => setEditTitle(e.target.value)}
                  />
                  Content:
                  <input
                    type="text"
                    placeholder="Content"
                    value={editContent}
                    onChange={(e) => setEditContent(e.target.value)}
                  />
                  <button type="submit">Update</button>
                  <button onClick={() => setEditNoteId(null)}>Cancel</button>
                </form>
              ) : (
                <>
                  <h5>Id: {note.id}</h5>
                  <h3>{note.title}</h3>
                  <p>{note.content}</p>
                  <button onClick={() => handleEditNote(note)}>Edit</button>
                  <button onClick={() => deleteNote(note.id)}>Delete</button>
                </>
              )}
            </li>
          ))}
        </ul>
      </div>
    </div>
  );
};

export default App;
