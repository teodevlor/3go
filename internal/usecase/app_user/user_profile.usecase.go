package appuser

import (
	"context"
	"errors"
	"fmt"

	dto "go-structure/internal/dto/app_user"
	"go-structure/internal/repository"
	appuserrepo "go-structure/internal/repository/app_user"
	"go-structure/internal/repository/model"
	"go-structure/internal/usecase"
	"go-structure/pkg/validator"
)

type (
	IUserProfileUsecase interface {
		RegisterUserProfile(ctx context.Context, req *dto.UserRegisterRequestDto) (*dto.UserRegisterResponseDto, error)
	}

	userProfileUsecase struct {
		userProfileRepository appuserrepo.IUserProfileRepository
		accountRepository     repository.IAccountRepository
		notifyUsecase         usecase.INotifyUsecase
	}
)

var (
	ErrUserAlreadyRegistered = errors.New("user already registered")
)

func NewUserProfileUsecase(
	userProfileRepo appuserrepo.IUserProfileRepository,
	accountRepo repository.IAccountRepository,
	notifyUc usecase.INotifyUsecase,
) IUserProfileUsecase {
	return &userProfileUsecase{
		userProfileRepository: userProfileRepo,
		accountRepository:     accountRepo,
		notifyUsecase:         notifyUc,
	}
}

func (u *userProfileUsecase) RegisterUserProfile(ctx context.Context, req *dto.UserRegisterRequestDto) (*dto.UserRegisterResponseDto, error) {
	hashedPassword, err := validator.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	// 1. Kiểm tra account theo phone
	account, err := u.accountRepository.GetByPhone(ctx, req.Phone)
	if err != nil {
		return nil, err
	}

	// 2. Nếu chưa có account thì tạo mới
	if account == nil {
		account = &model.Account{
			Phone:        req.Phone,
			PasswordHash: hashedPassword,
		}
		account, err = u.accountRepository.CreateAccount(ctx, account)
		if err != nil {
			return nil, err
		}

		// Account vừa tạo thì chắc chắn chưa có profile, nên nhảy xuống tạo profile mới.
	} else {
		// Nếu đã có account thì kiểm tra user profile hiện có
		existingProfile, err := u.userProfileRepository.GetByAccountID(ctx, account.ID)
		if err != nil {
			return nil, err
		}

		// Nếu account đã có user profile thì không cho đăng ký nữa
		if existingProfile != nil {
			return nil, ErrUserAlreadyRegistered
		}
	}

	// 3. Nếu đến đây nghĩa là: hoặc account mới tạo, hoặc account cũ nhưng chưa có profile → tạo profile mới
	userProfile := &model.UserProfile{
		AccountID: account.ID,
		FullName:  req.FullName,
		IsActive:  true,
	}

	saved, err := u.userProfileRepository.RegisterUserProfile(ctx, userProfile)
	if err != nil {
		return nil, err
	}

	// 4. Gửi OTP qua kênh notify (Telegram adapter hiện tại)
	if u.notifyUsecase != nil {
		otp := "123456"
		msg := fmt.Sprintf("OTP đăng ký tài khoản cho số %s là: %s", req.Phone, otp)
		_ = u.notifyUsecase.SendMessage(ctx, msg)
	}

	return &dto.UserRegisterResponseDto{
		ID:        saved.ID,
		Phone:     req.Phone,
		FullName:  saved.FullName,
		CreatedAt: saved.CreatedAt,
		UpdatedAt: saved.UpdatedAt,
	}, nil
}
