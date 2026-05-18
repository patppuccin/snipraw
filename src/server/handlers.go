package server

import (
	"archive/zip"
	"fmt"
	"io"
	"io/fs"
	"mime"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"text/template"

	"github.com/go-chi/chi/v5"
	"github.com/patppuccin/snipraw/src/helpers"
	"github.com/patppuccin/snipraw/src/web"
	"github.com/patppuccin/snipraw/src/web/pages"
)

func notFoundHandler(w http.ResponseWriter, r *http.Request) {
	appCtx, ok := r.Context().Value(appCtxKey).(*appCtx)
	if !ok {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusNotFound)
	pages.NotFoundPage(appCtx.config, "").Render(r.Context(), w)
}

func methodNotAllowedHandler(w http.ResponseWriter, r *http.Request) {
	appCtx, ok := r.Context().Value(appCtxKey).(*appCtx)
	if !ok {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusMethodNotAllowed)
	pages.MethodNotAllowedPage(appCtx.config, "").Render(r.Context(), w)
}

func homePageHandler(w http.ResponseWriter, r *http.Request) {
	appCtx, ok := r.Context().Value(appCtxKey).(*appCtx)
	if !ok {
		return
	}

	projects, err := helpers.ListProjects(appCtx.dir)
	if err != nil {
		if os.IsNotExist(err) {
			appCtx.logger.Warn().Str("dir", appCtx.dir).Msg("snippets directory does not exist")
		} else {
			appCtx.logger.Error().Err(err).Str("dir", appCtx.dir).Msg("failed to list projects")
		}
		projects = []web.Project{}
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	pages.Home(appCtx.config, pages.HomePageProps{Projects: projects}).Render(r.Context(), w)
}

func renderedProjectPageRedirectHandler(w http.ResponseWriter, r *http.Request) {
	project, err := url.PathUnescape(chi.URLParam(r, "project"))
	if err != nil {
		http.NotFound(w, r)
		return
	}
	http.Redirect(w, r, "/project/"+project+"/view/", http.StatusFound)
}

func renderedProjectPageHandler(w http.ResponseWriter, r *http.Request) {
	appCtx, ok := r.Context().Value(appCtxKey).(*appCtx)
	if !ok {
		return
	}

	project, err := url.PathUnescape(chi.URLParam(r, "project"))
	if err != nil {
		http.NotFound(w, r)
		return
	}

	filePath, err := url.PathUnescape(chi.URLParam(r, "*"))
	if err != nil {
		http.NotFound(w, r)
		return
	}

	root := filepath.Join(appCtx.dir, project)
	fullPath := filepath.Join(root, filePath)

	if slices.Contains(strings.Split(filepath.ToSlash(filePath), "/"), ".git") {
		http.NotFound(w, r)
		return
	}

	info, err := os.Stat(fullPath)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	props := pages.ProjectPageProps{
		Project:  project,
		FilePath: filePath,
	}

	name := filepath.Base(fullPath)

	if info.IsDir() {
		entries, err := os.ReadDir(fullPath)
		if err != nil {
			http.Error(w, "failed to read dir", http.StatusInternalServerError)
			return
		}

		var filtered []os.DirEntry
		for _, e := range entries {
			if e.Name() == ".git" {
				continue
			}
			filtered = append(filtered, e)
		}
		props.Entries = filtered

		if readme := helpers.FindReadme(fullPath); readme != nil {
			props.ReadmeHTML = helpers.RenderMarkdown(readme, project, filePath)
		}
	} else if helpers.IsImage(name) || helpers.IsPDF(name) {
		props.FileName = name
	} else {
		props.FileName = name

		const maxRenderSize = 5 * 1024 * 1024 // 5MB

		if info.Size() > maxRenderSize {
			props.TooLarge = true
		} else {
			content, err := os.ReadFile(fullPath)
			if err != nil {
				http.Error(w, "failed to read file", http.StatusInternalServerError)
				return
			}

			var highlighted string
			if helpers.IsMarkdown(name) {
				highlighted = helpers.RenderMarkdown(content, project, filepath.ToSlash(filepath.Dir(filePath)))
			} else {
				highlighted, err = helpers.CodeHighlight(content, name)
				if err != nil {
					highlighted = fmt.Sprintf("<pre>%s</pre>", template.HTMLEscapeString(string(content)))
				}
			}

			props.FileContent = content
			props.HighlightedHTML = highlighted
		}
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	pages.Project(appCtx.config, props).Render(r.Context(), w)
}

func rawContentHandler(w http.ResponseWriter, r *http.Request) {
	appCtx, ok := r.Context().Value(appCtxKey).(*appCtx)
	if !ok {
		return
	}

	project, err := url.PathUnescape(chi.URLParam(r, "project"))
	if err != nil {
		http.NotFound(w, r)
		return
	}

	filePath, err := url.PathUnescape(chi.URLParam(r, "*"))
	if err != nil {
		http.NotFound(w, r)
		return
	}

	fullPath := filepath.Join(appCtx.dir, project, filePath)

	info, err := os.Stat(fullPath)
	if err != nil || info.IsDir() {
		http.NotFound(w, r)
		return
	}

	content, err := os.ReadFile(fullPath)
	if err != nil {
		http.Error(w, "failed to read file", http.StatusInternalServerError)
		return
	}

	ext := strings.ToLower(filepath.Ext(fullPath))
	contentType := mime.TypeByExtension(ext)
	if contentType == "" {
		contentType = http.DetectContentType(content)
	}

	w.Header().Set("Content-Type", contentType)
	w.Header().Set("Content-Disposition", "inline")
	w.Write(content)
}

func downloadContentHandler(w http.ResponseWriter, r *http.Request) {
	appCtx, ok := r.Context().Value(appCtxKey).(*appCtx)
	if !ok {
		return
	}

	project, err := url.PathUnescape(chi.URLParam(r, "project"))
	if err != nil {
		http.NotFound(w, r)
		return
	}

	filePath, err := url.PathUnescape(chi.URLParam(r, "*"))
	if err != nil {
		http.NotFound(w, r)
		return
	}

	fullPath := filepath.Join(appCtx.dir, project, filePath)

	info, err := os.Stat(fullPath)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	if info.IsDir() {
		name := filepath.Base(fullPath)
		w.Header().Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s.zip"`, name))
		w.Header().Set("Content-Type", "application/zip")

		zw := zip.NewWriter(w)
		defer zw.Close()

		err := filepath.WalkDir(fullPath, func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if d.IsDir() {
				return nil
			}

			rel, err := filepath.Rel(fullPath, path)
			if err != nil {
				return err
			}

			f, err := zw.Create(rel)
			if err != nil {
				return err
			}

			file, err := os.Open(path)
			if err != nil {
				return err
			}
			defer file.Close()

			_, err = io.Copy(f, file)
			return err
		})

		if err != nil {
			appCtx.logger.Warn().Err(err).Str("path", fullPath).Msg("failed to zip dir")
		}
	} else {
		name := filepath.Base(fullPath)
		w.Header().Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, name))
		w.Header().Set("Content-Type", "application/octet-stream")

		file, err := os.Open(fullPath)
		if err != nil {
			http.Error(w, "failed to read file", http.StatusInternalServerError)
			return
		}
		defer file.Close()

		io.Copy(w, file)
	}
}
