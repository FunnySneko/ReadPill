package api

import (
	"encoding/json"
	"net/http"
)

func (s *ServerHandler) GetBooks(w http.ResponseWriter, r *http.Request) {
	books, err := s.Db.GetBooks(r.Context(), nil, nil)
	if err != nil {
		ErrorOut(w, err)
	}

	b := BookArray{}
	for _, book := range books {
		b = append(b, mapBookToAPI(book))
	}

	response, err := json.Marshal(b)
	if err != nil {
		ErrorOut(w, err)
	}

	w.Write(response)
}

func (s *ServerHandler) PostBooks(w http.ResponseWriter, r *http.Request) {
	userId := int(r.Context().Value("userID").(float64))

	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		ErrorOut(w, err)
		return
	}

	bp := BookPost{}
	data := r.FormValue("data")
	err = json.Unmarshal([]byte(data), &bp)
	if err != nil {
		ErrorOut(w, err)
		return
	}

	file, header, err := r.FormFile("image")
	if err != nil {
		ErrorOut(w, err)
		return
	}

	imagePath, err := s.Ah.UploadImage(file, *header)
	if err != nil {
		ErrorOut(w, err)
		return
	}

	s.Ah.PostBook(r.Context(), bp.Title, bp.AuthorName, bp.Description, bp.YearOfRelease, bp.Tags, imagePath, userId)
}
