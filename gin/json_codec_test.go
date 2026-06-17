package gin

import (
	"encoding/json"
	"errors"
	"math"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestJSONEncoder_HttpResponse_Struct(t *testing.T) {
	rec := httptest.NewRecorder()
	enc := JSONEncoder{}

	type payload struct {
		Hello string `json:"hello"`
		N     int    `json:"n"`
	}

	enc.HttpResponse(rec, http.StatusCreated, payload{Hello: "world", N: 7})

	if got := rec.Code; got != http.StatusCreated {
		t.Fatalf("status = %d, want %d", got, http.StatusCreated)
	}
	if ct := rec.Header().Get("Content-Type"); ct != "application/json" {
		t.Fatalf("Content-Type = %q, want application/json", ct)
	}

	var got payload
	if err := json.Unmarshal(rec.Body.Bytes(), &got); err != nil {
		t.Fatalf("unmarshal body: %v", err)
	}
	if got.Hello != "world" || got.N != 7 {
		t.Fatalf("body = %+v, want {hello:world n:7}", got)
	}
}

func TestJSONEncoder_HttpResponse_UnencodableReturns500(t *testing.T) {
	rec := httptest.NewRecorder()
	enc := JSONEncoder{}

	type payload struct {
		Val float64 `json:"val"`
	}

	// NaN cannot be represented in JSON; encoding must fail and the response
	// must NOT commit the requested success status.
	enc.HttpResponse(rec, http.StatusOK, payload{Val: math.NaN()})

	if got := rec.Code; got != http.StatusInternalServerError {
		t.Fatalf("status = %d, want %d (success status must not be committed on encode failure)", got, http.StatusInternalServerError)
	}
	if ct := rec.Header().Get("Content-Type"); ct == "application/json" {
		t.Fatalf("Content-Type should not be application/json on encode failure, got %q", ct)
	}
}

func TestJSONEncoder_HttpResponse_ErrorPayload(t *testing.T) {
	rec := httptest.NewRecorder()
	enc := JSONEncoder{}

	enc.HttpResponse(rec, http.StatusBadRequest, errors.New("something failed"))

	if got := rec.Code; got != http.StatusBadRequest {
		t.Fatalf("status = %d, want %d", got, http.StatusBadRequest)
	}
	if ct := rec.Header().Get("Content-Type"); ct != "application/json" {
		t.Fatalf("Content-Type = %q, want application/json", ct)
	}

	var got struct {
		Error string `json:"error"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &got); err != nil {
		t.Fatalf("unmarshal body: %v", err)
	}
	if got.Error != "something failed" {
		t.Fatalf("error field = %q, want %q", got.Error, "something failed")
	}
}

func TestJSONEncoder_HttpResponseError_OutOfRangeStatus(t *testing.T) {
	rec := httptest.NewRecorder()
	enc := JSONEncoder{}

	// code 0 is out of the valid range and must be normalized to 500, while the
	// error still renders as the JSON error envelope.
	enc.HttpResponseError(rec, 0, errors.New("boom"))

	if got := rec.Code; got != http.StatusInternalServerError {
		t.Fatalf("status = %d, want %d", got, http.StatusInternalServerError)
	}
	if ct := rec.Header().Get("Content-Type"); ct != "application/json" {
		t.Fatalf("Content-Type = %q, want application/json", ct)
	}

	var got struct {
		Error string `json:"error"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &got); err != nil {
		t.Fatalf("unmarshal body: %v", err)
	}
	if got.Error != "boom" {
		t.Fatalf("error field = %q, want %q", got.Error, "boom")
	}
}
