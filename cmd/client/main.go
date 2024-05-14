package main

import (
	"context"
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/nicholasmarco27/UTS_5027221042_Nicholas-Marco-Weinandra/common/genproto/menulist"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

type httpServer struct {
	addr string
}

func NewHttpServer(addr string) *httpServer {
	return &httpServer{addr: addr}
}

func init() {
	args := os.Args[1:]
	var configname string = "default-config"
	if len(args) > 0 {
		configname = args[0] + "-config"
	}
	log.Printf("loading config file %s.yml \n", configname)

	viper.SetConfigName(configname)
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./configs")
	err := viper.ReadInConfig()
	if err != nil {
		panic("Fatal error config file: " + err.Error())
	}
}

func (s *httpServer) Run() error {
	http.HandleFunc("/", s.handleIndex)
	http.HandleFunc("/create", s.handleCreate)
	http.HandleFunc("/update", s.handleUpdate)
	http.HandleFunc("/delete", s.handleDelete)
	http.HandleFunc("/list", s.handleList)
	log.Printf("Starting HTTP server on localhost%s/list\n", s.addr)
	return http.ListenAndServe(s.addr, nil)
}

func (s *httpServer) handleIndex(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.New("index").Parse(menusTemplate))
	if err := tmpl.Execute(w, nil); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// Handlers for CRUD operations
func (s *httpServer) handleCreate(w http.ResponseWriter, r *http.Request) {
	// Mendapatkan data dari form HTML
	title := r.FormValue("title")
	description := r.FormValue("description")

	// Menginisialisasi koneksi gRPC
	port := ":" + viper.GetString("app.grpc.port")
	client, err := grpc.Dial(port, grpc.WithInsecure())
	if err != nil {
		http.Error(w, "Could not connect to gRPC server", http.StatusInternalServerError)
		return
	}
	defer client.Close()

	// Membuat klien API
	menuClient := menulist.NewMenuApiClient(client)

	// Membuat menu baru
	_, err = menuClient.CreateMenu(context.Background(), &menulist.Menu{
		Title:       title,
		Description: description,
	})
	if err != nil {
		http.Error(w, "Failed to create menu", http.StatusInternalServerError)
		return
	}

	// Tetapkan pengalihan arahan ke halaman list
	http.Redirect(w, r, "/list", http.StatusSeeOther)
}

func (s *httpServer) handleUpdate(w http.ResponseWriter, r *http.Request) {
	// Mendapatkan data dari form HTML
	id := r.FormValue("id")
	title := r.FormValue("title")
	description := r.FormValue("description")

	// Menginisialisasi koneksi gRPC
	port := ":" + viper.GetString("app.grpc.port")
	client, err := grpc.Dial(port, grpc.WithInsecure())
	if err != nil {
		http.Error(w, "Could not connect to gRPC server", http.StatusInternalServerError)
		return
	}
	defer client.Close()

	// Membuat klien API menu
	menuClient := menulist.NewMenuApiClient(client)

	// Update menu
	_, err = menuClient.UpdateMenu(context.Background(), &menulist.Menu{
		Id:          id,
		Title:       title,
		Description: description,
	})
	if err != nil {
		http.Error(w, "Failed to update menu", http.StatusInternalServerError)
		return
	}

	// Tetapkan pengalihan arahan ke halaman menu
	http.Redirect(w, r, "/list", http.StatusSeeOther)
}

func (s *httpServer) handleDelete(w http.ResponseWriter, r *http.Request) {
	// Mendapatkan ID menu dari parameter URL
	id := r.URL.Query().Get("id")

	// Menginisialisasi koneksi gRPC
	port := ":" + viper.GetString("app.grpc.port")
	client, err := grpc.Dial(port, grpc.WithInsecure())
	if err != nil {
		http.Error(w, "Could not connect to gRPC server", http.StatusInternalServerError)
		return
	}
	defer client.Close()

	// Membuat klien API menu
	menuClient := menulist.NewMenuApiClient(client)

	// Menghapus menu
	_, err = menuClient.DeleteMenu(context.Background(), &wrapperspb.StringValue{Value: id})
	if err != nil {
		http.Error(w, "Failed to delete menu", http.StatusInternalServerError)
		return
	}

	// Tetapkan pengalihan arahan ke halaman list
	http.Redirect(w, r, "/list", http.StatusSeeOther)
}

func (s *httpServer) handleList(w http.ResponseWriter, r *http.Request) {
	// Menginisialisasi koneksi gRPC
	port := ":" + viper.GetString("app.grpc.port")
	client, err := grpc.Dial(port, grpc.WithInsecure())
	if err != nil {
		http.Error(w, "Could not connect to gRPC server", http.StatusInternalServerError)
		return
	}
	defer client.Close()

	// Membuat klien API tugas
	menuClient := menulist.NewMenuApiClient(client)

	// Mengambil daftar menu dari server
	menus, err := menuClient.ListMenus(context.Background(), &emptypb.Empty{})
	if err != nil {
		http.Error(w, "Failed to fetch menus: "+err.Error(), http.StatusInternalServerError)
		log.Printf("Failed to fetch menus: %v\n", err)
		return
	}

	// Menyiapkan data menu untuk ditampilkan di halaman HTML
	type ViewData struct {
		Menus []*menulist.Menu
	}
	data := ViewData{
		Menus: menus.List,
	}

	// Membuat template HTML
	tmpl := template.Must(template.New("index").Parse(menusTemplate))

	// Menampilkan template HTML dengan data menu yang telah disiapkan
	err = tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func main() {
	httpServer := NewHttpServer(":1000")
	httpServer.Run()
}

var menusTemplate = `
<!DOCTYPE html>
<html>

<head>
    <title>Menu List</title>
    <link rel="stylesheet" href="https://stackpath.bootstrapcdn.com/bootstrap/4.5.2/css/bootstrap.min.css">
    <link rel="stylesheet" href="https://cdn.jsdelivr.net/npm/bootstrap-icons@1.3.0/font/bootstrap-icons.css">
    <style>
        body {
            background-color: #CDE8E5;
        }
    </style>
</head>

<body>
    <div class="container-fluid">
		<div class="container mt-5">
   			 <div class="card mx-auto py-3 rounded" style="max-width: 500px; max-length: 800px;">
        		<div class="card-header text-center bg-info text-white">
            		<h1 class="my-4">RESTAURANT MENU LIST</h1>
        		</div>
        		<div class="card-body">
            	<form action="/create" method="post">
                <div class="form-group">
                    <label for="title">Menu name:</label>
                    <input type="text" class="form-control" id="title" name="title">
                </div>
                <div class="form-group">
                    <label for="description">Price (Rpxx.xxx,-):</label>
                    <textarea class="form-control" id="description" name="description"></textarea>
                </div>
                <button type="submit" class="btn btn-info btn-block">Create Menu</button>
            </form>
        </div>
    </div>
</div>





        <h2 class="my-5 mx-auto"></h2>
        <div class="row">
            {{if not (eq (len .Menus) 0)}}
            {{range .Menus}}
            <div class="col-md-6">
                <div class="card mb-3">
                    <div class="card-body">
                        <h5 class="card-title">{{.Title}}</h5>
                        <p class="card-text text-muted">{{.Description}}</p>
                        <div class="text-right">
                            <button class="btn btn-info" onclick="showUpdateForm('{{.Id}}')">Update</button>
                            <button class="btn btn-danger ml-2" onclick="confirmDelete('{{.Id}}')">Delete</button>
                        </div>
                    </div>

                    <form id="updateForm{{.Id}}" class="update-form mt-3" action="/update" method="post"
                        style="display: none;">
                        <input type="hidden" name="id" value="{{.Id}}">
                        <div class="form-group">
                            <label for="title{{.Id}}">New Menu:</label>
                            <input type="text" class="form-control" id="title{{.Id}}" name="title"
                                value="{{.Title}}">
                        </div>
                        <div class="form-group">
                            <label for="description{{.Id}}">New Price:</label>
                            <textarea class="form-control" id="description{{.Id}}" name="description">{{.Description}}</textarea>
                        </div>
                        <input type="submit" class="btn btn-primary" value="Update Menu">
                        <a href="/list" class="btn btn-warning ml-2">Back</a>
                    </form>
                </div>
            </div>
            {{end}}
            {{else}}
            <div class="col-md-12">
                <div class="card mb-3">
                    <div class="card-body">
                        No menus available
                    </div>
                </div>
            </div>
            {{end}}
        </div>
    </div>

    <script>
        function showUpdateForm(menuId) {
            var formId = 'updateForm' + menuId;
            var form = document.getElementById(formId);
            if (form.style.display === 'none') {
                form.style.display = 'block';
            } else {
                form.style.display = 'none';
            }
        }

        function confirmDelete(menuId) {
            if (confirm("Are you sure to delete this menu?")) {
                window.location.href = "/delete?id=" + menuId;
            }
        }
    </script>
</body>
</html>
`
