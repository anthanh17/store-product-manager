# Online Store Product Management System

## Overview

The Online Store Product Management System is a RESTful API built to manage products in an online store. The system allows users to perform CRUD (Create, Read, Update, Delete) operations on products, manage categories, reviews, and wishlists.

> **Note:** This project was developed as part of an `8-hour coding challenge` to build an entire system from design to implementation. Due to the time constraints, not all API endpoints are fully implemented, and there may be some limitations. The challenge of building a complete system in such a short timeframe was an exciting exercise in prioritization and efficient development.

## Features

- Authentication System

  - User registration with password hashing (bcrypt)
  - Login with JWT token authentication
  - Access token and refresh token mechanism
  - Token-based authorization for protected routes

- Rate Limiting - Security: DDOS

  - All APIs are protected with rate limiting (30 requests per minute)
  - Prevents abuse and brute force attacks
  - Redis-based implementation for distributed rate limiting

- SQL Injection Prevention
  - Using sqlc for type-safe database access
  - All SQL queries are parameterized and compiled at build time
  - No string concatenation or dynamic SQL generation
  - Automatic parameter binding prevents injection attacks
  - Strict type checking ensures data integrity

## Technologies and Tools

- Language: Golang
- Framework: Gin (Web Framework)
- Database: PostgreSQL
- Compile SQL: Sqlc
- Authentication: JWT (JSON Web Tokens)
- Container: Docker
- API Documentation: Swagger
- Caching: Redis
- Tools: Makefile, golang-migrate
- Tests: REST Client for Visual Studio Code - test api

## `I have written all the sample runs in file ./tests/api.http`

### Testing API with REST Client

For easy API testing, we recommend installing the REST Client extension for Visual Studio Code:

1. Open Visual Studio Code
2. Go to Extensions (Ctrl+Shift+X or Cmd+Shift+X on Mac)
3. Search for "REST Client" by Huachao Mao
4. Click Install

After installation, you can use the test files in the `/tests/api.http` directory:

- Simply open any `.http` file
- Click the "Send Request" link that appears above each request
- View the response directly in VS Code

Benefits of using REST Client:

- Test API endpoints without leaving your editor
- Save and version control your API tests alongside code
- Share request collections with your team
- Environment variables support for different configurations
- Response syntax highlighting and formatting

## Setup Guide

### System Requirements

- Go 1.23 or higher
- PostgreSQL 14 or higher
- Docker and Docker Compose

#### How to use

1. Clone the repository:

   ```bash
   git clone https://github.com/yourusername/store-product-manager.git
   cd store-product-manager
   ```

2. Start system:

   ```bash
   make systemup
   ```

## Tests

- Please check the api tests in the folder `/tests/api.http`

## Database Schema

![Database Schema](./docs/online_store_product_management_system.png)

### Users

```sql
CREATE TABLE users (
  id SERIAL PRIMARY KEY,
  ext_id uuid NOT NULL,
  username VARCHAR(255) UNIQUE NOT NULL,
  email VARCHAR(255) UNIQUE NOT NULL,
  password_hash VARCHAR(255) NOT NULL,
  full_name VARCHAR(255),
  role VARCHAR(50) DEFAULT 'user',
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  INDEX idx_user_username (emausernameil)
);
```

### Products

```sql
CREATE TABLE products (
  id SERIAL PRIMARY KEY,
  name VARCHAR(255) NOT NULL,
  description TEXT,
  price DECIMAL(10, 2),
  stock_quantity INT,
  status VARCHAR(50),
  image_url VARCHAR(255),
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  INDEX idx_product_name (name)
);
```

### Categories

```sql
CREATE TABLE categories (
  id SERIAL PRIMARY KEY,
  name VARCHAR(255) NOT NULL,
  description TEXT,
  INDEX idx_category (name)
);
```

### Product_Categories

```sql
CREATE TABLE product_categories (
  product_id INT NOT NULL,
  category_id INT NOT NULL,
  PRIMARY KEY (product_id, category_id),
  FOREIGN KEY (product_id) REFERENCES products(id) ON DELETE CASCADE,
  FOREIGN KEY (category_id) REFERENCES categories(id) ON DELETE CASCADE
);
```

### Reviews

