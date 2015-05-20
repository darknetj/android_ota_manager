package lib

import (
  "sync"
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

  t := template.New("layout.html").Funcs(funcs)

  t = template.Must(t.ParseFiles(
      "views/layout.html",
      filepath.Join("views", name),
  ))

  // r.SetHTMLTemplate(template.Must(template.ParseFiles(baseTemplate, "templates/releases_form.html")))
  // c.HTML(200, "base", data)
  cachedTemplates[name] = t

  return t
}
