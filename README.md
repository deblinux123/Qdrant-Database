# Qdrant-Database

# 🚀 یادگیری Qdrant با Golang — برنامه ۱۴ روزه

یک مسیر عملی ۱۴ روزه برای یادگیری **Qdrant Vector Database** با زبان **Golang**.

هدف این دوره این است که در پایان بتوانیم یک سیستم واقعی مبتنی بر **Vector Search / Semantic Search / RAG** با Go و Qdrant بسازیم.

---

# 📚 فهرست مطالب

* [Qdrant چیست؟](#qdrant-چیست)
* [چرا Qdrant؟](#چرا-qdrant)
* [مفاهیم اصلی](#مفاهیم-اصلی-qdrant)
* [معماری کلی](#معماری-کلی)
* [پیش‌نیازها](#پیشنیازها)
* [برنامه ۱۴ روزه](#برنامه-۱۴-روزه)
* [مفاهیم مهم Qdrant](#مفاهیم-مهم-qdrant)
* [Embedding چیست؟](#embedding-چیست)
* [Collection چیست؟](#collection-چیست)
* [Point چیست؟](#point-چیست)
* [Payload چیست؟](#payload-چیست)
* [Vector Search](#vector-search)
* [Filtering](#filtering)
* [Indexing](#indexing)
* [Distance Metrics](#distance-metrics)
* [Top-K Search](#top-k-search)
* [Semantic Search](#semantic-search)
* [Hybrid Search](#hybrid-search)
* [Qdrant در RAG](#qdrant-در-rag)
* [ساختار پروژه](#ساختار-پروژه)
* [پروژه نهایی](#پروژه-نهایی)

---

# Qdrant چیست؟

**Qdrant** یک Vector Database متن‌باز است که برای ذخیره، مدیریت و جستجوی Vectorها طراحی شده است.

در پروژه‌های هوش مصنوعی معمولاً داده‌هایی مانند:

* متن
* تصویر
* صوت
* Document
* Product
* مقاله
* پیام
* رزومه

را به Vector تبدیل می‌کنیم.

سپس این Vectorها را در Qdrant ذخیره می‌کنیم تا بتوانیم داده‌هایی را پیدا کنیم که از نظر معنایی به Query کاربر نزدیک هستند.

مثلاً:

```text
Query:

چطور می‌توانم در Golang یک API بسازم؟
```

ممکن است در دیتابیس متنی داشته باشیم:

```text
ساخت REST API با زبان Go
```

حتی اگر کلمات دقیقاً یکسان نباشند، Vector Search می‌تواند ارتباط معنایی آن‌ها را پیدا کند.

---

# چرا Qdrant؟

Qdrant برای پروژه‌های AI و RAG امکانات مهمی دارد:

* Vector Search
* Semantic Search
* Filtering
* Payload
* Metadata
* Multiple Vector
* Sparse Vector
* Dense Vector
* Quantization
* Indexing
* Collection Management
* Similarity Search
* Recommendation
* Hybrid Search
* REST API
* gRPC
* SDK
* مناسب برای RAG

Qdrant همچنین می‌تواند به‌صورت:

```text
Docker
Linux
Cloud
Self-hosted
```

استفاده شود.

---

# معماری کلی

یک سیستم ساده AI معمولاً به شکل زیر است:

```text
User
 │
 ▼
Go Application
 │
 ├── Generate Embedding
 │
 ▼
Embedding Model
 │
 ▼
Vector
 │
 ▼
Qdrant
 │
 ├── Vector
 ├── Payload
 └── Metadata
 │
 ▼
Similarity Search
 │
 ▼
Relevant Documents
 │
 ▼
LLM
 │
 ▼
Final Answer
```

در یک سیستم RAG:

```text
User Question
      │
      ▼
   Embedding
      │
      ▼
    Qdrant
      │
      ▼
Relevant Chunks
      │
      ▼
     LLM
      │
      ▼
    Answer
```

---

# پیش‌نیازها

قبل از شروع بهتر است با موارد زیر آشنا باشید:

## Golang

```text
Variables
Functions
Struct
Interface
Error Handling
JSON
HTTP
Goroutine
Context
Packages
Modules
```

## مفاهیم AI

لازم نیست Machine Learning را حرفه‌ای بدانید، اما باید با این مفاهیم آشنا شوید:

```text
Embedding
Vector
Similarity
Cosine Similarity
Semantic Search
RAG
LLM
```

---

# نصب Qdrant

ساده‌ترین روش استفاده از Docker است:

```bash
docker run -p 6333:6333 -p 6334:6334 \
  qdrant/qdrant
```

بعد Qdrant روی این آدرس در دسترس است:

```text
http://localhost:6333
```

---

# بررسی Qdrant

```bash
curl http://localhost:6333
```

---

# برنامه ۱۴ روزه

---

# 🟢 Day 01 — Vector Database چیست؟

## هدف

در روز اول باید بفهمیم چرا Vector Database به وجود آمده است.

موضوعات:

* Database چیست؟
* Vector چیست؟
* Vector Database چیست؟
* تفاوت SQL Database و Vector Database
* Semantic Search چیست؟
* Keyword Search چیست؟
* Embedding چیست؟

مثال:

```text
"برنامه‌نویسی Go"

↓ Embedding

[0.12, -0.42, 0.91, ...]
```

### تمرین

یک متن ساده انتخاب کنید و مفهوم تبدیل Text به Vector را بررسی کنید.

---

# 🟢 Day 02 — نصب و راه‌اندازی Qdrant

## هدف

راه‌اندازی Qdrant و آشنایی با API.

با Docker:

```bash
docker run -p 6333:6333 -p 6334:6334 \
  qdrant/qdrant
```

بررسی:

```bash
curl http://localhost:6333
```

موضوعات:

* Docker
* REST API
* Qdrant Dashboard
* REST
* gRPC

### تمرین

Qdrant را اجرا کنید و وضعیت آن را بررسی کنید.

---

# 🟢 Day 03 — Collection

## هدف

یادگیری Collection.

Collection در Qdrant تقریباً مشابه یک container برای Vectorها است.

مثلاً:

```text
papers
```

یا:

```text
documents
```

یا:

```text
products
```

هر Collection تنظیمات Vector خودش را دارد.

مثلاً:

```text
Vector Size: 768

Distance: Cosine
```

### موضوعات

* Collection
* Vector Size
* Distance
* Cosine
* Dot Product
* Euclidean

### تمرین

یک Collection به نام:

```text
documents
```

بسازید.

---

# 🟢 Day 04 — Point

## هدف

یادگیری Point.

هر داده‌ای که در Qdrant ذخیره می‌کنیم معمولاً یک Point است.

مثال:

```text
Point
 ├── ID
 ├── Vector
 └── Payload
```

مثلاً:

```json
{
  "id": 1,
  "vector": [0.12, 0.44, 0.91],
  "payload": {
    "title": "Golang Tutorial",
    "category": "programming"
  }
}
```

### تمرین

حداقل ۱۰ Point در Qdrant ذخیره کنید.

---

# 🟢 Day 05 — Golang + Qdrant

## هدف

امروز Qdrant را از طریق Go کنترل می‌کنیم.

ساختار ساده:

```text
Go
 │
 ▼
Qdrant Client
 │
 ▼
Qdrant
```

موضوعات:

* Go Module
* Qdrant Client
* Context
* HTTP
* Error Handling

ساخت پروژه:

```bash
mkdir qdrant-go
cd qdrant-go

go mod init qdrant-go
```

---

# 🟢 Day 06 — Insert / Upsert

## هدف

یادگیری ذخیره Vector.

مفهوم مهم:

```text
Upsert
```

یعنی:

```text
Insert
+
Update
```

ساختار داده:

```text
ID
Vector
Payload
```

مثال:

```text
ID: 1001

Vector:
[0.12, 0.22, 0.91, ...]

Payload:
{
    title: "Learning Go",
    author: "John"
}
```

### تمرین

یک برنامه Go بسازید که:

```text
10 Documents
      ↓
Embedding
      ↓
Qdrant
```

را ذخیره کند.

---

# 🟢 Day 07 — Payload

## هدف

Payload یکی از مهم‌ترین بخش‌های Qdrant است.

Payload اطلاعات اضافی مربوط به Vector است.

مثلاً:

```json
{
  "title": "Golang Tutorial",
  "author": "John",
  "category": "programming",
  "language": "en",
  "year": 2026
}
```

Vector برای similarity search استفاده می‌شود.

Payload برای:

```text
Metadata
Filtering
Filtering Search Results
```

استفاده می‌شود.

---

# 🟢 Day 08 — Vector Search

## هدف

امروز مهم‌ترین قابلیت Qdrant را یاد می‌گیریم:

```text
Similarity Search
```

فرض کنید Query داریم:

```text
How to write REST API in Go?
```

آن را به Vector تبدیل می‌کنیم:

```text
Query
 ↓
Embedding
 ↓
Vector
```

سپس:

```text
Vector
 ↓
Qdrant
 ↓
Top K Similar Vectors
```

مثلاً:

```text
1. Build REST API with Go
2. Golang HTTP Server
3. Fiber Web Framework
4. Go Backend Tutorial
5. Building APIs
```

---

# 🟢 Day 09 — Filtering

## هدف

ترکیب Vector Search با Filter.

مثلاً فقط مقالات:

```text
language = "fa"
```

یا:

```text
category = "AI"
```

یا:

```text
year >= 2025
```

را جستجو کنیم.

مثلاً:

```text
Semantic Search
+
language = fa
+
category = AI
```

این یکی از مهم‌ترین قابلیت‌های Qdrant برای سیستم‌های واقعی است.

---

# 🟢 Day 10 — Indexing

## هدف

یادگیری Index.

وقتی دیتابیس کوچک است، Search ساده است.

اما اگر داشته باشیم:

```text
100K vectors
1M vectors
10M vectors
```

Performance اهمیت زیادی پیدا می‌کند.

موضوعات:

* HNSW
* Payload Index
* Vector Index
* Search Performance
* Recall
* Latency

---

# 🟢 Day 11 — Embedding

## هدف

یادگیری کامل ارتباط:

```text
Text
 ↓
Embedding Model
 ↓
Vector
 ↓
Qdrant
```

مدل Embedding می‌تواند مثلاً:

```text
EmbeddingGemma
BGE
E5
Sentence Transformers
```

باشد.

نکته مهم:

### Vector Dimension

اگر مدل Embedding خروجی:

```text
768 dimensions
```

تولید کند، Collection باید با:

```text
size = 768
```

ساخته شود.

مثلاً:

```text
Embedding Model
       │
       │ 768
       ▼
Qdrant Collection
       │
       └── size: 768
```

---

# 🟢 Day 12 — Qdrant + RAG

## هدف

ساخت اولین RAG واقعی.

Architecture:

```text
User
 │
 ▼
Question
 │
 ▼
Embedding Model
 │
 ▼
Qdrant Search
 │
 ▼
Top K Documents
 │
 ▼
Context
 │
 ▼
LLM
 │
 ▼
Answer
```

مثلاً:

```text
User:

شرایط استفاده از پژوهش‌یار چیست؟
```

سیستم:

```text
Question
   ↓
Embedding
   ↓
Qdrant
   ↓
Relevant Documents
   ↓
LLM
   ↓
Persian Answer
```

---

# 🟢 Day 13 — Production Architecture

## هدف

امروز پروژه را به شکل Production نزدیک می‌کنیم.

Architecture:

```text
             ┌─────────────┐
             │   Client    │
             └──────┬──────┘
                    │
                    ▼
             ┌─────────────┐
             │ Go / Fiber  │
             └──────┬──────┘
                    │
          ┌─────────┴─────────┐
          ▼                   ▼
   Embedding Model           Redis
          │
          ▼
       Qdrant
          │
          ▼
       Retrieved
       Documents
          │
          ▼
          LLM
```

موضوعات:

* Docker
* Environment Variables
* Connection Pooling
* Context
* Logging
* Error Handling
* Health Check
* Qdrant Backup
* Performance
* Security

---

# 🟢 Day 14 — Final Project

## 🎯 پروژه نهایی

ساخت یک:

# Persian Semantic Search + RAG API

با:

```text
Golang
Fiber
Qdrant
Embedding Model
Ollama
Redis
```

Architecture:

```text
                 User
                   │
                   ▼
              Fiber API
                   │
          ┌────────┴────────┐
          │                 │
          ▼                 ▼
       Redis             Embedding
                            │
                            ▼
                         Qdrant
                            │
                            ▼
                     Relevant Chunks
                            │
                            ▼
                          Ollama
                            │
                            ▼
                       Final Answer
```

---

# 🧠 مفاهیم مهم Qdrant

در این بخش مهم‌ترین مفاهیمی که باید واقعاً بلد باشید را مرور می‌کنیم.

---

# 1. Vector

Vector یک آرایه عددی است.

مثلاً:

```text
[0.12, 0.43, -0.21, 0.87]
```

Embedding Model متن را به Vector تبدیل می‌کند.

مثلاً:

```text
"Go programming"
```

ممکن است تبدیل شود به:

```text
[0.12, 0.43, -0.21, ...]
```

---

# 2. Embedding

Embedding نمایش عددی یک داده است.

مثلاً:

```text
Text
 ↓
Embedding Model
 ↓
Vector
```

Embedding باعث می‌شود بتوانیم شباهت معنایی داده‌ها را محاسبه کنیم.

---

# 3. Collection

Collection محل نگهداری Vectorها است.

مثلاً:

```text
papers
```

می‌تواند شامل میلیون‌ها Vector باشد.

هر Collection معمولاً تنظیماتی مانند:

```text
Vector Size
Distance Metric
HNSW
Quantization
```

دارد.

---

# 4. Point

Point واحد اصلی داده در Qdrant است.

ساختار کلی:

```text
Point
├── ID
├── Vector
└── Payload
```

مثال:

```json
{
  "id": 123,
  "vector": [0.12, 0.42, 0.87],
  "payload": {
    "title": "AI Research",
    "language": "fa"
  }
}
```

---

# 5. Payload

Payload اطلاعات اضافی مربوط به Point است.

مثلاً:

```json
{
  "title": "Introduction to RAG",
  "author": "Farid",
  "language": "fa",
  "category": "AI",
  "year": 2026
}
```

Payload برای Filter بسیار مهم است.

---

# 6. Distance Metric

Qdrant باید بداند Vectorها چگونه با هم مقایسه شوند.

مهم‌ترین Metricها:

```text
Cosine
Dot
Euclidean
Manhattan
```

---

# Cosine Similarity

برای بسیاری از پروژه‌های NLP و Embedding گزینه رایجی است.

مفهوم آن بیشتر روی:

```text
Direction
```

Vector تمرکز دارد.

مثلاً:

```text
Vector A
   ↗

Vector B
   ↗
```

هرچه جهت آن‌ها مشابه‌تر باشد، Similarity بیشتر است.

---

# 7. Top-K

فرض کنید کاربر یک Query ارسال کرده است.

می‌خواهیم:

```text
Top 5
```

نتیجه مشابه را دریافت کنیم.

یعنی:

```text
limit = 5
```

نتیجه:

```text
1 → Score 0.94
2 → Score 0.91
3 → Score 0.88
4 → Score 0.84
5 → Score 0.81
```

---

# 8. Score

Score میزان شباهت Query و Vector را نشان می‌دهد.

مثلاً:

```text
Document A → 0.95
Document B → 0.83
Document C → 0.61
```

معمولاً نتیجه با Score بالاتر مرتبط‌تر است.

اما تفسیر دقیق Score به Distance Metric و مدل Embedding بستگی دارد.

---

# 9. Filtering

Filter اجازه می‌دهد Search را محدود کنیم.

مثلاً:

```text
category = AI
```

یا:

```text
language = fa
```

یا:

```text
year >= 2025
```

ترکیب بسیار قدرتمند:

```text
Vector Search
+
Payload Filter
```

است.

---

# 10. HNSW

یکی از مهم‌ترین مفاهیم Performance در Vector Search است.

HNSW مخفف:

```text
Hierarchical Navigable Small World
```

است.

به زبان ساده، یک ساختار graph-based برای پیدا کردن Vectorهای نزدیک است.

به جای اینکه تمام Vectorها را یکی‌یکی بررسی کنیم:

```text
1
2
3
4
5
...
1,000,000
```

از ساختار Graph برای پیدا کردن سریع‌تر Neighborها استفاده می‌شود.

---

# 11. Quantization

وقتی تعداد Vectorها بسیار زیاد می‌شود، مصرف RAM و Storage اهمیت پیدا می‌کند.

Quantization می‌تواند حجم Vectorها را کاهش دهد.

مثلاً:

```text
Full Precision
      ↓
Quantization
      ↓
Less Memory
```

اما باید بین:

```text
Memory
Speed
Recall
```

تعادل برقرار کنیم.

---

# 12. Dense Vector

Vectorهای معمولی Embedding معمولاً Dense هستند.

مثلاً:

```text
[0.12, 0.44, 0.87, 0.21, ...]
```

بسیاری از مدل‌های Embedding از Dense Vector استفاده می‌کنند.

---

# 13. Sparse Vector

در Sparse Vector بیشتر مقادیر صفر هستند.

مثلاً:

```text
[0, 0, 0.8, 0, 0, 0, 0.4, 0]
```

Sparse Vector برای بعضی سیستم‌های Search بسیار مفید است.

---

# 14. Hybrid Search

Hybrid Search ترکیبی از روش‌های مختلف Search است.

مثلاً:

```text
Dense Search
+
Sparse Search
```

Dense Search:

```text
Semantic Meaning
```

Sparse Search:

```text
Keyword / Exact Matching
```

ترکیب این دو می‌تواند Search قدرتمندتری ایجاد کند.

---

# Qdrant در RAG

یکی از مهم‌ترین کاربردهای Qdrant، ساخت RAG است.

## مرحله اول — Documents

```text
PDF
DOCX
TXT
Web Pages
Database
```

---

## مرحله دوم — Chunking

Document به قسمت‌های کوچک تقسیم می‌شود:

```text
Document
 │
 ├── Chunk 1
 ├── Chunk 2
 ├── Chunk 3
 └── Chunk 4
```

---

## مرحله سوم — Embedding

هر Chunk تبدیل می‌شود به:

```text
Chunk
 ↓
Embedding Model
 ↓
Vector
```

---

## مرحله چهارم — ذخیره در Qdrant

```text
ID
Vector
Payload
```

مثلاً:

```json
{
  "id": 1,
  "vector": [...],
  "payload": {
    "text": "متن سند...",
    "source": "document.pdf",
    "page": 12
  }
}
```

---

## مرحله پنجم — Query

کاربر سؤال می‌پرسد:

```text
پژوهش‌یار چیست؟
```

---

## مرحله ششم — Search

```text
Question
 ↓
Embedding
 ↓
Qdrant
 ↓
Top-K Documents
```

---

## مرحله هفتم — Context

نتایج به LLM داده می‌شوند:

```text
Question
+
Retrieved Context
```

---

## مرحله هشتم — LLM

```text
Ollama
   ↓
LLM
   ↓
Answer
```

---

# ساختار پیشنهادی پروژه Go

```text
qdrant-go-rag/
│
├── cmd/
│   └── server/
│       └── main.go
│
├── internal/
│   │
│   ├── qdrant/
│   │   ├── client.go
│   │   ├── collection.go
│   │   ├── points.go
│   │   └── search.go
│   │
│   ├── embedding/
│   │   └── service.go
│   │
│   ├── rag/
│   │   └── service.go
│   │
│   └── api/
│       ├── handler.go
│       └── routes.go
│
├── config/
│   └── config.go
│
├── docker-compose.yml
├── go.mod
├── go.sum
└── README.md
```

---

# API پیشنهادی

## Health

```http
GET /health
```

---

## Create Collection

```http
POST /collections
```

---

## Insert Document

```http
POST /documents
```

Request:

```json
{
  "title": "Learning Golang",
  "text": "Golang is a programming language..."
}
```

---

## Semantic Search

```http
POST /search
```

Request:

```json
{
  "query": "How can I build an API with Go?",
  "limit": 5
}
```

---

## RAG Chat

```http
POST /chat
```

Request:

```json
{
  "message": "Qdrant چیست؟"
}
```

Response:

```json
{
  "answer": "Qdrant یک Vector Database است..."
}
```

---

# 🔥 مهم‌ترین چیزهایی که باید بلد باشید

اگر قرار باشد فقط مهم‌ترین بخش‌های Qdrant را یاد بگیرید، این موارد اولویت بالاتری دارند:

```text
1. Vector
2. Embedding
3. Collection
4. Point
5. Payload
6. Distance Metrics
7. Similarity Search
8. Top-K
9. Filtering
10. HNSW
11. Indexing
12. Quantization
13. Dense Vector
14. Sparse Vector
15. Hybrid Search
16. RAG
17. Performance
18. Production Deployment
```

---

# 🎯 مسیر یادگیری پیشنهادی

بهتر است این مسیر را دنبال کنید:

```text
Vector
   ↓
Embedding
   ↓
Vector Database
   ↓
Qdrant
   ↓
Collection
   ↓
Point
   ↓
Payload
   ↓
Insert / Upsert
   ↓
Search
   ↓
Filter
   ↓
Index
   ↓
HNSW
   ↓
Quantization
   ↓
Hybrid Search
   ↓
RAG
   ↓
Production
```

---

# 🏆 پروژه نهایی

در پایان ۱۴ روز باید بتوانید چنین سیستمی بسازید:

```text
                    ┌──────────────┐
                    │    User      │
                    └──────┬───────┘
                           │
                           ▼
                    ┌──────────────┐
                    │  Go + Fiber  │
                    └──────┬───────┘
                           │
             ┌─────────────┴─────────────┐
             │                           │
             ▼                           ▼
      ┌──────────────┐            ┌──────────────┐
      │  Embedding   │            │    Redis     │
      │    Model     │            └──────────────┘
      └──────┬───────┘
             │
             ▼
      ┌──────────────┐
      │   Qdrant     │
      │              │
      │ Vector       │
      │ Payload      │
      │ Index        │
      └──────┬───────┘
             │
             ▼
       Top-K Results
             │
             ▼
      ┌──────────────┐
      │    Ollama    │
      │     LLM      │
      └──────┬───────┘
             │
             ▼
      ┌──────────────┐
      │ Final Answer │
      └──────────────┘
```

---

# ✅ نتیجه مورد انتظار

بعد از پایان این دوره باید بتوانید:

* Qdrant را نصب و مدیریت کنید.
* Collection ایجاد کنید.
* Vector ذخیره کنید.
* Point ایجاد و Update کنید.
* Payload طراحی کنید.
* Semantic Search پیاده کنید.
* Filtering انجام دهید.
* Top-K Search انجام دهید.
* Distance Metricها را درک کنید.
* HNSW را بشناسید.
* Index طراحی کنید.
* Quantization را درک کنید.
* Dense و Sparse Vector را بشناسید.
* Hybrid Search پیاده کنید.
* Qdrant را با Golang استفاده کنید.
* Qdrant را با Fiber ترکیب کنید.
* Embedding Model به Qdrant متصل کنید.
* RAG بسازید.
* Qdrant را در Production استفاده کنید.

---

# 🚀 Stack پروژه نهایی

```text
Golang
Fiber
Qdrant
Ollama
Embedding Model
Redis
Docker
REST API
RAG
```

هدف نهایی:

```text
Build → Understand → Optimize → Deploy
```

نه فقط استفاده از Qdrant، بلکه درک کامل معماری Vector Search و استفاده از آن در پروژه‌های واقعی AI.
