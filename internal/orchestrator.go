package internal

import (
	"fmt"
	"html/template"
	"log/slog"
	"net/http"
	"strconv"
)

func HtmlPage(w http.ResponseWriter, r *http.Request) {
	type data struct {
		Opers []string
	}

	tmpl := template.Must(template.ParseFiles("./templates/index.html"))
	tmpl.Execute(w, data{Opers: []string{"+", "-", "*", "/"}})
}

func SendServers(w http.ResponseWriter, r *http.Request) {
	type machine struct {
		id        int
		operation string
	}

	res, err := glob_db.db.Query("SELECT * FROM `servers` WHERE operation IS NOT NULL")
	if err != nil {
		http.Error(
			w,
			http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError,
		)
	}

	defer res.Close()

	for res.Next() {
		var machine machine
		err := res.Scan(&machine.id, &machine.operation)
		if err != nil {
			fmt.Println(err)
		}

		var li string

		li = fmt.Sprintf(
			"<li><div class='server'><h2>SERVER %d</h2><p>calculating %s</p></div></li>",
			machine.id,
			machine.operation,
		)

		w.Write([]byte(li))

	}
}

func SendOpers(w http.ResponseWriter, r *http.Request) {
	type operation struct {
		operation string
		exec_time int
	}

	opers := []string{"+", "-", "*", "/"}

	for _, oper := range opers {
		exec_time, _ := strconv.Atoi(r.PostFormValue(oper))

		glob_db.db.Exec(fmt.Sprintf(
			"UPDATE `operations` SET exec_time = %d WHERE operation = '%s' ",
			exec_time,
			oper,
		))
	}

	// res, err := glob_db.db.Query("SELECT operation, exec_time FROM `operations` ")
	// if err != nil {
	// 	http.Error(
	// 		w,
	// 		http.StatusText(http.StatusInternalServerError),
	// 		http.StatusInternalServerError,
	// 	)
	// }
	//
	// defer res.Close()
	//
	// for res.Next() {
	// 	var operation operation
	// 	err := res.Scan(&operation.operation, &operation.exec_time)
	// 	if err != nil {
	// 		fmt.Println("error while reading problems")
	// 	}
	//
	// 	li := fmt.Sprintf(
	// 		"<div class='oper-list-elem'><h2>%s</h2> <input type='text' value=%d name=%s placeholder='enter execution time'/></div>",
	// 		operation.operation,
	// 		operation.exec_time,
	// 		operation.operation,
	// 	)
	//
	// 	w.Write([]byte(li))
	//
	// }
}

func SendProbs(w http.ResponseWriter, r *http.Request) {
	type problem struct {
		id              int
		text            string
		interm_val      int
		answer          int
		operations_left int
	}
	var cnt int
	glob_db.db.QueryRow("SELECT COUNT(*) FROM `problems`").Scan(&cnt)

	if cnt == 0 {
		return
	}

	res, err := glob_db.db.Query("SELECT * FROM `problems` ORDER BY id DESC ")
	if err != nil {
		http.Error(
			w,
			http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError,
		)
	}

	defer res.Close()

	for res.Next() {
		var problem problem
		err := res.Scan(
			&problem.id,
			&problem.text,
			&problem.interm_val,
			&problem.answer,
			&problem.operations_left,
		)
		if err != nil {
			fmt.Println(err)
		}

		var li string

		if problem.operations_left != 0 {
			li = fmt.Sprintf(
				"<li><div class='problem'><h2 class='cursive'>%s = ...</h2><p>operations_left: %d</p></div></li>",
				problem.text,
				problem.operations_left,
			)
		} else {
			li = fmt.Sprintf(
				"<li><div class='problem'><h2>%s = %d</h2></div></li>",
				problem.text,
				problem.answer,
			)
		}

		w.Write([]byte(li))

	}
}

func ProblemHandler(w http.ResponseWriter, r *http.Request) {
	problem := r.PostFormValue("problem-enter")

	if problem == "" {
		http.Error(
			w,
			http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError,
		)
	}

	var id, cnt int

	glob_db.db.QueryRow(fmt.Sprintf("SELECT COUNT(*), id FROM `problems` WHERE text = '%s' GROUP BY id LIMIT 1;", problem)).
		Scan(&cnt, &id)

	if cnt == 0 {

		glob_db.db.QueryRow("SELECT COUNT(*) FROM `problems`").Scan(&id)

		id++

		parsedProb, err := parseProblem(problem)

		groupes, err := createGroups(id, parsedProb)

		glob_db.db.Exec(
			fmt.Sprintf(
				"INSERT INTO `problems` (`text`, `operations_left`) values ('%s', %d)",
				problem,
				len(groupes),
			),
		)

		if err != nil {
			slog.Error("err", err)
			http.Error(
				w,
				http.StatusText(http.StatusBadRequest),
				http.StatusBadRequest,
			)
		}

		TransportOperations(groupes)
	}

	fmt.Fprint(w, id)
}
