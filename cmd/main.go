package main

import (
	"fmt"
	"lapas/pkg/controllers"
	"log"
	"net/http"

	mw "lapas/middleware"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter().StrictSlash(true)
	origins := handlers.AllowedOrigins([]string{"*"})

	//	Login
	router.HandleFunc("/login", controllers.Login).Methods("POST")

	api := router.PathPrefix("").Subrouter()
	api.Use(mw.AuthToken)

	//	User
	api.HandleFunc("/user/{idUser}", controllers.GetUser).Methods("GET")
	api.HandleFunc("/user", mw.IsIT(controllers.GetUsers)).Methods("GET")
	api.HandleFunc("/user", mw.IsIT(controllers.CreateUser)).Methods("POST")
	api.HandleFunc("/user/{idUser}", mw.IsIT(controllers.UpdateUser)).Methods("PUT")
	api.HandleFunc("/password", mw.IsIT(controllers.ResetPassword)).Methods("POST")
	api.HandleFunc("/password/{idUser}", controllers.ChangePassword).Methods("PUT")

	//	Survei
	api.HandleFunc("/survei-detail/{slug}", mw.IsSDM(controllers.GetSurvei)).Methods("GET")
	api.HandleFunc("/survei/{slug}", controllers.GetSurveiActived).Methods("GET")
	api.HandleFunc("/survei", controllers.GetSurveis).Methods("GET")
	api.HandleFunc("/survei", mw.IsSDM(controllers.CreateSurvei)).Methods("POST")
	api.HandleFunc("/survei/{idSurvei}", mw.IsSDM(controllers.DeleteSurvei)).Methods("DELETE")
	api.HandleFunc("/survei/{idSurvei}", mw.IsSDM(controllers.UpdateSurvei)).Methods("PUT")
	api.HandleFunc("/survei-duplikasi/{idSurvei}", mw.IsSDM(controllers.DuplicateSurvei)).Methods("POST")
	api.HandleFunc("/survei-status/{idSurvei}", mw.IsSDM(controllers.ChangeStatus)).Methods("PUT")
	api.HandleFunc("/survei-statistik-responden/{idSurvei}", controllers.GetStatistikResponden).Methods("GET")
	api.HandleFunc("/survei-statistik-survei/{idSurvei}/{direktorat}", controllers.GetStatistikJawaban).Methods("GET")
	api.HandleFunc("/survei-responden/{idSurvei}", mw.IsSDM(controllers.GetDataResponden)).Methods("GET")

	// Sub Survei
	api.HandleFunc("/subsurvei", controllers.GetSubSurvei).Methods("GET")
	api.HandleFunc("/subsurvei", mw.IsSDM(controllers.CreateSubSurvei)).Methods("POST")
	api.HandleFunc("/subsurvei/{idSub}", mw.IsSDM(controllers.DeleteSubSurvei)).Methods("DELETE")

	//	Jawaban
	api.HandleFunc("/jawaban/{idSurvei}/{idUser}", controllers.GetTanggapan).Methods("GET")
	api.HandleFunc("/jawaban/{idSurvei}/{idUser}", controllers.SaveJawaban).Methods("POST")

	//	Evaluasi
	api.HandleFunc("/evaluasi/{idSurvei}", controllers.GetEvaluasi).Methods("GET")
	api.HandleFunc("/evaluasi", mw.IsSDM(controllers.CreateEvaluasi)).Methods("POST")
	api.HandleFunc("/evaluasi/{idEvaluasi}", mw.IsSDM(controllers.UpdateEvaluasi)).Methods("PUT")

	//	Laporan
	api.HandleFunc("/laporan/{idLaporan}", controllers.GetLaporan).Methods("GET")
	api.HandleFunc("/laporan", mw.IsSDM(controllers.GetLaporans)).Methods("GET")
	api.HandleFunc("/mylaporan/{idUser}", controllers.GetMyLaporan).Methods("GET")
	api.HandleFunc("/laporan", controllers.CreateLaporan).Methods("POST")
	api.HandleFunc("/tanggapan/{idLaporan}", mw.IsSDM(controllers.CreateTanggapan)).Methods("POST")

	// Surat
	api.HandleFunc("/surat/{idSurat}", controllers.GetSurat).Methods("GET")
	api.HandleFunc("/surat", mw.RSurat(controllers.GetSurats)).Methods("GET")
	api.HandleFunc("/surat", mw.CUDSurat(controllers.CreateSurat)).Methods("POST")
	api.HandleFunc("/surat/{idSurat}", mw.CUDSurat(controllers.UpdateSurat)).Methods("PUT")
	api.HandleFunc("/surat/{idSurat}", mw.CUDSurat(controllers.DeleteSurat)).Methods("DELETE")
	api.HandleFunc("/surat-status/{idSurat}", mw.RSurat(controllers.BeriStatusSurat)).Methods("PUT")

	// Disposisi
	api.HandleFunc("/disposisi/{idDisposisi}", controllers.GetDisposisi).Methods("GET")
	api.HandleFunc("/disposisi", controllers.GetDisposisis).Methods("GET")
	api.HandleFunc("/mydisposisi", controllers.GetMyDisposisis).Methods("GET")
	api.HandleFunc("/disposisi", mw.CDispo(controllers.CreateDisposisi)).Methods("POST")
	api.HandleFunc("/disposisi-status/{idDisposisi}", controllers.BeriStatusDisposisi).Methods("PUT")

	fmt.Println("Server running at: 5000")
	log.Fatal(http.ListenAndServe(":5000", handlers.CORS(origins)(router)))

}
