# Task Session 1

TASK_URL: https://docs.kodingworks.io/s/17137b9a-ed7a-4950-ba9e-eb11299531c2#h-%F0%9F%8E%AF-tugas

## Overview

Implementasikan CRUD Kategori pada project API

## Model

### Category
- **ID**
- **Name**
- **Description**

## API Endpoints

The application provides several API endpoints for the functionalities mentioned above. Below are some key endpoints:

- **Ambil semua kategori**: `GET /categories`
- **Tambah kategori**: `POST /categories`
- **Update kategori**: `PUT /categories/{id}`
- **Ambil detail satu kategori**: `PGET /categories/{id}`
- **Hapus kategori**: `DELETE /categories/{id}`

## Getting Started

1. **Clone the Repository**:
   ```bash
   git clone <repository-url>
   ```

2. **Install Dependencies**:
   ```bash
   cd <project-directory>
   go mod tidy
   ```

3. **Run the Application**:
   ```bash
   go run main.go 
   ```

4. **Access the API**: Use tools like Postman or cURL to interact with the API endpoints.
    ```bash
   postman collection: docs/categories-api.postman_collection.json
   ```
   Health Check Endpoint:
   ```bash
   curl --location '{Hosted API}/api/v1/categories/health'
   ```
   Display All Categories Endpoint:
   ```bash
   curl --location '{Hosted API}/api/v1/categories'
   ```
   Display Category By ID Endpoint:
   ```bash
   curl --location '{Hosted API}/api/v1/categories/6'
   ```
   Create New Category Endpoint:
   ```bash
   curl --location '{Hosted API}/api/v1/categories' \
   --header 'Content-Type: application/json' \
   --data '{
   "name": "Susu",
   "description": "Kategori Susu"
   }'
   ```
   Update Existing Category Endpoint:
   ```bash
   curl --location --request PUT '{Hosted API}/api/v1/categories/9' \
   --header 'Content-Type: application/json' \
   --data '{
   "name": "Minuman",
   "description": "Kategori Minuman"
   }'
   ```
   Delete Existing Category Endpoint:
   ```bash
   curl --location --request DELETE '{Hosted API}/api/v1/categories/9'
   ```

5. **Hosted API**:

   Localhost:
   ```bash
   http://localhost:8000
   ```
   Railway:
   ```bash
   https://pandusatrianura-categories-api-production.up.railway.app
   ```