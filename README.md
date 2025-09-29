# Price 
I used Gin as http framework because it's much easier. and used slog for structured log.
## postgres
```shell
sudo docker compose up --wait -d
```


### migrate
````shell
make db-migrate-up
````

### run

```shell
make run
```

### we use --cron to run the job trigger.
```shell
make cron
```
## Swagger
http://localhost:8080/swagger/index.html

### Get latest price

```bash
curl -X GET "http://localhost:8080/prices/latest?symbol=btc" \
  -H "Accept: application/json"
```

### 3. Get history (24h default, no interval provided)

```bash
curl -X GET "http://localhost:8080/prices/history?symbol=btc" \
  -H "Accept: application/json"
```

### 4. Get history with valid interval

```bash
curl -X GET "http://localhost:8080/prices/history?symbol=btc&interval=1m" \
  -H "Accept: application/json"
```

### Get history with explicit time range

```bash
FROM=$(date -d "6 hours ago" +%s)
TO=$(date +%s)

curl -X GET "http://localhost:8080/prices/history?symbol=btc&interval=1m&from=$FROM&to=$TO" \
  -H "Accept: application/json"
```
