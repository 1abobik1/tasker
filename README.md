# Привет)
## Запуск программы ``` docker-compose up ```
## имеются две ручки 
### post localhost:8080/task, пример request json body 
```json 
{
    "type": "fetch_url",
    "payload": "https://example.com/"
} 
```
### get localhost:8080/task?id=task_id
## на данный момент имеется только метод fetch_url, можно добавлять самостоятельно новые io bound задачи в папке io_bounds
## После этого нужно объявить это в main, пример
 ``` registry.Register("fetch_url", fetcher.NewFetchURLProcessor()) ```
