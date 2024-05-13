package main

import (
	"context"
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/GabriellaErlinda/UTS_5027221018_Gabriella-Erlinda/common/genproto/habittracker"
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
	tmpl := template.Must(template.New("index").Parse(habitsTemplate))
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

	// Membuat klien API tugas
	habitClient := habittracker.NewHabitApiClient(client)

	// Membuat habit baru
	_, err = habitClient.CreateHabit(context.Background(), &habittracker.Habit{
		Title:       title,
		Description: description,
	})
	if err != nil {
		http.Error(w, "Failed to create habit", http.StatusInternalServerError)
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

	// Membuat klien API tugas
	habitClient := habittracker.NewHabitApiClient(client)

	// Update habit
	_, err = habitClient.UpdateHabit(context.Background(), &habittracker.Habit{
		Id:          id,
		Title:       title,
		Description: description,
	})
	if err != nil {
		http.Error(w, "Failed to update habit", http.StatusInternalServerError)
		return
	}

	// Tetapkan pengalihan arahan ke halaman list
	http.Redirect(w, r, "/list", http.StatusSeeOther)
}

func (s *httpServer) handleDelete(w http.ResponseWriter, r *http.Request) {
	// Mendapatkan ID habit dari parameter URL
	id := r.URL.Query().Get("id")

	// Menginisialisasi koneksi gRPC
	port := ":" + viper.GetString("app.grpc.port")
	client, err := grpc.Dial(port, grpc.WithInsecure())
	if err != nil {
		http.Error(w, "Could not connect to gRPC server", http.StatusInternalServerError)
		return
	}
	defer client.Close()

	// Membuat klien API tugas
	habitClient := habittracker.NewHabitApiClient(client)

	// Menghapus habit
	_, err = habitClient.DeleteHabit(context.Background(), &wrapperspb.StringValue{Value: id})
	if err != nil {
		http.Error(w, "Failed to delete habit", http.StatusInternalServerError)
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
	habitClient := habittracker.NewHabitApiClient(client)

	// Mengambil daftar habit dari server
	habits, err := habitClient.ListHabits(context.Background(), &emptypb.Empty{})
	if err != nil {
		http.Error(w, "Failed to fetch habits: "+err.Error(), http.StatusInternalServerError)
		log.Printf("Failed to fetch habits: %v\n", err)
		return
	}

	// Menyiapkan data habit untuk ditampilkan di halaman HTML
	type ViewData struct {
		Habits []*habittracker.Habit
	}
	data := ViewData{
		Habits: habits.List,
	}

	// Membuat template HTML
	tmpl := template.Must(template.New("index").Parse(habitsTemplate))

	// Menampilkan template HTML dengan data habit yang telah disiapkan
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

var habitsTemplate = `
<!DOCTYPE html>
<html>
<head>
    <title>Habit Tracker</title>
    <link rel="stylesheet" href="https://stackpath.bootstrapcdn.com/bootstrap/4.5.2/css/bootstrap.min.css">
	<link rel="stylesheet" href="https://cdn.jsdelivr.net/npm/bootstrap-icons@1.3.0/font/bootstrap-icons.css">
    <style>
        body {
            background-color: #F6F3E3;
        }
    </style>
</head>
<body>
	<nav class="navbar navbar-expand-lg navbar-dark bg-info">
		<a class="navbar-brand" href="#">Habit Tracker</a>
		<button class="navbar-toggler" type="button" data-bs-toggle="collapse" data-bs-target="#navbarNav" aria-controls="navbarNav" aria-expanded="false" aria-label="Toggle navigation">
			<span class="navbar-toggler-icon"></span>
		</button>
		<div class="collapse navbar-collapse" id="navbarNav">
			<ul class="navbar-nav">
				<li class="nav-item">
					<a class="nav-link" href="#">Home</a>
				</li>
				<li class="nav-item">
					<a class="nav-link" href="#">Habits</a>
				</li>
				<li class="nav-item">
					<a class="nav-link" href="#">About</a>
				</li>
			</ul>
		</div>
	</nav>

    <div class="container-fluid">
		<h1 class="my-4">Daily Habit Tracker</h1>
		<form action="/create" method="post">
			<div class="form-group">
				<label for="title">Habit name:</label>
				<input type="text" class="form-control" id="title" name="title">
			</div>
			<div class="form-group">
				<label for="description">Description:</label>
				<textarea class="form-control" id="description" name="description"></textarea>
			</div>
			<input type="submit" class="btn btn-success" value="Create Habit">
		</form>
		<hr>
		<!-- Progress Bar -->
		<div class="progress mb-4">
			<div id="progress-bar" class="progress-bar bg-success" role="progressbar" style="width: 0%;" aria-valuenow="0" aria-valuemin="0" aria-valuemax="100"></div>
		</div>
		<!-- Reset Button -->
		<button id="reset-btn" class="btn btn-danger mb-3" onclick="resetHabits()">Reset</button>
		<h2 class="mb-4">Habit List</h2>
		<a href="/list" class="btn btn-primary mb-3">
			<svg width="1em" height="1em" viewBox="0 0 16 16" class="bi bi-arrow-repeat" fill="currentColor" xmlns="http://www.w3.org/2000/svg">
				<path fill-rule="evenodd" d="M8 0a8 8 0 1 0 8 8A8 8 0 0 0 8 0zM3.08 8A5.96 5.96 0 0 1 8 2.97v1.02A5 5 0 1 0 9 8.13H8a.5.5 0 0 1 0-1h3a.5.5 0 0 1 .5.5v3a.5.5 0 0 1-1 0V7.03A7 7 0 1 1 3.08 8z"/>
			</svg>
			<span class="sr-only">Refresh Habit List</span>
		</a>
		<ul class="list-group">
			{{if not (eq (len .Habits) 0)}}
				{{range .Habits}}
				<li class="list-group-item">
					<div class="card-body">
						<h5 class="card-title">{{.Title}}</h5>
						<p class="card-text text-muted">{{.Description}}</p>
						<div class="text-right">
							<input type="checkbox" id="habit{{.Id}}" onchange="updateProgress()">
							<label for="habit{{.Id}}" class="mr-3">Done</label>
							<button class="btn btn-primary" onclick="showUpdateForm('{{.Id}}')">Update</button>
							<button class="btn btn-danger ml-2" onclick="confirmDelete('{{.Id}}')">Delete</button>
						</div>
					</div>
			
					<form id="updateForm{{.Id}}" class="update-form mt-3" action="/update" method="post" style="display: none;">
						<input type="hidden" name="id" value="{{.Id}}">
						<div class="form-group">
							<label for="title{{.Id}}">New Habit:</label>
							<input type="text" class="form-control" id="title{{.Id}}" name="title" value="{{.Title}}">
						</div>
						<div class="form-group">
							<label for="description{{.Id}}">New Description:</label>
							<textarea class="form-control" id="description{{.Id}}" name="description">{{.Description}}</textarea>
						</div>
						<input type="submit" class="btn btn-success" value="Update Habit">
						<a href="/list" class="btn btn-danger ml-2">Back</a>
					</form>
				</li>
				{{end}}
			{{else}}
				<li class="list-group-item">No habits available</li>
			{{end}}
		</ul>
    </div>

<script>
    function showUpdateForm(habitId) {
        var formId = 'updateForm' + habitId;
        var form = document.getElementById(formId);
        if (form.style.display === 'none') {
            form.style.display = 'block';
        } else {
            form.style.display = 'none';
        }
    }

    function updateProgress() {
        var checkboxes = document.querySelectorAll('.list-group-item input[type="checkbox"]');
        var checkedCount = 0;
        checkboxes.forEach(function(checkbox) {
            if (checkbox.checked) {
                checkedCount++;
            }
        });
        var totalCount = checkboxes.length;
        var progress = (checkedCount / totalCount) * 100;
        document.getElementById('progress-bar').style.width = progress + '%';
        document.getElementById('progress-bar').setAttribute('aria-valuenow', progress.toString());
    }

    function resetHabits() {
        var checkboxes = document.querySelectorAll('.list-group-item input[type="checkbox"]');
        checkboxes.forEach(function(checkbox) {
            checkbox.checked = false;
        });
        updateProgress(); // Reset progress bar
    }

    function confirmDelete(habitId) {
        if (confirm("Are you sure to delete this habit?")) {
            window.location.href = "/delete?id=" + habitId;
        }
    }
</script>
</body>
</html>
`
