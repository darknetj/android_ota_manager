package lib

import (
  "sync"
  "strings"
  "html/template"
  "path/filepath"
)

// Cached templates
var cachedTemplates = map[string]*template.Template{}
var cachedMutex sync.Mutex
var funcs = template.FuncMap{}

func T(name string) *template.Template {
  cachedMutex.Lock()
  defer cachedMutex.Unlock()

  if t, ok := cachedTemplates[name]; ok {
      return t
  }

  layout := "layout.html"
  layout_path := "views/layout.html"
  if strings.Contains(name, "login") {
    layout = "layout_auth.html"
    layout_path = "views/layout_auth.html"
  }

  t := template.New(layout).Funcs(funcs)
  t = template.Must(t.ParseFiles(
      layout_path,
      filepath.Join("views", name),
  ))

  cachedTemplates[name] = t
  return t
}
