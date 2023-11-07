FROM golang:1.19

# Устанавливаем рабочую директорию внутри контейнера
WORKDIR /app

RUN apt-get update && apt-get install -y tzdata

# Задаем часовой пояс контейнера
ENV TZ=Asia/Tokyo

# Копируем файл go.mod и go.sum внутрь контейнера
COPY . .

# Загружаем зависимости проекта
RUN go mod download

CMD []