package stats

import (
	"html/template"
	"net/http"

	// "go.mongodb.org/mongo-driver/bson"
	// "gitlab.com/cpanova/forex/domain/admin"
	"gitlab.com/cpanova/excentral/domain/adminstats"
)

const tpl = `
<html>
    <head>
        <link rel="stylesheet" href="https://unpkg.com/purecss@1.0.1/build/base-min.css">
        <link rel="stylesheet" href="https://unpkg.com/purecss@2.0.3/build/pure-min.css" integrity="sha384-cg6SkqEOCV1NbJoCu11+bm0NvBRc8IYLRGXkmNrqUBfTjmMYwNKPWBTIKyw9mHNJ" crossorigin="anonymous">
    </head>
    <body style="padding: 15">
        <div style="margin: 10 15">
            <table class="pure-table">
                <thead>
                    <tr>
                        <th>Date</th>
                        <th>Leads</th>
                    </tr>
                </thead>
                <tbody>
                    {{range $row := .ByDay}}
                    <tr>
                        <td>{{$row.Day}} </td>
                        <td>{{$row.Count}} </td>
                    </tr>
                    {{end}}
                </tbody>
            </table>
        </div>

		<div style="margin: 10 15">
            <table class="pure-table">
                <thead>
                    <tr>
                        <th>PID</th>
                        <th>Leads</th>
                    </tr>
                </thead>
                <tbody>
                    {{range $row := .ByPID}}
                    <tr>
                        <td>{{$row.PID}} </td>
                        <td>{{$row.Count}} </td>
                    </tr>
                    {{end}}
                </tbody>
            </table>
        </div>

    </body>
</html>
`

// StatsHandler ...
type StatsHandler interface {
	Stats(http.ResponseWriter, *http.Request)
}

type statsHandler struct {
	adminStatsRepo adminstats.Repo
}

// NewHandler ...
func NewHandler(adminStatsRepo adminstats.Repo) StatsHandler {
	return &statsHandler{adminStatsRepo}
}

type statsData struct {
	ByDay []adminstats.DailyReport
	ByPID []adminstats.PIDReport
}

func (h *statsHandler) Stats(w http.ResponseWriter, r *http.Request) {
	// get dates
	// layout := "2006-01-02"

	// var d1, d2 string
	// d1s, ok := r.URL.Query()["from"]
	// if !ok || len(d1s[0]) < 1 {
	// 	d1 = time.Now().UTC().Format("2006-01-02")
	// } else {
	// 	_, err = time.Parse(layout, d1s[0])
	// 	if err != nil {
	// 		http.Error(w, err.Error(), http.StatusBadRequest)
	// 		return
	// 	}
	// 	d1 = d1s[0]
	// }
	// d2s, ok := r.URL.Query()["from"]
	// if !ok || len(d2s[0]) < 1 {
	// 	d2 = time.Now().UTC().Format("2006-01-02")
	// } else {
	// 	_, err = time.Parse(layout, d2s[0])
	// 	if err != nil {
	// 		http.Error(w, err.Error(), http.StatusBadRequest)
	// 		return
	// 	}
	// 	d2 = d2s[0]
	// }

	dailyReport, err := h.adminStatsRepo.ByDay()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	pidReport, err := h.adminStatsRepo.ByPID()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// totalLimit, err := h.adminStatsRepo.TotalLimit()
	// if err != nil {
	// 	http.Error(w, err.Error(), http.StatusInternalServerError)
	// 	return
	// }

	// t, err := template.ParseFiles("http/rest/admin/stats.html")
	t, err := template.New("foo").Parse(tpl)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// d := statsData{
	// 	Stats: []map[string]int{
	// 		map[string]int{"partner": 1, "leads": 100},
	// 		map[string]int{"partner": 2, "leads": 200},
	// 	},
	// 	Cap: 256,
	// }
	d := statsData{
		ByDay: dailyReport,
		ByPID: pidReport,
	}
	t.Execute(w, d)
}
