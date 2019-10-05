package main

import (
	"bytes"
	"encoding/json"
	"github.com/21stio/go-playground/gifts/pkg/boot"
	"github.com/21stio/go-playground/gifts/pkg/communication"
	"github.com/21stio/go-playground/gifts/pkg/types"
	"net/http"
	"sync"
	"testing"
)

func TestServer(t *testing.T) {
	addr := ":22000"

	err := func() (err error) {
		d, err := boot.GetDependencies("employees.json", "gifts.json", addr)
		if err != nil {
			return
		}
		defer d.Server.Close()

		go func() {
			err := d.Server.Serve()
			if err != nil {
				return
			}
		}()

		es, err := d.Storage.GetEmployees()
		if err != nil {
		    return
		}

		names := []string{"wilhelm"}
		for _, e := range es {
			names = append(names, e.Name)
		}
		for _, e := range es {
			names = append(names, e.Name)
		}

		wg := sync.WaitGroup{}
		for _, n := range names {
			wg.Add(1)

			go func(name string) {
				defer wg.Add(-1)

				err := func() (err error) {
					var b bytes.Buffer
					err = json.NewEncoder(&b).Encode(types.AssignEmployeeToGiftRequest{
						EmployeeName: name,
					})
					if err != nil {
						return
					}

					var rsp *http.Response
					rsp, err = http.Post("http://localhost" + addr + communication.AssignPath, "application/json", &b)
					if err != nil {
						return
					}

					rsp0 := types.AssignEmployeeToGiftResponse{}
					err = json.NewDecoder(rsp.Body).Decode(&rsp0)
					if err != nil {
						return
					}

					if rsp0.Error != "" {
						t.Log("rsp.error:", rsp0.Error)
					}

					return
				}()
				if err != nil {
				    t.Error(err)
				}
			}(n)
		}

		wg.Wait()

		return
	}()
	if err != nil {
		t.Error(err)
	}

}
