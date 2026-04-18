# QRIS Payment System API Specification

## Overview
This API provides QRIS (Quick Response Code Indonesian Standard) payment processing capabilities, including merchant inquiry, payment processing, transaction management, and administrative functions.

## Base URL
```
http://localhost:3000
```

## Authentication
All API endpoints (except `/health` and `/swagger/*`) require JWT authentication using the Authorization header:
- `Authorization: Bearer {jwt_token}`

Obtain JWT tokens through the login endpoint.

## Error Response Format
All error responses follow this structure:
```json
{
  "status": "error",
  "errors": "Error message"
}
```

## Endpoints

### Health Check
- **Endpoint**: `GET /health`
- **Description**: Check if the service is running
- **Response**:
  ```json
  {
    "status": "ok"
  }
  ```

### Authentication

#### Register User
- **Endpoint**: `POST /api/auth/register`
- **Description**: Create a new user account with initial balance
- **Request Body**:
  ```json
  {
    "username": "string (alphanum, required)",
    "password": "string (min 8 chars, required)",
    "initial_balance": "number (gte 0, required)"
  }
  ```
- **Response** (201):
  ```json
  {
    "status": "success",
    "data": {
      "token": "string",
      "account_id": "string",
      "balance": "number"
    }
  }
  ```

#### Login User
- **Endpoint**: `POST /api/auth/login`
- **Description**: Authenticate user and get JWT token
- **Request Body**:
  ```json
  {
    "username": "string (alphanum, required)",
    "password": "string (required)"
  }
  ```
- **Response** (200):
  ```json
  {
    "status": "success",
    "data": {
      "token": "string",
      "account_id": "string",
      "balance": "number"
    }
  }
  ```

### QRIS Operations

#### QRIS Inquiry
- **Endpoint**: `GET /api/qris/inquiry/{qris_payload}`
- **Description**: Get merchant data from QRIS payload
- **Headers**: `Authorization: Bearer {jwt_token}`
- **Parameters**:
  - `qris_payload`: string (QRIS payload, can be URL-encoded)
- **Response** (200):
  ```json
  {
    "status": "success",
    "data": {
      "merchant_id": "string",
      "merchant_name": "string",
      "terminal_id": "string",
      "city": "string",
      "fixed_amount": "number",
      "inquiry_id": "string"
    },
    "metadata": {
      "latency_ms": "number",
      "source": "string"
    }
  }
  ```

#### QRIS Inquiry from Image
- **Endpoint**: `POST /api/qris/inquiry/image`
- **Description**: Get merchant data from QRIS image
- **Headers**: `Authorization: Bearer {jwt_token}`
- **Request**: multipart/form-data
  - `image`: file (QR code image)
- **Response** (200): Same as QRIS Inquiry

#### Add Merchant from Image
- **Endpoint**: `POST /api/qris/merchant/image`
- **Description**: Create or reactivate merchant from QRIS image
- **Headers**: `Authorization: Bearer {jwt_token}`
- **Request**: multipart/form-data
  - `image`: file (QR code image)
- **Response** (200):
  ```json
  {
    "status": "success",
    "data": {
      "merchant_id": "string",
      "merchant_name": "string",
      "city": "string",
      "mcc": "string",
      "is_active": "boolean",
      "is_new": "boolean"
    }
  }
  ```

#### Create Payment
- **Endpoint**: `POST /api/qris/payment`
- **Description**: Submit payment for QRIS transaction
- **Headers**: `Authorization: Bearer {jwt_token}`
- **Request Body**:
  ```json
  {
    "inquiry_id": "string (required)",
    "amount": "number (gt 0, required)",
    "payment_method": "string (balance, required)",
    "pincode": "string (len 6, required)"
  }
  ```
- **Response** (202):
  ```json
  {
    "status": "success",
    "data": {
      "status": "string",
      "transaction_id": "string",
      "message": "string",
      "estimated_completion": "string"
    }
  }
  ```

