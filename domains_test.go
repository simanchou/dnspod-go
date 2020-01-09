package dnspod

import (
	"fmt"
	"net/http"
	"reflect"
	"testing"
)

func TestDomainsService_List(t *testing.T) {
	client, mux, teardown := setupClient()
	defer teardown()

	mux.HandleFunc("/Domain.List", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "unsupported method", http.StatusBadRequest)
			return
		}

		_, _ = fmt.Fprint(w, `{
			"status": {"code":"1","message":""},
			"domains": [
				{
					"id": 2238269,
					"status": "enable"

				},
				{
					"id": 10360095,
					"status": "enable"

				}
			]}`)
	})

	domains, _, err := client.Domains.List()
	if err != nil {
		t.Fatal(err)
	}

	want := []Domain{{ID: "2238269", Status: "enable"}, {ID: "10360095", Status: "enable"}}
	if !reflect.DeepEqual(domains, want) {
		t.Errorf("got %+v, want %+v", domains, want)
	}
}

func TestDomainsService_List_Ambiguous_Value(t *testing.T) {
	client, mux, teardown := setupClient()
	defer teardown()

	mux.HandleFunc("/Domain.List", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "unsupported method", http.StatusBadRequest)
			return
		}

		_, _ = fmt.Fprint(w, `{
			"status": {"code":"1","message":""},
			"domains": [
				{
					"id": 2238269,
					"status": "enable",
					"group_id": 9
				},
				{
					"id": 10360095,
					"status": "enable",
					"group_id": "9"
				}
			]}`)
	})

	domains, _, err := client.Domains.List()
	if err != nil {
		t.Fatal(err)
	}

	want := []Domain{{ID: "2238269", Status: "enable", GroupID: "9"}, {ID: "10360095", Status: "enable", GroupID: "9"}}
	if !reflect.DeepEqual(domains, want) {
		t.Errorf("got %+v, want %+v", domains, want)
	}
}

func TestDomainsService_Create(t *testing.T) {
	client, mux, teardown := setupClient()
	defer teardown()

	mux.HandleFunc("/Domain.Create", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "unsupported method", http.StatusBadRequest)
			return
		}

		w.WriteHeader(http.StatusCreated)
		_, _ = fmt.Fprintf(w, `{"status": {"code":"1","message":""},"domain":{"id":1, "name":"example.com"}}`)
	})

	domainValues := Domain{Name: "example.com"}
	domain, _, err := client.Domains.Create(domainValues)
	if err != nil {
		t.Fatal(err)
	}

	want := Domain{ID: "1", Name: "example.com"}
	if !reflect.DeepEqual(domain, want) {
		t.Errorf("got %+v, want %+v", domain, want)
	}
}

func TestDomainsService_Get(t *testing.T) {
	client, mux, teardown := setupClient()
	defer teardown()

	mux.HandleFunc("/Domain.Info", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "unsupported method", http.StatusBadRequest)
			return
		}

		_, _ = fmt.Fprint(w, `{"status": {"code":"1","message":""},"domain": {"id":1, "name":"example.com"}}`)
	})

	domain, _, err := client.Domains.Get(1)
	if err != nil {
		t.Errorf("Domains.Get returned error: %v", err)
	}

	want := Domain{ID: "1", Name: "example.com"}
	if !reflect.DeepEqual(domain, want) {
		t.Errorf("got %+v, want %+v", domain, want)
	}
}

func TestDomainsService_Delete(t *testing.T) {
	client, mux, teardown := setupClient()
	defer teardown()

	mux.HandleFunc("/Domain.Remove", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "unsupported method", http.StatusBadRequest)
			return
		}

		_, _ = fmt.Fprint(w, `{"status": {"code":"1","message":""}}`)
	})

	_, err := client.Domains.Delete(1)
	if err != nil {
		t.Fatal(err)
	}
}
