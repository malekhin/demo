# Пакет для работы c Clickhouse

Подключение пакета выполняется с помощью автомонтирования:

```go
type History struct {
	cl clickhouse.IClickhouse
}

func NewHistory(cl clickhouse.IClickhouse) *History {
	return &TariffHistory{cl: cl}
}
```

После чего в экземляре пакета будут достпны следующие методы:

- `Rows` - получение нескольких записей
- `Row` - получение одной записи
- `Insert` - асинхронная вставка данных, реализация на уровне кода
- `AsyncInsert` - асинхронная вставка данных, реализация на стороне clickhouse
- `Exec` - для выполнения произвольных запросов

### Вставка данных

```go

type History struct {
    SkId      int64                 `ch:"sk_id"`
    ItemJson   database.JsonWrapper `ch:"item_json"`
    CreatedAt time.Time             `ch:"created_at"`
}

history := History{SkId: 123, ItemJson(someitem), CreatedAt: time.Now()}

err := s.Ch.Insert(ctx, `
    insert into history 
        (sk_id, item_json, created_at)
`, &history)
```

Как видно из примера поля в структуре должны быть помечены тегом `ch` с соответствующим названием поля из запроса.

