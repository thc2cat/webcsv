package main

import (
	"html/template"
	"io/fs"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"sort"
)

// Structure pour représenter un fichier/dossier
type FileInfo struct {
	Name    string
	IsDir   bool
	Size    int64 // Taille en octets
	ModTime string
}

func main() {
	re, err := regexp.Compile("20.*.csv")
	if err != nil {
		log.Fatal("Erreur lors de la compilation regex :", err)
	}

	// Récupérer le répertoire courant
	currentDir, err := os.Getwd()
	if err != nil {
		log.Fatal("Erreur lors de la récupération du répertoire courant :", err)
	}

	// currentDir := "/home/cat/projects/webcvs/js-source/data"

	http.HandleFunc("/view/", func(w http.ResponseWriter, r *http.Request) {
		data := path.Base(r.URL.Path)

		tmpl, err := template.New("view").Parse(`
			<!DOCTYPE html>
<html lang="en">

<head>
    <meta charset="utf-8">
    <meta content="width=device-width, initial-scale=1, shrink-to-fit=no" name="viewport">
    <meta http-equiv="X-UA-Compatible" content="IE=edge">
    <title>CSV to HTML Table</title>
    <meta name="author" content="Derek Eder">
    <meta content="Display any CSV file as a searchable, filterable, pretty HTML table">

    <!-- Bootstrap core CSS -->
    <link rel="stylesheet" href="https://stackpath.bootstrapcdn.com/bootstrap/4.2.1/css/bootstrap.min.css"
        integrity="sha384-GJzZqFGwb1QTTN6wy59ffF1BuGJpLSa9DkKMp0DgiMDm4iYMj70gZWKYbI706tWS" crossorigin="anonymous">
    <link rel="stylesheet" href="https://cdn.datatables.net/1.10.19/css/dataTables.bootstrap4.min.css">
</head>

<body>
 <div class="container-fluid">
        <main class="row">
            <div class="col">
                <h1> shadowserver.org results for {{.}}</h1>

                <div id="table-container"></div>
            </div>
        </main>
        <footer class="row">
            <div class="col">
                <hr>
                <p class="text-right"><a href="https://github.com/derekeder/csv-to-html-table">CSV to HTML Table</a> by
                    <a href="http://derekeder.com">Derek
                        Eder</a>
                </p>
            </div>
        </footer>
    </div>
   
    <script src="https://code.jquery.com/jquery-3.3.1.min.js"></script>
    <script src="https://cdnjs.cloudflare.com/ajax/libs/twitter-bootstrap/4.2.1/js/bootstrap.bundle.min.js"></script>
    <script src="/js/jquery.csv.min.js"></script>
    <script src="https://cdn.datatables.net/1.10.19/js/jquery.dataTables.min.js"></script>
    <script src="https://cdn.datatables.net/1.10.19/js/dataTables.bootstrap4.min.js"></script>
    <script src="/js/csv_to_html_table.js"></script>

    <script>
        function format_link(link) {
            if (link)
                return "<a href='" + link + "' target='_blank'>" + link + "</a>";
            else return "";
        }

        CsvToHtmlTable.init({
            csv_path: "/data/{{.}}",
            element: "table-container",
            allow_download: true,
            csv_options: {
                separator: ",",
                delimiter: '"'
            },
            datatables_options: {
                paging: false
            },
            custom_formatting: [
                [4, format_link]
            ]
        });
    </script>
</body>

</html>`)
		if err != nil {
			http.Error(w, "Erreur de template : "+err.Error(), http.StatusInternalServerError)
			return
		}

		err = tmpl.Execute(w, data)
		if err != nil {
			http.Error(w, "Erreur d'exécution du template : "+err.Error(), http.StatusInternalServerError)
			return
		}
	})

	// Handler pour la liste des fichiers
	http.HandleFunc("/list", func(w http.ResponseWriter, r *http.Request) {
		files := []FileInfo{}

		// Utilisation de WalkDir pour une meilleure gestion des erreurs et des sous-répertoires
		err := filepath.WalkDir(currentDir, func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				// Gérer les erreurs d'accès aux fichiers/dossiers
				log.Printf("Erreur d'accès à %s : %v", path, err)
				return nil // Continuer la marche même en cas d'erreur
			}
			relPath, _ := filepath.Rel(currentDir, path)
			if relPath == "." {
				return nil
			}

			info, err := d.Info()
			if err != nil {
				log.Printf("Erreur info %s : %v", path, err)
				return nil
			}

			if re.MatchString(d.Name()) {

				files = append(files, FileInfo{
					// Name:    filepath.ToSlash(relPath), // Utiliser des slashs pour les URL
					Name:    d.Name(),
					IsDir:   d.IsDir(),
					Size:    info.Size(),
					ModTime: info.ModTime().Format("2006-01-02 15:04:05"),
				})
			}
			return nil
		})

		if err != nil {
			http.Error(w, "Erreur lors de la lecture du répertoire : "+err.Error(), http.StatusInternalServerError)
			return
		}

		sort.Slice(files, func(i, j int) bool {
			if files[i].IsDir && !files[j].IsDir {
				return true
			}
			if !files[i].IsDir && files[j].IsDir {
				return false
			}
			return files[i].Name < files[j].Name
		})
		// Utilisation d'un template pour générer le HTML
		tmpl, err := template.New("filelist").Parse(`
                        <!DOCTYPE html>
                        <html>
                        <head><title>Liste des fichiers CSV </title></head>
                        <body>
                        <h1>Répertoire : data </h1>
                        <ul>
                                {{range .Files}}
                                <li><a href="/view/{{.Name}}">{{.Name}}</a> ({{if .IsDir}}Dossier{{else}}{{.Size}} octets, modifié le {{.ModTime}}{{end}})</li>
                                {{end}}
                        </ul>
                        </body>
                        </html>
                `)
		if err != nil {
			http.Error(w, "Erreur de template : "+err.Error(), http.StatusInternalServerError)
			return
		}

		data := struct {
			CurrentDir string
			Files      []FileInfo
		}{
			CurrentDir: currentDir,
			Files:      files,
		}
		err = tmpl.Execute(w, data)
		if err != nil {
			http.Error(w, "Erreur d'exécution du template : "+err.Error(), http.StatusInternalServerError)
			return
		}
	})

	// Servir les fichiers statiques (pour les liens)
	fs1 := http.FileServer(http.Dir("./data"))
	http.Handle("/data/", http.StripPrefix("/data/", fs1))
	fs2 := http.FileServer(http.Dir("./js"))
	http.Handle("/js/", http.StripPrefix("/js/", fs2))

	var port = ":8080"
	log.Printf("Serveur démarré sur %s", port)
	log.Fatal(http.ListenAndServe(port, nil))
}
