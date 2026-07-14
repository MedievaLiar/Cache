
# In-memory кэш с TTL

Потокобезопасный In-memory кэш с поддержкой TTL (time-to-live) записей


## Возможности

**Set(key, value, ttl)** - добавление значения в кэш

**Get(key) (value, bool)** - получение значения из кэша

**Delete(key)** - удаление значения из кэша

**Cleanup** - очистка просроченных записей

**Exists(key) bool** - проверка существования ключа

**Keys()** - получение всех ключей

**MaxSize** - ограничение максимального размера кэша


## Особенности

- Потокобезопасность - кэш безопасен для использования из нескольких горутин
- TTL для каждой записи
- Автоматическая очистка - фоновая горутина периодически удаляет просроченные записи
- Ограничение размера кэша
- Поддержка любых типов данных



## Установка

go get github.com/MedievaLiar/Cache

## Использование



    // Создание кэша

    // defaultTTL - время жизни по умолчанию (5 минут)
    // cleanupInterval - интервал очистки (10 секунд)
    // maxSize - максимальный размер (1000 записей)

    c := cache.NewCache(5*time.Minute, 10*time.Second, 1000)
    defer c.Stop()

    // Добавление записей
    c.Set("user1", "Ann", 0)                  // использует дефолтный TTL
    c.Set("user2", "Jim", 30*time.Second)     // свой TTL

    // Получение записи
    if val, ok := c.Get("user1"); ok {
        fmt.Printf("User 1: %v\n", val)
    }

    // Проверка существования
    if c.Exists("user2") {
        fmt.Println("User 2 exists")
    }

    // Получение всех ключей
    keys := c.Keys()
    fmt.Printf("Все ключи: %v\n", keys)

    // Удаление записи
    c.Delete("user1")

## Тестирование

Запуск тестов:

```go test ./cache -v```

Пример вывода:

```
=== RUN   TestSetAndGet
--- PASS: TestSetAndGet (0.00s)
=== RUN   TestTTL
--- PASS: TestTTL (0.15s)
=== RUN   TestDelete
--- PASS: TestDelete (0.00s)
=== RUN   TestMaxSize
--- PASS: TestMaxSize (0.00s)
=== RUN   TestKeys
--- PASS: TestKeys (0.00s)
PASS
ok      github.com/MedievaLiar/Cache/cache      0.405s
```


## Стратегия при переполнении кэша

При достижении максимального размера кэша:

1. Сначала удаляются просроченные записи

2. Затем удаляется запись с самым близким истечением TTL

3. Если не получается – удаляется любая
## Зависимости
Go 1.21+
