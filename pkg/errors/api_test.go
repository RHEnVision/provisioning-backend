// Copyright Red Hat

package errors

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/onsi/gomega"
)

func TestNewInternalServerError(t *testing.T) {
	g := gomega.NewGomegaWithT(t)

	err := NewInternalServerError("Something went wrong")

	g.Expect(err.Code).To(gomega.Equal("ERROR"))
	g.Expect(err.Title).To(gomega.Equal("Something went wrong"))
	g.Expect(err.Status).To(gomega.Equal(http.StatusInternalServerError))
}

func TestNewBadRequest(t *testing.T) {
	g := gomega.NewGomegaWithT(t)

	err := NewBadRequest("Missing required field: Name")

	g.Expect(err.Code).To(gomega.Equal("BAD_REQUEST"))
	g.Expect(err.Title).To(gomega.Equal("Missing required field: Name"))
	g.Expect(err.Error()).To(gomega.Equal("Missing required field: Name"))
	g.Expect(err.Status).To(gomega.Equal(http.StatusBadRequest))
}

func TestNewNotFound(t *testing.T) {
	g := gomega.NewGomegaWithT(t)

	err := NewNotFound("Not found")

	g.Expect(err.Code).To(gomega.Equal("NOT_FOUND"))
	g.Expect(err.Title).To(gomega.Equal("Not found"))
	g.Expect(err.Status).To(gomega.Equal(http.StatusNotFound))
}

func TestRespondWithBadRequest(t *testing.T) {
	g := gomega.NewGomegaWithT(t)

	responseRecorder := httptest.NewRecorder()
	RespondWithBadRequest("Test bad request", responseRecorder)
	var badRequest BadRequest
	json.Unmarshal(responseRecorder.Body.Bytes(), &badRequest)

	g.Expect(responseRecorder.Code).To(gomega.Equal(http.StatusBadRequest))
	g.Expect(badRequest.Title).To(gomega.Equal("Test bad request"))
}

func TestRespondWithInternalServerError(t *testing.T) {
	g := gomega.NewGomegaWithT(t)

	responseRecorder := httptest.NewRecorder()
	RespondWithInternalServerError("Test internal server error", responseRecorder)
	var badRequest BadRequest
	json.Unmarshal(responseRecorder.Body.Bytes(), &badRequest)

	g.Expect(responseRecorder.Code).To(gomega.Equal(http.StatusInternalServerError))
	g.Expect(badRequest.Title).To(gomega.Equal("Test internal server error"))
}
