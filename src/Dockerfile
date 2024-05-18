FROM  node:21-alpine as frontend
WORKDIR /code
COPY ./package.json package-lock.json /code/
RUN npm install
COPY ./*.js /code
COPY ./*.json /code
RUN mkdir /code/static
RUN mkdir /code/templates
COPY ./public/tailwind.css ./public/tailwind.css
COPY ./templates/*.html ./templates/
RUN npm run build
RUN cp ./node_modules/htmx.org/dist/htmx.min.js /code/static/htmx..min.js

RUN npm run build


FROM golang:1.21 as backend

WORKDIR /workspace
COPY go.mod ./
RUN go mod download
COPY . ./
COPY --from=frontend /code/static ./static
RUN go build -o ./todo ./cmd/


FROM busybox:latest
WORKDIR /
COPY --from=backend /workspace/todo /todo
RUN ls
CMD ["/todo"]
