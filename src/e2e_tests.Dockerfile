FROM python:3.8

# Устанавливаем рабочую директорию внутри контейнера
WORKDIR /app
# Копируем файл go.mod и go.sum внутрь контейнера
COPY ./e2e_tests .

# Загружаем зависимости проекта
RUN pip install -r requirements.txt

CMD []