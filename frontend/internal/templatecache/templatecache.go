package templatecache

import (
	"html/template"
	"path/filepath"

	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

func New(dir string) (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}

	pages, err := filepath.Glob(filepath.Join(dir, "*.page.html"))
	if err != nil {
		return nil, err
	}

	for _, page := range pages {
		name := filepath.Base(page)

		ts, err := template.New(name).Funcs(template.FuncMap{
			"formattedPrice": formattedPrice,
		}).ParseFiles(page)
		if err != nil {
			return nil, err
		}

		ts, err = ts.ParseGlob(filepath.Join(dir, "*.layout.html"))
		if err != nil {
			return nil, err
		}

		ts, err = ts.ParseGlob(filepath.Join(dir, "*.partial.html"))
		if err != nil {
			return nil, err
		}

		cache[name] = ts
	}

	return cache, nil
}

func formattedPrice(price int64) string {
	return message.NewPrinter(language.English).Sprint(price)
}
