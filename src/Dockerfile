FROM  node:21-alpine as frontend
WORKDIR /code
COPY ./package.json /code/
RUN npm install
COPY ./*.js /code
COPY ./*.json /code
RUN mkdir /code/static
RUN mkdir /code/public
COPY ./static/input.css /code/static
RUN npm run build


FROM golang:1.21 as backend

WORKDIR /workspace
COPY go.mod ./
RUN go mod download
COPY --from=frontend /code/public ./public
COPY . ./
RUN go build -o ./todo ./cmd/


FROM scratch
WORKDIR /
COPY --from=backend /workspace/todo todo
ENTRYPOINT ["/todo"]
