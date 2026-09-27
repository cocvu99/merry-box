  # merry-box

  A minimal file storage service exposing a REST API to upload, retrieve,
  and delete files — conceptually a simplified S3/Dropbox-style backend.

  Files with identical content are deduplicated via content-addressable
  storage (SHA-256 hashing), similar to how Git stores blob objects.

  Infrastructure (AWS, Terraform, CI/CD) lives in a separate repo:
  [merry-infra](<link>).

  > Status: work in progress — rebuilding after first review pass.


