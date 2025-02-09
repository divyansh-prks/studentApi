package student

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"

	"github.com/divyansh/students-api/internal/storage"
	"github.com/divyansh/students-api/internal/types"
	"github.com/divyansh/students-api/internal/utils/response"
	"github.com/go-playground/validator"
)


func New(storage storage.Storage) http.HandlerFunc {
	slog.Info("creating a student ")
	return func (w http.ResponseWriter , r *http.Request){
		var student types.Student
		err := json.NewDecoder(r.Body).Decode(&student)

		if errors.Is(err , io.EOF){
			response.WriteJson(w , http.StatusBadRequest , response.GeneralError(err))
			return

		}

		//validation 
		if err := validator.New().Struct(student); err != nil{
			validateErrs := err.(validator.ValidationErrors)
			response.WriteJson(w , http.StatusBadGateway , response.ValidationError(validateErrs))
		}

		lastid , err := storage.CreateStudent(
			student.Name,
			student.Email,
			student.Age,

		)

		if err != nil {
			response.WriteJson(w , http.StatusInternalServerError, err)
			return 
		}

		slog.Info("user created succefully" , slog.String("userId" , fmt.Sprint(lastid)))


		response.WriteJson(w , http.StatusCreated , map[string]int64 {"id" : lastid} )

		w.Write([]byte("welcome to students api"))
	}
}