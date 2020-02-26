package ikeaiip

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"cloud.google.com/go/bigquery"
	"google.golang.org/api/iterator"
	"github.com/gorilla/schema"
)

// HandleRequest  new blah
func HandleRequest(w http.ResponseWriter, r *http.Request) {
	//fmt.Fprintf(w, "Body, %s", r.ParseForm())
	 // fmt.Fprintf(w, "Body, %s", r.ParseForm())
	//http.HandleFunc("/", IntegrationList)
	//http.HandleFunc("/integrations", IntegrationList)
	//fmt.Fprint(w, r.Header.Get("Content-Type"))
	var intno string
	switch r.Header.Get("Content-Type") {
	case "application/json":
		var d struct {
			IntNo string `json:"intno"`
		}
		err := json.NewDecoder(r.Body).Decode(&d)
		if err != nil {
			fmt.Fprintf(w, "error parsing application/json: %v \n", err)
		} else {
			intno = d.IntNo
		}
	case "text/plain":
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			fmt.Fprintf(w, "error parsing text/plain: %v", err)
		} else {
			intno = string(body)
			fmt.Fprint(w, body)
		}

	}
	
	//Query Param parsing-begin
	type INT_ID struct {
	Id  string `schema:"id"`  // custom name
	
	}
	 var decoder = schema.NewDecoder()
	 var int_id INT_ID
	 // Parse the request from query string
    if err := decoder.Decode(&int_id,r.URL.Query());err != nil {
        // Report any parsing errors
        //w.WriteHeader(http.StatusUnprocessableEntity)
        fmt.Fprintf(w, "Error: %s", err)
        return
    }
	intno = int_id.Id
	
	//Query Param parsing-end
	
	if intno == "" {
		fmt.Fprint(w, "Parameter unavailable ! Fetching entire data set !")
		//return
	}
	//fmt.Fprintf(w, "Final value of intno:, %s\n!", html.EscapeString(intno))
	GetIntegrationList(w, r, intno)

}

// GetIntegrationList blah
func GetIntegrationList(w http.ResponseWriter, r *http.Request, intno string) {
	ctx := context.Background()
	client, err := bigquery.NewClient(ctx, "cloudpoc-267017")
	var whereClause string
	if intno != "" {
		whereClause = " where int_no = '" + intno + "'"
	}
	qrystr := "SELECT int_no,branch,folder,git_loc FROM `cloudpoc-267017.accounts.interfaces` " + whereClause
	query := client.Query(qrystr)
	iter, err := query.Read(ctx)
	PrintResult(w, iter)
	if err != nil {
		// TODO: Handle error.
		fmt.Fprint(w, err, "BAD")
	}
	return

}

// IIPIntegration blah
type IIPIntegration struct {
	Intno  string `bigquery:"int_no"`
	Branch  string  `bigquery:"branch"`
	Folder   string `bigquery:"folder"`
	Gitloc string `bigquery:"git_loc"`
}
// PrintResult blah
func PrintResult(w http.ResponseWriter, iter *bigquery.RowIterator) {
	for {
		var row IIPIntegration
		err := iter.Next(&row)
		if err == iterator.Done {
			//return nil
		}
		if err != nil {
			return
		}

		//fmt.Fprintf(w, "%s %s %s %s\n", row.Gitloc, row.Intno, row.Branch,row.Folder)
		json.NewEncoder(w).Encode(row)
	}

}
