package main

import (
	"fmt"
	"lapas/pkg/controllers"
	"log"
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter().StrictSlash(true)
	origins := handlers.AllowedOrigins([]string{"*"})

	//	User
	router.HandleFunc("/user/{idUser}", controllers.GetUser).Methods("GET")
	router.HandleFunc("/user", controllers.GetUsers).Methods("GET")
	router.HandleFunc("/user", controllers.CreateUser).Methods("POST")
	router.HandleFunc("/user/{idUser}", controllers.UpdateUser).Methods("PUT")
	router.HandleFunc("/password", controllers.ResetPassword).Methods("POST")
	router.HandleFunc("/password/{idUser}", controllers.ChangePassword).Methods("PUT")

	//	Survei
	router.HandleFunc("/survei-detail/{slug}", controllers.GetSurvei).Methods("GET")
	router.HandleFunc("/survei/{slug}", controllers.GetSurveiActived).Methods("GET")
	router.HandleFunc("/survei", controllers.GetSurveis).Methods("GET")
	router.HandleFunc("/survei", controllers.CreateSurvei).Methods("POST")
	router.HandleFunc("/survei/{idSurvei}", controllers.DeleteSurvei).Methods("DELETE")
	router.HandleFunc("/survei/{idSurvei}", controllers.UpdateSurvei).Methods("PUT")
	router.HandleFunc("/survei-duplikasi/{idSurvei}", controllers.DuplicateSurvei).Methods("POST")
	router.HandleFunc("/survei-status/{idSurvei}", controllers.ChangeStatus).Methods("PUT")
	router.HandleFunc("/survei-statistik-responden/{idSurvei}", controllers.GetStatistikResponden).Methods("GET")
	router.HandleFunc("/survei-statistik-survei/{idSurvei}/{direktorat}", controllers.GetStatistikJawaban).Methods("GET")
	router.HandleFunc("/survei-responden/{idSurvei}", controllers.GetDataResponden).Methods("GET")

	// Sub Survei
	router.HandleFunc("/subsurvei", controllers.GetSubSurvei).Methods("GET")
	router.HandleFunc("/subsurvei", controllers.CreateSubSurvei).Methods("POST")
	router.HandleFunc("/subsurvei/{idSub}", controllers.DeleteSubSurvei).Methods("DELETE")

	//	Jawaban
	router.HandleFunc("/jawaban/{idSurvei}/{idUser}", controllers.GetTanggapan).Methods("GET")
	router.HandleFunc("/jawaban/{idSurvei}/{idUser}", controllers.SaveJawaban).Methods("POST")

	//	Evaluasi
	router.HandleFunc("/evaluasi/{idSurvei}", controllers.GetEvaluasi).Methods("GET")
	router.HandleFunc("/evaluasi", controllers.CreateEvaluasi).Methods("POST")
	router.HandleFunc("/evaluasi/{idEvaluasi}", controllers.UpdateEvaluasi).Methods("PUT")

	fmt.Println("Server running at: 5000")
	log.Fatal(http.ListenAndServe(":5000", handlers.CORS(origins)(router)))

}