```sql
CREATE TABLE reviews (
  id SERIAL PRIMARY KEY,
  product_id INT NOT NULL,
  user_id INT NOT NULL,
  rating INT NOT NULL CHECK (rating >= 1 AND rating <= 5),
  comment TEXT,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (product_id) REFERENCES products(id) ON DELETE CASCADE,
  FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
  INDEX idx_review_product (product_id),
  INDEX idx_review_user (user_id)
);
```

### Wishlist

```sql
CREATE TABLE wishlist (
  user_id INT NOT NULL,
  product_id INT NOT NULL,
  added_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (user_id, product_id),
  FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
  FOREIGN KEY (product_id) REFERENCES products(id) ON DELETE CASCADE
);
```

## API Documentation

### Authentication

#### Register

- **Endpoint**: `POST /api/auth/register`
- **Description**: Register a new user
- **Request Body**:
  ```json
  {
    "username": "user123",
    "email": "user@example.com",
    "password": "securepassword",
    "full_name": "John Doe"
  }
  ```
- **Response**:

  ```json
  HTTP/1.1 200 OK
  Content-Type: application/json; charset=utf-8
  Date: Sat, 19 Apr 2025 12:16:32 GMT
  Content-Length: 183
  Connection: close

  {
    "username": "testuser2drrr3",
    "full_name": "Test User",
    "email": "tesd1dss3@example.com",
    "password_changed_at": "0001-01-01T07:06:30+07:06",
    "created_at": "2025-04-19T19:16:32.032964+07:00"
  }
  ```

#### Login

- **Endpoint**: `POST /api/auth/login`
- **Description**: Login and receive JWT token
- **Request Body**:
  ```json
  {
    "username": "user123",
    "password": "securepassword"
  }
  ```
- **Response**:

  ```json
  {
    "session_id": "006b80e4-7b3a-4d42-b813-3444b79d21a7",
    "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6IjYwZGIyM2Y3LTU2YmUtNDJiOC05MDJlLTYyMzZmNmUyYTRlNSIsInVzZXJuYW1lIjoidGVzdHVzZXIyZDMiLCJyb2xlIjoidXNlciIsImlzc3VlZF9hdCI6IjIwMjUtMDQtMTlUMTk6MTg6MjQuODM1NzE1KzA3OjAwIiwiZXhwaXJlZF9hdCI6IjIwMjUtMDQtMTlUMTk6MzM6MjQuODM1NzE1KzA3OjAwIn0.ahuNNWM4Xrs85SVAcAmfWVTNtMMlGDkeiMhp68DtpF0",
    "access_token_expires_at": "2025-04-19T19:33:24.835715+07:00",
    "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6IjAwNmI4MGU0LTdiM2EtNGQ0Mi1iODEzLTM0NDRiNzlkMjFhNyIsInVzZXJuYW1lIjoidGVzdHVzZXIyZDMiLCJyb2xlIjoidXNlciIsImlzc3VlZF9hdCI6IjIwMjUtMDQtMTlUMTk6MTg6MjQuODM1Nzg0KzA3OjAwIiwiZXhwaXJlZF9hdCI6IjIwMjUtMDQtMjBUMTk6MTg6MjQuODM1Nzg0KzA3OjAwIn0.Z6WPIN9CK7901aQkGFNRPUYysfDrG4cN6a4YzgysE_Q",
    "refresh_token_expires_at": "2025-04-20T19:18:24.835784+07:00",
    "user": {
      "username": "testuser2d3",
      "full_name": "Test User",
      "email": "tesd13@example.com",
      "password_changed_at": "0001-01-01T07:06:30+07:06",
      "created_at": "2025-04-19T17:26:55.901829+07:00"
    }
  }
  ```

#### Renew Access Token

- **Endpoint**: `POST /tokens/renew_access`
- **Description**: Renew Access Token
- **Request Body**:
  ```json
  {
    "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6IjAwNmI4MGU0LTdiM2EtNGQ0Mi1iODEzLTM0NDRiNzlkMjFhNyIsInVzZXJuYW1lIjoidGVzdHVzZXIyZDMiLCJyb2xlIjoidXNlciIsImlzc3VlZF9hdCI6IjIwMjUtMDQtMTlUMTk6MTg6MjQuODM1Nzg0KzA3OjAwIiwiZXhwaXJlZF9hdCI6IjIwMjUtMDQtMjBUMTk6MTg6MjQuODM1Nzg0KzA3OjAwIn0.Z6WPIN9CK7901aQkGFNRPUYysfDrG4cN6a4YzgysE_Q"
  }
  ```
