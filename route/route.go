package route

import (
	"net/http"

	"github.com/Safayet-Shawn/Parking-Lot/api"
	cfg "github.com/Safayet-Shawn/Parking-Lot/config"
	"github.com/go-chi/chi/v5"
)

func NewRouter() http.Handler {
	r := chi.NewRouter()
	SetHashToken(cfg.JwtSecret())
	r.Use(JwtAuthentication)
	r.Get("/", api.ApiHome)
	r.Post("/v1/register", api.RegistrationUser)                  // checked// auth done
	r.Post("/v1/login", api.UserLogin)                            // checked// auth done
	r.Post("/v1/create-parking", api.CreateParking)               // checked// auth done
	r.Post("/v1/create-parking/{id}/slot", api.CreateParkingSlot) // checked //auth done
	r.Get("/v1/get-slot/{latlon}", api.GetSlot)                   // checked [user]// no auth need
	r.Put("/v1/update-slot/{sid}", api.UpdateSlot)                //  cheched//auth

	r.Post("/v1/parking-lot/{pid}/slot/{sid}", api.Park)  // // checked [user]
	r.Put("/v1/parking-lot/{pid}/slot/{sid}", api.Unpark) // checked
	r.Get("/v1/parking-slot", api.GetSlotInfo)
	return r
}
