{{range . -}}
pub mod {{ .TableName }};
{{ end }}

// Every Table Must Have Primary Keys to Be Included In This Output
// Primiay Keys must be one column (no compostion types yet)
// Primiay Keys can be 1) Auto Increment 2) Other self Inserted

// Implemention is simple NOT many features is suported in Rust version:
// Keep mysql data types in int, bigint, text, varchar, bool, blob
// No signed integer is supported
// For now Primary key should only be numbers
// Not fully ORM is supported: limited to CRUD on rows + Indexes querys