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
  {
    "status": "success",
    "data": {
      "username": "testuser2d3",
      "full_name": "Test User",
      "email": "tesd13@example.com",
      "password_changed_at": "0001-01-01T07:06:30+07:06",
      "created_at": "2025-04-19T14:10:11.906889+07:00"
    }
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
    "status": "success",
    "message": "Login successful",
    "data": {
      "session_id": "355c6cfe-b6ac-4c7f-87f7-7f0fe0c1803a",
      "access_token": "eyJhb...",
      "access_token_expires_at": "2025-04-19T14:26:05.339973+07:00",
      "refresh_token": "eyJ..",
      "refresh_token_expires_at": "2025-04-20T14:11:05.340039+07:00",
      "user": {
        "username": "testuser2d3",
        "full_name": "Test User",
        "email": "tesd13@example.com",
        "password_changed_at": "0001-01-01T07:06:30+07:06",
        "created_at": "2025-04-19T14:10:11.906889+07:00"
      }
    }
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