#### Get Transaction Status
- **Endpoint**: `GET /api/qris/status/{transaction_id}`
- **Description**: Get current transaction status
- **Headers**: `Authorization: Bearer {jwt_token}`
- **Parameters**:
  - `transaction_id`: string
- **Response** (200):
  ```json
  {
    "status": "success",
    "data": {
      "transaction_id": "string",
      "status": "string",
      "final_balance": "number",
      "timestamp": "string (ISO 8601)"
    }
  }
  ```

### Admin Operations

#### List Transactions
- **Endpoint**: `GET /api/admin/transactions`
- **Description**: Get all transactions for admin review
- **Headers**: `Authorization: Bearer {jwt_token}`
- **Response** (200):
  ```json
  {
    "status": "success",
    "data": [
      {
        "transaction_id": "string",
        "trace_id": "string",
        "account_id": "string",
        "merchant_id": "string",
        "amount": "number",
        "status": "string",
        "created_at": "string (ISO 8601)",
        "updated_at": "string (ISO 8601)"
      }
    ]
  }
  ```

#### Update Transaction
- **Endpoint**: `PUT /api/admin/transactions/{transaction_id}`
- **Description**: Update transaction status or amount
- **Headers**: `Authorization: Bearer {jwt_token}`
- **Parameters**:
  - `transaction_id`: string
- **Request Body**:
  ```json
  {
    "amount": "number (gt 0, optional)",
    "status": "string (PENDING/SUCCESS/FAILED, optional)"
  }
  ```
- **Response** (200):
  ```json
  {
    "status": "success",
    "data": {
      "transaction_id": "string",
      "trace_id": "string",
      "account_id": "string",
      "merchant_id": "string",
      "amount": "number",
      "status": "string",
      "created_at": "string (ISO 8601)",
      "updated_at": "string (ISO 8601)"
    }
  }
  ```

#### List API Clients
- **Endpoint**: `GET /api/admin/api-clients`
- **Description**: Get all configured API clients
- **Headers**: `Authorization: Bearer {jwt_token}`
- **Response** (200):
  ```json
  {
    "status": "success",
    "data": [
      {
        "client_id": "string",
        "status": "string",
        "created_at": "string (ISO 8601)"
      }
    ]
  }
  ```

#### Create API Client
- **Endpoint**: `POST /api/admin/api-clients`
- **Description**: Add new API client
- **Headers**: `Authorization: Bearer {jwt_token}`
- **Request Body**:
  ```json
  {
    "client_id": "string (required)",
    "client_secret": "string (required)",
    "status": "string (ACTIVE/INACTIVE, required)"
  }
  ```
- **Response** (201):
  ```json
  {
    "status": "success",
    "data": {
      "client_id": "string",
      "status": "string",
      "created_at": "string (ISO 8601)"
    }
  }
  ```

#### Update API Client
- **Endpoint**: `PUT /api/admin/api-clients/{client_id}`
- **Description**: Update API client secret or status
- **Headers**: `Authorization: Bearer {jwt_token}`
- **Parameters**:
  - `client_id`: string
- **Request Body**:
  ```json
  {
    "client_secret": "string (optional)",
    "status": "string (ACTIVE/INACTIVE, optional)"
  }
  ```
- **Response** (200):
  ```json
  {
    "status": "success",
    "data": {
      "client_id": "string",
      "status": "string",
      "created_at": "string (ISO 8601)"
    }
  }
  ```

## Data Types

### Primitives
- `string`: Text data
- `number`: Float64 for amounts, int64 for counts
- `boolean`: True/false values
- `timestamp`: ISO 8601 formatted date-time strings

### Enums
- `payment_method`: "balance"
- `transaction_status`: "PENDING", "SUCCESS", "FAILED"
- `api_client_status`: "ACTIVE", "INACTIVE"

### Validation Rules
- Username: alphanumeric only
- Password: minimum 8 characters
- Pincode: exactly 6 digits
- Amount: greater than 0
- Initial balance: greater than or equal to 0