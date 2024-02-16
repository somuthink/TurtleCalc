package internal

import (
	"fmt"
	"log/slog"
	"math/rand"
	"time"

	ev "github.com/apaxa-go/eval"
)

var servers []chan Operation

func server(id int, ch chan Operation) {
	slog.Info(fmt.Sprintf("SERVER %d started", id))

	for prob := range ch {
		_, err := glob_db.db.Exec(
			fmt.Sprintf(
				"UPDATE `servers` SET `operation` = '%s for problem %d'  WHERE id = %d",
				prob.text,
				prob.id,
				id,
			),
		)
		if err != nil {
			slog.Error("err", err)
		}

		exec_time := ExecutuionTime(prob.text)
		timer := time.NewTimer(time.Duration(exec_time) * time.Millisecond)

		slog.Info(
			fmt.Sprintf("SERVER %d PROCESSING %v WHICH WIL TAKE %d MILISEC", id, prob, exec_time),
		)

		expr, _ := ev.ParseString(prob.text, "")

		r, _ := expr.EvalToData(nil)

		res, _ := r.AsInt()

		slog.Info(
			fmt.Sprintf("SERVER %d CALCULATED %v = %d", id, prob, res),
		)

		<-timer.C

		fmt.Println(glob_db.db.Exec(
			fmt.Sprintf(
				"UPDATE `problems` SET `interm_val` = `interm_val` + %d ,`operations_left`= `operations_left` - 1 WHERE id = %d; UPDATE `problems` SET `answer` = `interm_val` WHERE `operations_left` = 0 AND `id` = %d;",
				res,
				prob.id,
				prob.id,
			),
		))

		glob_db.db.Exec(
			fmt.Sprintf("UPDATE `servers` SET `operation` = 'nothing' WHERE id = %d", id),
		)

		slog.Info(fmt.Sprintf("SERVER %d UPDATED PROBLEM %d", id, prob.id))

	}
}

func CreateServers() {
	servers = make([]chan Operation, get_servers())
	for ind := range servers {
		servers[ind] = make(chan Operation)
	}

	for id, ch := range servers {
		go server(id+1, ch)
	}
}

func TransportOperations(operations []Operation) {
	for _, operation := range operations {
		servers[rand.Intn(len(servers))] <- operation
	}
}

func AgentPrint(str string) {
	fmt.Println(str)
}
