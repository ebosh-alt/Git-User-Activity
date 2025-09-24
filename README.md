## Git User Activity (git-cli)

Красивый и простой CLI для просмотра публичной активности пользователя GitHub прямо в терминале.
Разработан на Go, используем Cobra для CLI, net/http для работы с GitHub API (без внешних HTTP-клиентов).

### Возможности
- `git-cli activity <username>` — показать последние события пользователя
- Флаг `-t / --tail` — сколько событий вывести (с поддержкой пагинации)
- Флаг `-T / --type` — формат вывода: txt (по умолчанию) или json
- Флаг `-c / --commits-limit` — отображение коммитов
- Человекочитаемый вывод с типом события, репозиторием и временем
- Грамотная обработка ошибок: 404 (пользователь не найден), rate limit и т.п.
- Без внешних HTTP-библиотек (только стандартная библиотека)

### Использование
```bash
go run . git-cli activity [flags] <username>
```

## Пример
```bash
# последние 10 событий (по умолчанию)
go run . git-cli activity ebosh-alt

# последние 10 событий (по умолчанию) c комитами 
go run . git-cli activity -t 10 -c --commits-limit 5 ebosh-alt

# 50 событий
go run . git-cli activity -t 50 ebosh-alt

# JSON-вывод (сырой массив событий API)
go run . git-cli activity -T json ebosh-alt | jq .
```

## Пример текстового вывода:
```css
[2025-09-21 10:21] - Starred vincelwt/chatgpt-mac
[2025-09-20 19:15] - Pushed 1 commit to ebosh-alt/todo_api
• 2a8e599 update
[2025-09-20 19:15] - CreateEvent in ebosh-alt/todo_api
[2025-09-20 19:14] - CreateEvent in ebosh-alt/todo_api
[2025-09-12 11:22] - Pushed 1 commit to ebosh-alt/accounts
• a4ecccd update
[2025-09-12 11:17] - Pushed 1 commit to ebosh-alt/accounts
• 6a80fb0 update
[2025-09-12 09:27] - PublicEvent in ebosh-alt/accounts
[2025-09-12 08:53] - Pushed 1 commit to ebosh-alt/accounts
• f7abb62 update
[2025-09-12 08:50] - Pushed 1 commit to ebosh-alt/accounts
• 8e6019c update
```

## Структура проекта
```go 
github_user_activity/
├─ go.mod
├─ main.go
├─ cmd/
│  ├─ root.go           # корневая команда
│  ├─ git_cli.go        # сабкоманда: git-cli [command]
│  └─ activity.go       # сабкоманда: activity <username>
└─ internal/
    └─ gh/
        └─ client.go      # работа с GitHub API, модели, humanize
```
### Главные части:
- `internal/gh/client.go` — запросы к `https://api.github.com/users/<username>/events`, пагинация, парсинг JSON, форматирование событий
- `gh.Human(e)` — формирует человекочитаемую строку, включая дату/время created_at

## Поведение и ограничения
- Показываются публичные события пользователя
- GitHub возвращает недавние события (обычно ~90 дней и/или до 300 штук)

## Тест-драйв
```bash
go run .git-cli activity -t 15 ebosh-alt

go run . git-cli activity -t 25 -T json ebosh-alt | jq -r '.[].type' | sort | uniq -c

go run . git-cli activity -t 10 -c --commits-limit 5 ebosh-alt
```
