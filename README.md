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

### Authentication

### Product Management

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
