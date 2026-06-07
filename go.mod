module merry-box

// Go 1.22+ is required: ServeMux gained method+path patterns
// (e.g. "GET /files/{name}") and r.PathValue in this version.
go 1.24.5
