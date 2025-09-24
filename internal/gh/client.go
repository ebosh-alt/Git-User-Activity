package gh

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

const base = "https://api.github.com"

type Client struct {
	Timeout time.Duration
}

func (c Client) hc() *http.Client {
	t := c.Timeout
	if t == 0 {
		t = 10 * time.Second
	}
	return &http.Client{Timeout: t}
}

type Event struct {
	Type string `json:"type"`
	Repo struct {
		Name string `json:"name"`
	} `json:"repo"`
	CreatedAt time.Time       `json:"created_at"`
	Payload   json.RawMessage `json:"payload"`
}

func (c Client) UserEvents(ctx context.Context, username string, limit int) ([]Event, error) {
	if limit <= 0 {
		limit = 15
	}
	if limit > 100 {
		limit = 100
	}
	url := fmt.Sprintf("%s/users/%s/events?per_page=%d", base, username, limit)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "github-activity-cli")

	res, err := c.hc().Do(req)
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
		}
	}(res.Body)

	if res.StatusCode == 404 {
		return nil, fmt.Errorf("пользователь '%s' не найден (404)", username)
	}
	if res.StatusCode < 200 || res.StatusCode >= 300 {
		body, _ := io.ReadAll(io.LimitReader(res.Body, 1024))
		return nil, fmt.Errorf("ошибка API: %s — %s", res.Status, string(body))
	}

	var out []Event
	if err := json.NewDecoder(res.Body).Decode(&out); err != nil {
		return nil, err
	}
	return out, nil
}

func Human(e Event) string {
	ts := formatTime(e.CreatedAt)
	repo := e.Repo.Name
	switch e.Type {
	case "PushEvent":
		var p struct {
			Size int `json:"size"`
		}
		_ = json.Unmarshal(e.Payload, &p)
		if p.Size == 1 {
			return fmt.Sprintf("[%s] - Pushed 1 commit to %s", ts, repo)
		}
		return fmt.Sprintf("[%s] - Pushed %d commits to %s", ts, p.Size, repo)
	case "IssuesEvent":
		var p struct {
			Action string `json:"action"`
		}
		_ = json.Unmarshal(e.Payload, &p)
		return fmt.Sprintf("[%s] - %s a new issue in %s", ts, title(p.Action), repo)
	case "PullRequestEvent":
		var p struct {
			Action string `json:"action"`
			Number int    `json:"number"`
		}
		_ = json.Unmarshal(e.Payload, &p)
		return fmt.Sprintf("[%s] - %s pull request #%d in %s", ts, title(p.Action), p.Number, repo)
	case "WatchEvent":
		return fmt.Sprintf("[%s] - Starred %s", ts, repo)
	default:
		return fmt.Sprintf("[%s] - %s in %s", ts, e.Type, repo)
	}

}

// HumanDetailed — человекочитаемый вывод, опционально с коммитами для PushEvent
func HumanDetailed(e Event, withCommits bool, maxCommits int) string {
	base := Human(e)
	if !withCommits || e.Type != "PushEvent" {
		return base
	}
	// распарсим коммиты из payload
	var p struct {
		Commits []struct {
			Sha     string `json:"sha"`
			Message string `json:"message"`
			// Url string `json:"url"` // можно дергать для полного сообщения, но не обязательно
		} `json:"commits"`
	}
	_ = json.Unmarshal(e.Payload, &p)
	if len(p.Commits) == 0 {
		return base
	}

	// ограничим и красиво выведем
	if maxCommits > 0 && len(p.Commits) > maxCommits {
		p.Commits = p.Commits[:maxCommits]
	}
	var b strings.Builder
	b.WriteString(base)
	for _, c := range p.Commits {
		sha := c.Sha
		if len(sha) > 7 {
			sha = sha[:7]
		}
		// первая строка сообщения
		msg := c.Message
		if i := strings.IndexByte(msg, '\n'); i >= 0 {
			msg = msg[:i]
		}
		msg = strings.TrimSpace(msg)
		if msg == "" {
			msg = "(no message)"
		}
		b.WriteString("\n  • ")
		b.WriteString(sha)
		b.WriteString(" ")
		b.WriteString(msg)
	}
	return b.String()
}

func formatTime(t time.Time) string {
	if t.IsZero() {
		return "-"
	}
	// локальное время машины; при желании можно заменить на t.UTC()
	return t.UTC().Format("2006-01-02 15:04")
}

func title(s string) string {
	if s == "" {
		return "Did"
	}
	b := []rune(s)
	if b[0] >= 'a' && b[0] <= 'z' {
		b[0] -= 'a' - 'A'
	}
	return string(b)
}
