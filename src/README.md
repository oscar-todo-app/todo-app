# Go, htmx, tailwind and postgres based app for organizing a todo list.

 Project to get more comfortable in writting CRUD apps and interacting with a db.

## How to run and test the project.
 - Tests are set up with make test. Or go test.
 - To run a dev server we can run the make run or docker-compose up -d. And export the env var for the db.

## Structure.

 - Cmd folder:
    - Main.go file that sets up the routes and the env vars.
    - handlres.go file handles the handlers for the routes.
    - templates.go file the handles the rendering of the templates.
    - validate.go file that takes care of the validation for some forms to make sure we get the right data.
    - handlers_test.go file that contains the test for the package.
 - Internal folder:
    - main.go contains all of our interactions to the postgresdb, handled with pgx.
    - main_test.go integration test for the db interactions.
    - testutils.go set up a ephemeral postgres container to run tests against.
 - Migrations folder:
    - migrate.go package to perform the db migrations as part of go code that then can  be called on startup.
 - Migrate folder:
    - Contains all the migrations needed to run the db using goose library.
    - We embed the files as part of the go binary so we can be sure that they are always there.
 - Static folder:
    - Contains the css and js files for the app to run, generated from tailwind.
    - And a go file to embed as a filesystem.
 - Templates folder:
    -  Contains all the go templates for the html layouts, also embeded for ease of use.
 - On root folder:
    - Dockerfile: Multistage build that compiles the frontend and builds the go binary and then creates a small images with the binary
    - Go mod and sum for the go deps.
    - package.json for the frontend depencies
    - tailwind.config.js for the tailwind conf.
    - docker-compose for the local env build
    - Makefile to wrap some comands and make it easier to use.
