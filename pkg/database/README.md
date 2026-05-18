# database

Инициализация БД происходит через автомонтирование в
internal/app/app.go:54 и далее в нужно месте подключается по примеру:

```go
type Sk struct {
    Db database.Db
}

func NewSk(db database.Database) *Sk {
    return &Sk{
        Db: db.Db(),
    }
}
```

Пример: internal/domain/sk/storage/sk.go

Как работает автомонтирование зависимостей можно прочесть тут internal/app.

После инициализации доступен выбор базы данных:

```go
db1 := db.Db()
db2 := db.Slave()
```

В данном проекте конфиг настроен так что Db() отдает основную базу данных postgres, 
а Slave() базу данных mysql. Конфиг может быть сконфигурирован любым образом, например postgres + его реплика.

Для получения экземпляра пакета базы данных доступно несколько методов:

- `Rows()`
- `Row()`
- `Exec()`
- `NamedRows()`
- `NamedRow()`
- `NamedExec()`
- `NamedBatch()`

Описание: `Rows` - получить массив записей, `Row` - получить одну запись, `Exec` - выполнить запрос (без возврата результата),
`Batch` - пакетная вставка данных.

Все разбирать нет смысла, сигнатуры у них почти одинаковые, разберу только два из них - не именованные и именованные (с префиксом Named).

### Неименованный `Row()`
Предназначен для запросов с маленьким количеством аргументов (плейсхолдеров).

```go
var res model.Sk
err := db1.Row(ctx, &res, `select id, name from sk where id = ? and is_valid = ?`, 123, true)
fmt.Println(res) // печать готового результата
```

В этом примере: `ctx` - контект, `&res` - структура в которую записывается результат, 
`select id, name from sk where id = ? and is_valid = ?` - запрос с плейсхолдером, `123` и `true` - переменные которые подставляется вместо плейсхолдеров `?`.

Обратите внимание, что аргумент `&res` передан по ссылке. Это нужно что бы передать пакету `database` процесс копирования 
данных из базы данных в структуру, что бы на выходе получить готовый результат.
Плейсхолдер `?` нужен для безопасности запроса, т.е. предотвращения sql-инъекций и кеширования плана запроса, что ускоряет выполнение запроса.
Пакет так же поддерживает плейсхолдеры `$1, $2`, но в реальном проекте рекомендуется использования одного типа плейсхолдеров на всех запросах.

Результат записывается в структуру model.Sk которая обязательно должна иметь все поля тегом `db` из перечисленных в запросе.
```go
type Sk struct {
    Id        int       `db:"id" `
    Name      string    `db:"name"`
}
```

Т.е. если в запросе написано `select id, name`, то в структуре должно быть два поля с тегами `db:"id"` и `db:"name"`.

### Именованный `NameRow()`
Предназначены для сложных запросов с большим количеством аргументов (плейсхолдеров).

```go
filter := serviceModel.SkFilter{Name: "Test", IsActive: true}

var res model.Sk
err := db1.NamedRow(ctx, &res, `select id, name from sk where id = :id, is_active = :is_active`, filter)
fmt.Println(res) // печать готового результата
```

В этом примере используются именованные плйсхолдеры `:id` и `:is_active`, но содержит только один аргумент `filter` 
который в свою очередь содержит структуру со всеми необходимыми аргументами:

```go
type SkFilter struct {
    Id        int       `db:"id" `
    IsActive  bool      `db:"is_active"`
}
```

Пакет так же поддерживает мапу в качестве агрумента для именованного запроса, данный пример по функциональности идентичен:

```go
filter := map[string]interface{
	"id": 123,
	"is_active": true,
}

var res model.Sk
err := db1.NamedRow(ctx, &res, `select id, name from sk where id = :id, is_active = :is_active`, filter)
```

### JsonWrapper

Дополнительно в пакет встроен тип JsonWrapper (json_wrapper.go) для работы с json-совместимыми полями базы данных.
Использовать этот тип можно как в структуре результата запроса, так и в структуре агрумента.

```go
type Data struct {
	Name string `db:"name"`
	ItemJson database.JsonWrapper `db:"item_json"`
}

data := Data{Name: "Test", ItemJson: database.NewJsonWrapper(somedata)}
err := db1.NamedExec(ctx, `insert into history (name, item_json) values :name, :item_json`, data)
```

В данном примере происходит добавление новой записи в таблицу history с полями name (строка) и item_json (тип jsonb).
В качестве переменно somedata somedata может быть структура или мапа.

Структура JsonWrapper имеет встроенные методы для помощи в разработке:

- `func (j *JsonWrapper) GetInt(key string) (int, error)` - получение числового значения из плоской структуры;
- `func (j *JsonWrapper) GetInt(key string) (string, error)` - получение строкового значения из плоской структуры;
- `func (j *JsonWrapper) Get(output interface{}) error` - преобразование содержимого json в структуру или мапу;
- `func (j *JsonWrapper) GetJson() (string, error)` - получение json в виде строки.
