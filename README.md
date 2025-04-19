# Online Store Product Management System

## Overview

The Online Store Product Management System is a RESTful API built to manage products in an online store. The system allows users to perform CRUD (Create, Read, Update, Delete) operations on products, manage categories, reviews, and wishlists.

## Technologies and Tools

- Language: Golang
- Framework: Gin (Web Framework)
- Database: PostgreSQL
- Compile SQL: Sqlc
- Authentication: JWT (JSON Web Tokens)
- Container: Docker
- API Documentation: Swagger
- Caching: Redis

## Database Schema

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

2. Setup project:

   ```bash
   make databaseup
   ```

3. Migrate database:

   ```bash
   make migrateup
   ```

4. Run server:

   ```bash
   make server
   ```

## Tests

- Please check the api tests in the folder `/tests/api.http`
