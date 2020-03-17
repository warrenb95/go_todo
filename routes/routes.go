package routes

// DisplayOverviewEndPoint nto sure what this does oops
// func DisplayOverviewEndPoint(res http.ResponseWriter, req *http.Request) {
// 	resbody, err := http.Get(mongodbURL)
// 	if err != nil {
// 		// Handke error
// 	}
// 	defer resbody.Body.Close()

// 	resbodybytes, err := ioutil.ReadAll(resbody.Body)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	var todos []models.Todo
// 	json.Unmarshal(resbodybytes, &todos)

// 	data := IndexPageData{Todos: todos}

// 	err = templates.ExecuteTemplate(res, "overview", data)
// 	if err != nil {
// 		http.Error(res, err.Error(), http.StatusInternalServerError)
// 		return
// 	}
// }
