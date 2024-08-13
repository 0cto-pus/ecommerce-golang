package handlers

import (
	"ecommerce-golang/internal/api/rest"
	"ecommerce-golang/internal/dto"
	"ecommerce-golang/internal/repository"
	"ecommerce-golang/internal/service"
	"log"
	"net/http"

	"github.com/gofiber/fiber/v3"
)

type UserHandler struct {
	svc service.UserService
	//svc userservice
}

func SetupUserRoutes(rh *rest.RestHandler) {
	app:= rh.App

	//create instance of  user service & inject to handler
	svc := service.UserService{
		Repo: repository.NewUserRepository(rh.DB),
		Auth: rh.Auth,
		Config: rh.Config,
	}
	handler := UserHandler{
		svc: svc,
	}
	pubRoutes := app.Group("/users")
	// Public endpoints
	pubRoutes.Post("/signup", handler.SignUp)
	pubRoutes.Post("/signin", handler.SignIn)


	pvtRoutes := pubRoutes.Group("/", rh.Auth.Authorize)

	//Private endpoints
	pvtRoutes.Get("/verify", handler.GetVerificationCode)
	pvtRoutes.Post("/verify", handler.Verify)

	pvtRoutes.Post("/profile", handler.CreateProfile)
	pvtRoutes.Get("/profile", handler.GetProfile)

	pvtRoutes.Post("/cart", handler.AddToCart)
	pvtRoutes.Get("/cart", handler.GetCart)
	pvtRoutes.Get("/order", handler.GetOrders)
	pvtRoutes.Get("/order/:id", handler.GetOrder)

	pvtRoutes.Post("/become-seller", handler.BecomeSeller)
}

func (h *UserHandler) SignUp(ctx fiber.Ctx) error{

	user:= dto.UserSignUp{}
	err:= ctx.Bind().Body(&user)

	if err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{
			"message": "please provide valid input",
		})
	}

	token, err:= h.svc.SignUp(user)

	if err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{
			"message": "error on sign up",
		})
	}


	return ctx.Status(http.StatusOK).JSON(fiber.Map{
		"message": "register",
		"token" : token,
	})
}

func (h *UserHandler) SignIn(ctx fiber.Ctx) error{
	signInInput:= dto.UserSignIn{}
	err:= ctx.Bind().Body(&signInInput)

	if err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{
			"message": "please provide valid input",
		})
	}

	token, err :=h.svc.SignIn(signInInput.Email,signInInput.Password)
	if err != nil {
		return ctx.Status(http.StatusUnauthorized).JSON(fiber.Map{
			"message": "please provide correct user id password",
		})
	}

	return ctx.Status(http.StatusOK).JSON(fiber.Map{
		"message": "login",
		"token" : token,
	})
	
}

func (h *UserHandler) GetVerificationCode(ctx fiber.Ctx) error{

	user := h.svc.Auth.GetCurrentUser(ctx)
	log.Println(user)
	// create verification code and update to user profile in DB
	err := h.svc.GetVerificationCode(user)
	log.Println(err)

	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(&fiber.Map{
			"message": "unable to generate verification code",
		})
	}

	return ctx.Status(http.StatusOK).JSON(&fiber.Map{
		"message": "get verification code",
	})
}

func (h *UserHandler) Verify(ctx fiber.Ctx) error{
	user := h.svc.Auth.GetCurrentUser(ctx)

	// request
	var req dto.VerificationCodeInput

	if err := ctx.Bind().Body(&req); err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"message": "please provide a valid input",
		})
	}

	err := h.svc.VerifyCode(user.ID, req.Code)

	if err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"message": err.Error(),
		})
	}

	return ctx.Status(http.StatusOK).JSON(&fiber.Map{
		"message": "verified successfully",
	})
}
func (h *UserHandler) CreateProfile(ctx fiber.Ctx) error{
	return ctx.Status(http.StatusOK).JSON(fiber.Map{
		"message": "Sign In",
	})
}
func (h *UserHandler) GetProfile(ctx fiber.Ctx) error{
	user := h.svc.Auth.GetCurrentUser(ctx)
	//log.Println(user)
	return ctx.Status(http.StatusOK).JSON(&fiber.Map{
		"message": "get profile",
		"user":    user,
	})
}
func (h *UserHandler) AddToCart(ctx fiber.Ctx) error{
	return ctx.Status(http.StatusOK).JSON(fiber.Map{
		"message": "Sign In",
	})
}

func (h *UserHandler) GetCart(ctx fiber.Ctx) error{
	return ctx.Status(http.StatusOK).JSON(fiber.Map{
		"message": "Sign In",
	})
}

func (h *UserHandler) GetOrders(ctx fiber.Ctx) error{
	return ctx.Status(http.StatusOK).JSON(fiber.Map{
		"message": "Sign In",
	})
}

func (h *UserHandler) GetOrder(ctx fiber.Ctx) error{
	return ctx.Status(http.StatusOK).JSON(fiber.Map{
		"message": "Sign In",
	})
}
func (h *UserHandler) BecomeSeller(ctx fiber.Ctx) error{
	user := h.svc.Auth.GetCurrentUser(ctx)

	req := dto.SellerInput{}
	err := ctx.Bind().Body(&req)
	if err != nil {
		return ctx.Status(400).JSON(&fiber.Map{
			"message": "request parameters are not valid",
		})
	}

	token, err := h.svc.BecomeSeller(user.ID, req)

	if err != nil {
		return ctx.Status(http.StatusUnauthorized).JSON(&fiber.Map{
			"message": "fail to become seller",
		})
	}

	return ctx.Status(http.StatusOK).JSON(&fiber.Map{
		"message": "become seller",
		"token":   token,
	})
}