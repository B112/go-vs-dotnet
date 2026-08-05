namespace DotnetNotes;

public class NoteStore {
    private readonly Lock _lock = new Lock();
    private readonly List<Note> _notes = [];
    private int _nextId = 1;

    public NoteStore() {
        Add(new CreateNoteRequest("Setup Balena fleet", "Configure fleet variables for CM5 devices", "work"));
        Add(new CreateNoteRequest("Buy groceries", "Milk, bread, cheese, coffee", "personal"));
        Add(new CreateNoteRequest("WireGuard VPN config", "Generate new peer keys for test devices", "work"));
    }

    public List<Note> All() {
        lock (_lock) {
            return [.. _notes];
        }
    }

    public Note? ByID(int id) {
        lock (_lock) {
            return _notes.FirstOrDefault(n => n.ID == id);
        }
    }

    public Note Add(CreateNoteRequest req) {
        lock (_lock) {
            var note = new Note(_nextId++, req.Title, req.Content, req.Category, DateTime.UtcNow);
            _notes.Add(note);
            return note;
        }
    }

    public bool Delete(int id) {
        lock (_lock) {
            return _notes.RemoveAll(n => n.ID == id) > 0;
        }
    }
}
