# Клиент WMI метрик

Клиент для учебного клиент-серверного проекта по сбору метрик через WMI с Windows устройств на центральном облачном сервере

## Сборка

Собирается только с целевой OS windows из-за библиотеки github.com/yusufpapurcu/wmi

```bash
GOOS=windows GOARCH=amd64 go build -o agent.exe cmd/main.go
```

## Настройка

Переменные окружения

- `SERVER_URL` - Адрес сервера на который отправляются метрики
- `METRICS_PERIOD` - Периодичность сбора метрик в миллисекундах

## Поведение

После подключения первое сообщение содержит общую информацию об ОС, дисках и прочем.
Далее раз в METRICS_PERIOD миллисекунд клиент получает данные с WMI и отправляет на сервер.

## Функционал

Готово:
- [x] Подключение к серверу
- [x] Отправка общей информации о системе
- [x] Регулярная отправка логов

Не готово:
- [ ] Переподключение в случае обрыва соединения
- [ ] Переотправка метрик в случае неудачной отправки
- [ ] Авторизация/аутентификация клиента на сервере

В планах расширить также сами передаваемые метрики
Возможно добавить настройки списка передаваемых метрик

Текущие метрики:

```go
type CPU struct {
	LoadPercentage uint32 `json:"load"`
}

type RAM struct {
	FreePhysicalMemory     uint64 `json:"free_ram"`
	TotalVisibleMemorySize uint64 `json:"total_ram"`
}

type Disk struct {
	DeviceID  string `json:"device_id"`
	FreeSpace uint64 `json:"free"`
	Size      uint64 `json:"total"`
}

type Metrics struct {
	Disks []Disk `json:"disks"`
	OS    []RAM  `json:"ram"`
	Proc  []CPU  `json:"cpu"`
}
```