- **Response**:

  ```json
  {
    "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6IjhjMzBiNmNkLTdiOTctNDM0OC1hNjYwLWNhYTIyODIwNjBkNCIsInVzZXJuYW1lIjoidGVzdHVzZXIyZDMiLCJyb2xlIjoidXNlciIsImlzc3VlZF9hdCI6IjIwMjUtMDQtMTlUMTk6MTk6NTkuOTI4NTAzKzA3OjAwIiwiZXhwaXJlZF9hdCI6IjIwMjUtMDQtMTlUMTk6MzQ6NTkuOTI4NTAzKzA3OjAwIn0.W4iYA28Mpv9uEzKJM9VctSBldnecTQ52V5G9L6pbnAk",
    "access_token_expires_at": "2025-04-19T19:34:59.928503+07:00"
  }
  ```

---

### Product Management

#### Create New Product

- **Endpoint**: `POST /api/products`
- **Description**: Create a new product
- **Request Body**:
  ```json
  {
    "name": "Product B",
    "description": "Description of Product B",
    "price": 150000,
    "stock_quantity": 30,
    "status": "active",
    "image_url": "https://example.com/image2.jpg",
    "category_ids": [1, 2]
  }
  ```
- **Response**:
  ```json
  {
    "status": "success",
    "message": "Product created successfully",
    "data": {
      "id": 2,
      "name": "Product B",
      "description": "Description of Product B",
      "price": 150000,
      "stock_quantity": 30,
      "status": "active",
      "image_url": "https://example.com/image2.jpg",
      "categories": [
        {
          "id": 1,
          "name": "Category 1"
        },
        {
          "id": 2,
          "name": "Category 2"
        }
      ],
      "created_at": "2023-01-03T00:00:00Z",
      "updated_at": "2023-01-03T00:00:00Z"
    }
  }
  ```

#### Get Product Details

- **Endpoint**: `GET /api/products/{id}`
- **Description**: Get detailed information about a product
- **Response**:
  ```json
  {
    "status": "success",
    "data": {
      "id": 1,
      "name": "Product A",
      "description": "Description of Product A",
      "price": 100000,
      "stock_quantity": 50,
      "status": "active",
      "image_url": "https://example.com/image.jpg",
      "categories": [
        {
          "id": 1,
          "name": "Category 1"
        }
      ],
      "reviews": [
        {
          "id": 1,
          "user_id": 2,
          "username": "user456",
          "rating": 5,
          "comment": "Great product",
          "created_at": "2023-01-02T00:00:00Z"
        }
      ],
      "created_at": "2023-01-01T00:00:00Z",
      "updated_at": "2023-01-01T00:00:00Z"
    }
  }
  ```

#### Get Product List (Paginated)

- **Endpoint**: `GET /api/products`
- **Description**: Get a paginated list of products with filtering
- **Query Parameters**:
  - `page`: Page number (default: 1)
  - `limit`: Number of products per page (default: 10)
  - `sort`: Field to sort by (default: created_at)
  - `search_product_name`: Search product name keyword
  - `category_id`: Filter by category
  - `price_sort`: sort price (asc/desc, default: desc)
  - `status`: Product status
- **Response**:
  ```json
  {
    "status": "success",
    "data": {
      "products": [
        {
          "id": 1,
          "name": "Product A",
          "description": "Description of Product A",
          "price": 100000,
          "stock_quantity": 50,
          "status": "active",
          "image_url": "https://example.com/image.jpg",
          "categories": [
            {
              "id": 1,
              "name": "Category 1"
            }
          ],
          "created_at": "2023-01-01T00:00:00Z",
          "updated_at": "2023-01-01T00:00:00Z"
        }
      ],
      "pagination": {
        "total": 100,
        "page": 1,
        "limit": 10,
        "total_pages": 10
      }
    }
  }
  ```

#### Update Product

