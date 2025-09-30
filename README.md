# Airport Tools Backend

---

## **Build & Run (Locally)**

1. Склонировать в одну директорию **backend**, **compvis** и **frontend**

2. Подготовить **.env** файл с следующими переменными:

   ```
   DB_URL=
   POSTGRES_HOST=
   POSTGRES_PORT=
   POSTGRES_USER=
   POSTGRES_PASSWORD=
   POSTGRES_DB=
   
   ML_SERVICE_URL=http://ml:8000/api/v1
   
   HTTP_PORT=
   HTTP_READ_TIMEOUT=
   HTTP_WRITE_TIMEOUT=
   
   BUCKET_NAME=
   AWS_ACCESS_KEY_ID=
   AWS_SECRET_ACCESS_KEY=
   AWS_ENDPOINT_URL=
   AWS_REGION=
   ```

3. Написать в терминал `docker-compose up --build` из корня репозитория **backend**

---

## **Requirements**

- **Go**: 1.24.5\+

- **PostgreSQL:** 17.4\+

- **S3 хранилище** или иной способ хранения изображений