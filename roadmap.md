# Roadmap — DevOps Hiring Challenge (File Storage Service)

## Tổng quan đề bài

Đề có **2 phần lớn**, nối tiếp nhau:

- **Requirement 1 — Application:** một service lưu/lấy/xóa file qua HTTP/REST API. Chứng minh khả năng *viết code*.
- **Requirement 2 — Deployment:** đưa app lên AWS, mọi thứ quản lý bằng code (IaC + CI/CD). Chứng minh khả năng *vận hành*.

## Tech stack đã chốt

| Thành phần | Lựa chọn |
|---|---|
| Backend | Golang (`net/http` chuẩn, không cần router ngoài từ Go 1.22+) |
| Metadata DB | PostgreSQL |
| Blob storage | Local filesystem (Phase đầu) → S3 (sau) |
| Container | Docker |
| IaC | Terraform (AWS provider) |
| CI/CD | GitHub Actions |
| Compute | EC2 (trước) → ECS (sau) |
| Frontend | ReactJS (optional, cuối cùng) |
| Repo | 2 repo tách biệt: `app` và `infra` |

## Triết lý làm việc

- **Functional-first:** chạy được trước, tối ưu sau.
- **Local → Container → Cloud:** dựng local, containerize, rồi mới deploy.
- **EC2 → ECS:** đơn giản trước, nâng cấp sau.
- **Manual → IaC:** làm tay để *hiểu* AWS, rồi Terraform hóa để *củng cố* (kiểu làm lab). Không đầu tư sâu vào bản làm tay vì sẽ bị vứt.
- **Deadline:** ~1 tháng, linh hoạt. Ưu tiên hiểu sâu hơn tốc độ.

---

## Phần A — Application milestone ladder (Requirement 1)

Mỗi nấc là một thứ **chạy được, demo được**; nấc sau thêm đúng một khái niệm.

| Milestone | Nội dung | Khái niệm mới |
|---|---|---|
| **M0** | HTTP server có `GET /health` trả 200 | Dựng server |
| **M1** | Lưu/lấy/xóa trong RAM (`map[string][]byte`) | Hình dạng 3 API, concurrency (mutex) |
| **M2** | Persist bytes xuống đĩa local | Lộ ra vấn đề trùng tên |
| **M3** | Metadata vào PostgreSQL (tách 2 bảng `files`/`blobs`) | Persistence, query, search |
| **M4** | Streaming upload + size limit (413) | Xử file lớn không OOM |
| **M5** | Content-addressable storage (băm SHA-256, lưu blob theo hash) | CAS |
| **M6** | Deduplication (check hash, `ref_count`) | Dedup, tiết kiệm dung lượng |
| **M7** | Trừu tượng `BlobStore` interface (`LocalFSBlobStore`) | Dependency inversion, mở đường S3 |
| **M8** | Làm cứng concurrency (unique constraint, atomic ref_count, cleanup temp, GC orphan) | Race conditions |
| **M9** | Hoàn thiện API (search, pagination, delete idempotency, error body chuẩn) | Production polish |

**Hết M9:** có một file-storage service hoàn chỉnh chạy local trong Docker.

**Vị trí hiện tại:** vừa xong blueprint (API contract + schema). Đang bắt tay **M1**.

---

## Phần B — Deployment phases (Requirement 2)

| Phase | Nội dung | Output |
|---|---|---|
| **P0** | Nền tảng: ôn HTTP, REST, Go `net/http`, PostgreSQL cơ bản, Git | Repo skeleton |
| **P1** | App lõi chạy local (tương ứng M1–M3) | API CRUD chạy local |
| **P2** | Bonus app: dedup + large file (M4–M6) | App có CAS/dedup |
| **P3** | Đóng gói: Dockerfile (multi-stage) + docker-compose | `docker compose up` chạy cả stack |
| **P4** | Deploy tay lên EC2: SG, SSH, chạy Docker, RDS PostgreSQL | API qua IP public; *hiểu* AWS |
| **P5** | IaC bằng Terraform: codify VPC/Subnet/SG/EC2/RDS | `apply` dựng / `destroy` xóa sạch |
| **P6** | CI/CD bằng GitHub Actions: build+test, build image → ECR, deploy | Push code → tự lên môi trường |
| **P7** | Nâng cấp EC2 → ECS (Task Def, Service, Cluster, Fargate) | App chạy trên ECS |
| **P8** | Chọn lọc bonus tiers (xem dưới) | Điểm cộng |
| **P9** | Frontend ReactJS (optional) | Trang upload/list/delete |
| **P10** | Hoàn thiện deliverables (xuyên suốt) | README, diagram, cost, teardown |

### Bonus tiers — ưu tiên theo điểm/độ dễ

1. **Secrets** (SSM Parameter Store) — dễ, ăn điểm bảo mật.
2. **Observability** (CloudWatch logs + structured logging) — dễ vừa.
3. **S3** làm blob store thay filesystem (viết `S3BlobStore`, lúc này mới cần AWS SDK Go) — hợp với dedup.
4. **IAM least-privilege** (bỏ `*` trong Action/Resource) — thể hiện hiểu bảo mật.
5. **HTTPS** (ALB + ACM) — trung bình.
6. **Autoscaling** (demo trigger under load) — khó nhất, để cuối.

---

## Deliverables (chốt cuối)

- [ ] Public GitHub/GitLab repo (app + infra)
- [ ] README: architecture decision, run locally, deploy, clean up
- [ ] Architecture diagram (draw.io / Excalidraw / vẽ tay)
- [ ] Estimated monthly cost nếu chạy 24/7 (EC2/RDS/ALB...)
- [ ] Teardown instructions verify bước xóa (dựa trên `terraform destroy`)

## Lưu ý cost & teardown

- Làm tay xong phải **xóa hết bằng tay** trước khi dựng lại bằng Terraform.
- Orphaned resources tốn tiền thật — reviewer sẽ chạy code trong tài khoản của họ.