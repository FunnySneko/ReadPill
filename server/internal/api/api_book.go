package api

import (
	"encoding/json"
	"net/http"
)

func (s *ServerHandler) GetBooks(w http.ResponseWriter, r *http.Request) {
	books, err := s.Db.GetBooks(nil, nil)
	if err != nil {
		ErrorOut(w, err, "internal server error", http.StatusInternalServerError)
	}

	b := BookArray{}
	for _, book := range books {
		b = append(b, mapBookToAPI(book))
	}

	response, err := json.Marshal(b)
	if err != nil {
		ErrorOut(w, err, "internal server error", http.StatusInternalServerError)
	}

	w.Write(response)
}

func (s *ServerHandler) PostBooks(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		ErrorOut(w, err, "internal server error", http.StatusInternalServerError)
	}

	bp := BookPost{}
	data := r.FormValue("data")
	err = json.Unmarshal([]byte(data), &bp)
	if err != nil {
		ErrorOut(w, err, "internal server error", http.StatusInternalServerError)
	}

	file, header, err := r.FormFile("image")
	if err != nil {
		ErrorOut(w, err, "internal server error", http.StatusInternalServerError)
	}

	imagePath, err := s.Ah.UploadImage(file, *header)

	s.Ah.PostBook(bp.Title, bp.AuthorName, bp.Description, bp.YearOfRelease, bp.Tags, imagePath)
}
