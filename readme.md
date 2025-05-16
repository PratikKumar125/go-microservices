Clone the Repository:
git clone https://github.com/PratikKumar125/go-microservices.git
cd go-microservices


Install Dependencies:

go mod tidy

Note := Set Up PostgreSQL instance
Note := Update your env yaml's accordingly in /users/env.dev.local and /graphql/graph/dev.env.yaml

Generate gRPC and GraphQL Code:

For Users microservice (gRPC):cd users/usersrpc
protoc --go_out=. --go-grpc_out=. users.proto

For GraphQL Gateway:cd graph
go run github.com/99designs/gqlgen generate


Running the Services
1. Users Microservice

Start the Service:cd users
go run cmd/server/main.go


Details:
Runs a gRPC server on localhost:5002 (configurable in users/internal/config/config.go).
Handles user CRUD operations (e.g., GetUserByEmail, CreateUser).



2. GraphQL Gateway

Start the Service:
cd graph
go run server.go


Details:

Runs a GraphQL server on http://localhost:8080.
Routes queries/mutations to the Users microservice via gRPC.
Access the GraphQL playground at http://localhost:8080.


Example Query:
query {
  users(name: 'pra', email: "") {
    id
    name
    email
  }
}


Using the Migration CLI
The pkg/migrations/cli package provides a CLI to manage database migrations for any microservice.

Navigate to a Microservice:

cd users/cmd/commands

Create a Migration:
go run make:migration --name=create_example_table


Run Migrations:

Applies all pending up migrations.
go run . run:migrations

Rollback Migrations:
go run . migrate:rollback --step=2 (steps are configurable)
Reverts the last applied migration.


Integrating Queue and Storage Packages
Enhance your microservices with my queue and storage packages, available on GitHub:

Queue Package: For message queues (e.g., RabbitMQ, Kafka).
GitHub: github.com/PratikKumar125/go-queue
Usage: go get github.com/PratikKumar125/go-queue@v0.1.0


Storage Package: For object storage (e.g., S3, MinIO).
GitHub: github.com/PratikKumar125/go-storage
Usage: go get github.com/PratikKumar125/go-storage@v0.1.0



See my Medium and LinkedIn posts for details (links in the repository description).
Contributing
I’m open to contributions and feedback! To contribute:

Fork the repository.
Create a feature branch (git checkout -b feature/your-feature).
Commit changes (git commit -m "Add your feature").
Push to the branch (git push origin feature/your-feature).
Open a pull request.


Contact
GitHub: PratikKumar125

Happy coding! 🚀
