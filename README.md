# sales-sphere
Sales record tracking

## Prerequisites

- **Go:** v1.18+
- **MYSQL:** v8.0+
- **Git**

## Getting Started

### 1. Clone the Repository

```bash
git clone https://github.com/NaveenKumar-3108/sales-sphere.git
```
### 2. Configure the Database

Update `dbconfig.toml` with your Mysql credentials:    

```toml
[database]
MysqlServer= "localhost"
MysqlPort= 3306
MysqlUser= "root"
MysqlPassword = "Root"
MysqlDatabase = "sys"
MysqlDBType = "mysql"
```

### 3. Initialize & Run

```bash
go mod tidy
go run main.go
```
App will start on: `http://localhost:29022` 
## 4. LOG file 
To track execution flow and for debugging purpose, the application generates logs during runtime.
LOCATION : log/logfile19042025.12.20.03.610924298.txt

## CSV Upload

Place your sales CSV file in the `./apps/file/` directory. The system loads it automatically on everyday 6am or when a refresh API is triggered.

| Endpoint                         | Method | Params                        |  Header       | Description                                                      
|----------------------------------|--------|-------------------------------|---------------|-------------------------------------------                       
| `/api/refresh`                   | GET    |   none                         |   None        | 
To Refersh data
| `/api/getrevenue`           | GET    | none            |  `StartDate,EndDate,ProductID`    | 
Total revenue 

## Data Refresh
- The refresh API reloads data from CSV to DB.
-Scheduler loads at several interval based on configuration -config.toml 

## Database details in separate file
database_schema.sql

## Schema diagram details in separate file
 Schema.png

## Author
Maintained by NaveenKumar A. 

## Sample API Requests 
1. API: Refresh to load sales data
Method: GET
Route: /api/refresh
Host: localhost:29022

sample request:
http://localhost:29022/api/refresh

**Sample Response** (Success):
```json
{
    "status": "S",
    "errmsg": "Data loaded successfully!!",
}
```

**Sample Response** (error):
```json
{
    "status": "E",
    "errmsg": "Error:RFT01 Error on db connection"
}
```
2. API: To get total revenue 
Method: GET
Route: /api/getrevenue
Host: localhost:29022
Header:StartDate,EndDate,ProductID

sample request:
http://localhost:29022/api/getrevenue

**Sample Response** (Success):
```json
{
    "totalrevenue":1213.21,
    "status": "S",
    "errmsg": "",
}
```

**Sample Response** (error):
```json
{
    "status": "E",
    "errmsg": "Error:RFT01 Error on whille fetching  "
}
```