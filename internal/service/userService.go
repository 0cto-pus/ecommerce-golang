package service

import (
	"ecommerce-golang/config"
	"ecommerce-golang/internal/domain"
	"ecommerce-golang/internal/dto"
	"ecommerce-golang/internal/helper"
	"ecommerce-golang/internal/repository"
	"ecommerce-golang/pkg/notification"
	"errors"
	"fmt"
	"log"
	"time"
)

type UserService struct {
	Config config.AppConfig
	Repo repository.UserRepository
	Auth helper.Auth
}
func (s UserService) findUserByEmail(email string) (*domain.User, error) {
	//perform some db operation
	//business logic
	user, err := s.Repo.FindUser(email)

	return &user, err
}

func (s UserService) SignUp(input dto.UserSignUp)(string, error){
	hPassword, err := s.Auth.CreateHashedPassword(input.Password)

	 if err != nil {
		return "", err
	}
 
	user, err := s.Repo.CreateUser(domain.User{
		Email:    input.Email,
		Password: hPassword,
		Phone:    input.Phone,
	})
	

	 if err != nil {
		return "", err
	}

	return s.Auth.GenerateToken(user.ID, user.Email, user.UserType)
}

func (s UserService) SignIn(email string, password string)(string, error){

	user,err := s.findUserByEmail(email)

	if err != nil {
		return "", errors.New("user does not exist with the provided email id");
	}

	err = s.Auth.VerifyPassword(password, user.Password)

	if err != nil {
		return "", err
	}

	// generate token
	return s.Auth.GenerateToken(user.ID, user.Email, user.UserType)
}

func (s UserService) isVerifiedUser(id uint) bool {

	currentUser, err := s.Repo.FindUserById(id)
	fmt.Print(currentUser.Verified)
	return err == nil && currentUser.Verified
}

func (s UserService) GetVerificationCode(e domain.User) error{
		// if user already verified
		if s.isVerifiedUser(e.ID) {
			return errors.New("user already verified")
		}
	
		// generate verification code
		code, err := s.Auth.GenerateCode()
		if err != nil {
			return err
		}
	
		// update user
		user := domain.User{
			Expiry: time.Now().Add(30 * time.Minute),
			Code:   code,
		}
	
		_, err = s.Repo.UpdateUser(e.ID, user)
	
		if err != nil {
			return errors.New("unable to update verification code")
		}
	
		user, _ = s.Repo.FindUserById(e.ID)
	
		// send SMS
		notificationClient := notification.NewNotificationClient(s.Config)
	
		msg := fmt.Sprintf("Your verification code is %v", code)
	
		err = notificationClient.SendSMS(user.Phone, msg)
		if err != nil {
			return errors.New("error on sending sms")
		}
	
		// return verification code
		return nil
}


func (s UserService) VerifyCode(id uint, code int) error{
	// if user already verified
	if s.isVerifiedUser(id) {
		log.Println("verified...")
		return errors.New("user already verified")
	}

	user, err := s.Repo.FindUserById(id)

	if err != nil {
		return err
	}

	if user.Code != code {
		return errors.New("verification code does not match")
	}

	if !time.Now().Before(user.Expiry) {
		return errors.New("verification code expired")
	}

	updateUser := domain.User{
		Verified: true,
	}

	_, err = s.Repo.UpdateUser(id, updateUser)

	if err != nil {
		return errors.New("unable to verify user")
	}

	return nil
}

func (s UserService) CreateProfile(id uint, input any)error{
	return nil
}

func (s UserService) getProfile(id uint)(*domain.User, error){
	return nil, nil
}

func (s UserService) updateProfile(id uint,input any) error{
	return nil
}
func (s UserService) BecomeSeller(id uint, input dto.SellerInput)(string, error){
	user, _ := s.Repo.FindUserById(id)

	if user.UserType == domain.SELLER {
		return "", errors.New("you have already joined seller program")
	}

	// update user
	seller, err := s.Repo.UpdateUser(id, domain.User{
		FirstName: input.FirstName,
		LastName:  input.LastName,
		Phone:     input.PhoneNumber,
		UserType:  domain.SELLER,
	})

	if err != nil {
		return "", err
	}

	// generating token
	token, err := s.Auth.GenerateToken(user.ID, user.Email, seller.UserType)

	// create bank account information

	err = s.Repo.CreateBankAccount(domain.BankAccount{
		BankAccount: input.BankAccountNumber,
		SwiftCode:   input.SwiftCode,
		PaymentType: input.PaymentType,
		UserId:      id,
	})

	return token, err
}

func (s UserService) findCart(id uint, input any)([]interface{}, error){
	return nil, nil
}

func (s UserService) createCart(input any, u domain.User)([]interface{}, error){
	return nil, nil
}

func (s UserService) createOrder(u domain.User)(int, error){
	return 0, nil
}

func (s UserService) getOrders(u domain.User)([]interface{}, error){
	return nil, nil
}

func (s UserService) getOrderById(id uint, uId int)(interface{}, error){
	return nil, nil
}
