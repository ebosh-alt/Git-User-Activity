## GitHub User Activity (git-cli)

Красивый и простой CLI для просмотра публичной активности пользователя GitHub прямо в терминале.
Разработан на Go, используем Cobra для CLI, net/http для работы с GitHub API (без внешних HTTP-клиентов).

### Возможности
- `git-cli activity <username>` — показать последние события пользователя
- Флаг `-t / --tail` — сколько событий вывести (с поддержкой пагинации)
- Флаг `-T / --type` — формат вывода: txt (по умолчанию) или json
- Человекочитаемый вывод с типом события, репозиторием и временем
- Грамотная обработка ошибок: 404 (пользователь не найден), rate limit и т.п.
- Без внешних HTTP-библиотек (только стандартная библиотека)

### Использование
```bash
git-cli activity [flags] <username>
```

Флаги:
- `-t, --tail int` — сколько событий показать (по умолчанию: 10; максимум 100 за страницу, есть пагинация)
- `-T, --type string` — формат вывода: txt или json (по умолчанию: txt)
- Стандартные `--help` и автодополнение shell от Cobra

## Пример
```bash
# последние 10 событий (по умолчанию)
git-cli activity ebosh-alt

# 50 событий
git-cli activity -t 50 ebosh-alt

# JSON-вывод (сырой массив событий API)
git-cli activity -T json ebosh-alt | jq .
```

## Пример текстового вывода:
```css 
[2025-09-23 12:41] - Starred vincelwt/chatgpt-mac
[2025-09-22 18:05] - Pushed 3 commits to ebosh-alt/developer-roadmap
[2025-09-21 10:12] - Opened a GitHub issue in ebosh-alt/todo_api
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
- `internal/gh/client.go` — запросы к https://api.github.com/users/<username>/events, пагинация, парсинг JSON, форматирование событий
- `gh.Human(e)` — формирует человекочитаемую строку, включая дату/время created_at

## Поведение и ограничения
- Показываются публичные события пользователя
- GitHub возвращает недавние события (обычно ~90 дней и/или до 300 штук)

## Тест-драйв
```bash
# текстовый вывод
go run .git-cli activity -t 15 ebosh-alt

# JSON для пайплайнов и интеграций
go run . git-cli activity -t 25 -T json ebosh-alt | jq -r '.[].type' | sort | uniq -c
```
