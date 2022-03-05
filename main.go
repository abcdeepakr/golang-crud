package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

// model
type Course struct {
	CourseId    string  `json:"courseid"`
	CourseName  string  `json:"coursename"`
	CoursePrice string  `json:"courseprice"`
	Author      *Author `json:"author"`
}

type Author struct {
	FullName string `json:"fullname"`
	Website  string `json:"website"`
}

var courses []Course

//middleware

func (c *Course) IsEmpty() bool {
	return c.CourseId == "" && c.CourseName == ""
}

func serveHome(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("<h1>WELCOME TO THE HOMEPAGE</h1>"))
}

func getAllCourses(w http.ResponseWriter, r *http.Request) {
	fmt.Println("get all courses")
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(courses)
}

func getOneCourse(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// grad id from the req

	params := mux.Vars(r)

	for _, course := range courses {
		if course.CourseId == params["id"] {
			json.NewEncoder(w).Encode(course)
			fmt.Println(course)
			return
		}
	}
	json.NewEncoder(w).Encode("NO COURSES FOUND")
}

func createCourse(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Body == nil {
		json.NewEncoder(w).Encode("please send some data")
		return
	}
	var course Course

	_ = json.NewDecoder(r.Body).Decode(&course)
	if course.IsEmpty() {
		json.NewEncoder(w).Encode("Empty data")
		return
	}

	// generate a random number and convert to string

	rand.Seed(time.Now().UnixNano())
	course.CourseId = strconv.Itoa(rand.Intn(100)) // converts to string

	courses = append(courses, course)
	json.NewEncoder(w).Encode(course)
	return
}
func updateCourse(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r) // getting params

	for index, course := range courses {
		if course.CourseId == params["id"] {
			var course Course
			_ = json.NewDecoder(r.Body).Decode(&course)
			courses[index] = course
			fmt.Println(courses)
			// courses = append(courses[:index], courses[index+1:]...)

		}
	}
	json.NewEncoder(w).Encode(courses)
}

// func updateOneCourse(w http.ResponseWriter, r http.Request){

// }
func main() {
	fmt.Println("API MODULE RUNNING")
	r := mux.NewRouter()

	// seeding

	courses = append(courses, Course{CourseId: "12345", CourseName: "React", CoursePrice: "222", Author: &Author{FullName: "Deepak", Website: "test.com"}})
	courses = append(courses, Course{CourseId: "123456", CourseName: "Reacts", CoursePrice: "111", Author: &Author{FullName: "Rajeev", Website: "rajeev.com"}})

	// Listening to a port

	//routing
	r.HandleFunc("/", serveHome).Methods("GET")
	r.HandleFunc("/courses", getAllCourses).Methods("GET")
	r.HandleFunc("/course/{id}", getOneCourse).Methods("POST")
	r.HandleFunc("/course", createCourse).Methods("POST")
	r.HandleFunc("/course/{id}", updateCourse).Methods("PUT")

	log.Fatal(http.ListenAndServe(":4000", r))
}
