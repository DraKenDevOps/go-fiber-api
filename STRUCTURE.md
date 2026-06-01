Here is the project structure based on your screenshot:

```plaintext
GO-FIBER-API/
│
├── config/
│   └── config.go
│
├── database/
│   └── db.go
│
├── handlers/
│   ├── account.go
│   ├── auth.go
│   ├── handler.go
│   ├── pocket.go
│   └── user.go
│
├── logs/
│   ├── demo-rest-api-error.log
│   └── demo-rest-api.log
│
├── middlewares/
│   ├── auth.go
│   └── logger.go
│
├── models/
│   ├── account.go
│   ├── pocket.go
│   ├── transaction.go
│   └── user.go
│
├── routes/
│   └── routes.go
│
├── uploads/
│
├── utils/
│   ├── functions.go
│   └── jwt.go
│
├── zplogger/
│   └── logger.go
│
├── .env
├── .env.example
├── .gitignore
├── go.mod
├── go.sum
├── main.go
├── Makefile
└── schemasql
```

### Brief Folder Purpose

* **config/** → Application configuration setup
* **database/** → Database connection and setup
* **handlers/** → Request handling / business logic
* **logs/** → Application log files
* **middlewares/** → Middleware functions (auth, logging, etc.)
* **models/** → Database models / structs
* **routes/** → API route definitions
* **uploads/** → Uploaded files storage
* **utils/** → Helper functions and JWT utilities
* **zplogger/** → Custom logging implementation
* **main.go** → Entry point of the application
* **Makefile** → Build and automation commands
* **schemasql** → Database schema file or SQL definitions
