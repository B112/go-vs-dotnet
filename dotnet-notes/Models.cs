// ReSharper disable ClassNeverInstantiated.Global
// ReSharper disable NotAccessedPositionalProperty.Global
namespace DotnetNotes;

public record Note(int ID, string Title, string Content, string Category, DateTime CreatedAt);

public record CreateNoteRequest(string Title, string Content = "", string Category = "general");

public record ImportRequest(string FilePath);
