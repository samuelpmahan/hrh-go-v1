
```
homeroomheroes/
├── cmd/
│   └── api/
│       └── main.go
├── internal/
│   ├── admin/
│   │   ├── model.go      # AdminUser (Entity), BulkImportJob (Entity)
│   │   ├── repository.go # Interface(s): AdminUserRepository, BulkImportJobRepository (for data *owned* by admin)
│   │   ├── handler.go    # Admin API handlers (e.g., /admin/teachers/{id}/verify)
│   │   └── service.go    # Admin domain logic (e.g., TeacherVerificationService uses teacherwishlist.TeacherRepository)
│   ├── publicsearch/
│   │   ├── handler.go    # Public search API handlers (e.g., /wishlists/search)
│   │   └── service.go    # Public search domain logic (e.g., WishlistSearchService)
│   ├── schooldirectory/
│   │   ├── model.go      # School (Entity), Address (Value Object)
│   │   └── service.go    # School lookup logic (GetSchool, SearchSchools)
│   ├── teacherwishlist/
│   │   ├── model.go      # Teacher (Entity), Wishlist (Entity), ValidationState (V)
│   │   ├── repository.go # Interfaces: TeacherRepository, WishlistRepository
│   │   └── service.go    # Core domain logic: CreateTeacher, CreateWishlist, UpdateWishlist, etc.
│   └── shared/           # Common utilities, constants, errors used across internal contexts
│       └── errors.go
│       └── constants.go
├── pkg/
│   ├── storage/
│   │   └── postgres/
│   │       └── teacher_repository.go    # Implements teacherwishlist.TeacherRepository
│   │       └── wishlist_repository.go   # Implements teacherwishlist.WishlistRepository
│   │       └── school_repository.go     # Implements schooldirectory.SchoolRepository
│   │       └── admin_repository.go      # Implements admin.AdminUserRepository (if needed)
│   │       └── db.go                    # Database connection setup, common DB utilities
│   └── instrumentation/
│       └── metrics.go                   # Centralized Prometheus metric definitions and helpers
│       └── logging.go                   # Centralized structured logging setup (slog/zap)
├── scripts/
│   ├── setup-db.sql                     # Initial SQL schema for all tables
│   ├── run.sh                           # Simple script to build and run the app
│   └── deploy.sh                        # Simple deployment script (scp, restart systemd)
├── web/
│   └── static/                          # Your static frontend files (HTML, CSS, JS, images)
│       ├── index.html
│       └── pages/
│           └── about.html
│       └── css/
│       └── js/
├── .env.example                         # Example for environment variables
├── .gitignore
├── Dockerfile                           # For containerization (optional, but good for consistent builds)
├── go.mod
├── go.sum
└── README.md
```
