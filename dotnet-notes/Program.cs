using DotnetNotes;
using Microsoft.AspNetCore.Http.Json;
using System.Text.Json;
var builder = WebApplication.CreateBuilder(args);
builder.Services.AddSingleton<NoteStore>();
builder.Services.Configure<JsonOptions>(o =>
                                            o.SerializerOptions.PropertyNamingPolicy = JsonNamingPolicy.CamelCase);
var app = builder.Build();

var startTime = DateTime.UtcNow;
var api = app.MapGroup("/api");

// GET /api/notes?category=work
api.MapGet("/notes", (NoteStore store, string? category) => {
    var notes = store.All();
    if (!string.IsNullOrEmpty(category)) notes = notes.Where(n => n.Category.Equals(category, StringComparison.OrdinalIgnoreCase)).ToList();
    return Results.Ok(notes);
});

// GET /api/notes/{id}
api.MapGet("/notes/{id:int}", (NoteStore store, int id) =>
               store.ByID(id) is { } note
                   ? Results.Ok(note)
                   : Results.NotFound(new {
                       error = "not found",
                   }));

// POST /api/notes
api.MapPost("/notes", (NoteStore store, CreateNoteRequest req) => {
    if (string.IsNullOrWhiteSpace(req.Title))
        return Results.BadRequest(new {
            error = "title required",
        });
    var note = store.Add(req with {
        Category = string.IsNullOrEmpty(req.Category) ? "general" : req.Category,
    });
    return Results.Created($"/api/notes/{note.ID}", note);
});

// DELETE /api/notes/{id}
api.MapDelete("/notes/{id:int}", (NoteStore store, int id) =>
                  store.Delete(id)
                      ? Results.NoContent()
                      : Results.NotFound(new {
                          error = "not found",
                      }));

// GET /api/stats
api.MapGet("/stats", (NoteStore store) => {
    var notes = store.All();
    return Results.Ok(new {
        total = notes.Count,
        byCategory = notes.GroupBy(n => n.Category).ToDictionary(g => g.Key, g => g.Count()),
        uptime = (DateTime.UtcNow - startTime).ToString(@"hh\:mm\:ss"),
    });
});

// POST /api/notes/import — error handling example
// C# uses try/catch: wrap risky code, catch specific exception types.
api.MapPost("/notes/import", async (NoteStore store, HttpRequest request) => {
    ImportRequest? body;
    try {
        body = await request.ReadFromJsonAsync<ImportRequest>();
        if (body is null || string.IsNullOrWhiteSpace(body.FilePath))
            return Results.BadRequest(new {
                error = "filePath required",
            });
    } catch (JsonException ex) {
        return Results.BadRequest(new {
            error = $"invalid request body: {ex.Message}",
        });
    }

    string fileContent;
    try {
        fileContent = await File.ReadAllTextAsync(body.FilePath);
    } catch (FileNotFoundException) {
        return Results.NotFound(new {
            error = $"file not found: {body.FilePath}",
        });
    } catch (UnauthorizedAccessException) {
        return Results.StatusCode(403);
    } catch (IOException) {
        return Results.StatusCode(500);
    }

    List<CreateNoteRequest>? notes;
    try {
        notes = JsonSerializer.Deserialize<List<CreateNoteRequest>>(fileContent,
                                                                    new JsonSerializerOptions {
                                                                        PropertyNameCaseInsensitive = true,
                                                                    });
        if (notes is null)
            return Results.BadRequest(new {
                error = "file contained null",
            });
    } catch (JsonException ex) {
        return Results.BadRequest(new {
            error = $"invalid JSON in file: {ex.Message}",
        });
    }

    var imported = 0;
    var errors = new List<string>();
    for(var i = 0; i < notes.Count; i++) {
        if (string.IsNullOrWhiteSpace(notes[i].Title)) {
            errors.Add($"note {i}: title required");
            continue;
        }
        store.Add(notes[i] with {
            Category = string.IsNullOrEmpty(notes[i].Category) ? "general" : notes[i].Category,
        });
        imported++;
    }

    return Results.Ok(new {
        imported,
        errors,
    });
});

// Serve UI from wwwroot/ — edit wwwroot/index.html like any normal file.
// Files in wwwroot/ are bundled into the publish output automatically.
app.UseStaticFiles();
app.MapFallbackToFile("index.html");

app.Urls.Add("http://localhost:5081");
Console.WriteLine("Dotnet Notes running on http://localhost:5081");
app.Run();
