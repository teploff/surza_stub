# Http-сервис (заглушка)

## Сборка и запуск программы
а) Необходимо установить Go. Подробное описание приведено тут: https://golang.org/doc/install

б) Склонировать проект с github:
```bash
$ git pull https://github.com/teploff/surza_stub.git
```

в) Перейти в репозиторий:
```bash
$ cd surza_stub
```

г) Скомпилировать программу: 
```bash
$ go build main.go
```
д) Запустить: 
```bash
./main
```

## Варианты запуска:
а) Если необходимо указать IP и Port для старта http-сервера, необходимо указать флаг при запуске сервера:
```bash
$ ./main -src="127.0.0.1:8080"
```
По умолчанию src = "127.0.0.1:8091" 
б) Если необходимо указать IP и Port для удаленного http-сервера, с которым предстоит обмен данных, необходимо указать флаг при запуске сервера:
```bash
$ ./main -dest="127.0.0.1:8081"
```
По умолчанию src="127.0.0.1:8092" 
в) Если необходимо указать частоту отправки сообщения удаленному серверу, необходимо указать флаг при запуске сервера:
```
$ ./main freq=275ms
```
По умолчанию freq=1s
г) Объединенный подход:
```bash
$ ./main -src="127.0.0.1:8080" -dest="127.0.0.1:8081" freq=275ms
```

## Режим работы
### Прием данных
Сервис ожидает POST запрос по url sourceAddr/surza, где sourceAdd - адрес Http-сервиса (флаг src при запуске программы), c телом запроса в виде json-а следующего вида:
```json
{
  "q": float64
}
```

### Отправка данных
Сервис отправляет POST запрос по url destAddr/surza, где destAdd - адрес удаленного Http-сервиса (флаг dest при запуске программы), c телом запроса:
```json
{
  "q": float64
}