- **Endpoint**: `PUT /api/products/{id}`
- **Description**: Update product information
- **Request Body**:
  ```json
  {
    "name": "Product B (Updated)",
    "description": "Updated description of Product B",
    "price": 160000,
    "stock_quantity": 25,
    "status": "active",
    "image_url": "https://example.com/image2_updated.jpg",
    "category_ids": [1, 3]
  }
  ```
- **Response**:
  ```json
  {
    "status": "success",
    "message": "Product updated successfully",
    "data": {
      "id": 2,
      "name": "Product B (Updated)",
      "description": "Updated description of Product B",
      "price": 160000,
      "stock_quantity": 25,
      "status": "active",
      "image_url": "https://example.com/image2_updated.jpg",
      "categories": [
        {
          "id": 1,
          "name": "Category 1"
        },
        {
          "id": 3,
          "name": "Category 3"
        }
      ],
      "updated_at": "2023-01-04T00:00:00Z"
    }
  }
  ```

#### Delete Product

- **Endpoint**: `DELETE /api/products/{id}`
- **Description**: Delete a product
- **Response**:
  ```json
  {
    "status": "success",
    "message": "Product deleted successfully"
  }
  ```

### Category Management

#### Create New Category

- **Endpoint**: `POST /api/categories`
- **Description**: Create a new category
- **Request Body**:
  ```json
  {
    "name": "Category 3",
    "description": "Description of Category 3"
  }
  ```
- **Response**:
  ```json
  {
    "status": "success",
    "message": "Category created successfully"
  }
  ```

#### Get Category List

- **Endpoint**: `GET /api/categories`
- **Description**: Get a list of all categories
- **Response**:
  ```json
  {
    "status": "success",
    "data": [
      {
        "id": 1,
        "name": "Category 1",
        "description": "Description of Category 1"
      },
      {
        "id": 2,
        "name": "Category 2",
        "description": "Description of Category 2"
      }
    ]
  }
  ```

### Dashboard

- **Endpoint**: `GET /api/dashboard/summary`
- **Description**: View dashboard summary

- **Response**:

  ```json
  {
    "status": "success",
    "data": {
      "total_categories": 5,
      "total_products": 25,
      "categories": [
        {
          "id": 1,
          "name": "Electronics",
          "description": "Electronic devices and gadgets",
          "product_count": 10
        },
        {
          "id": 2,
          "name": "Clothing",
          "description": "Apparel and fashion items",
          "product_count": 8
        },
        {
          "id": 3,
          "name": "Books",
          "description": "Books and publications",
          "product_count": 7
        }
      ],
      "status_summary": [
        {
          "status": "IN_STOCK",
          "count": 20
        },
        {
          "status": "OUT_OF_STOCK",
          "count": 5
        }
      ]
    }
  }
  ```

### Review Management

#### Add Review

- **Endpoint**: `POST /api/products/{id}/reviews`
- **Description**: Add a review for a product
- **Request Body**:
  ```json
  {
    "rating": 5,
    "comment": "Great product"
  }
  ```
- **Response**:
  ```json
  {
    "status": "success",
    "message": "Review added successfully",
    "data": {
      "id": 1,
      "product_id": 1,
      "user_id": 2,
      "username": "user456",
      "rating": 5,
      "comment": "Great product",
      "created_at": "2023-01-02T00:00:00Z"
    }
  }
  ```

### Wishlist Management

#### Get Wishlist

- **Endpoint**: `GET /api/wishlist`
- **Description**: Get the user's wishlist of products
- **Response**:
  ```json
  {
    "status": "success",
    "data": [
      {
        "id": 1,
        "name": "Product A",
        "description": "Description of Product A",
        "price": 100000,
        "image_url": "https://example.com/image.jpg",
        "added_at": "2023-01-05T00:00:00Z"
      }
    ]
  }
  ```

#### Add to Wishlist

- **Endpoint**: `POST /api/wishlist/{product_id}`
- **Description**: Add a product to the wishlist
- **Response**:
  ```json
  {
    "status": "success",
    "message": "Product added to wishlist"
  }
  ```

#### Remove from Wishlist

- **Endpoint**: `DELETE /api/wishlist/{product_id}`
- **Description**: Remove a product from the wishlist
- **Response**:
  ```json
  {
    "status": "success",
    "message": "Product removed from wishlist"
  }
  ```
