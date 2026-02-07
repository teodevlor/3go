package appuser

import (
	dto "go-structure/internal/dto/app_user"
	"go-structure/internal/repository/model"
)

func ToUserProfileResponse(account *model.Account, userProfile *model.UserProfile) dto.UserProfileResponseDto {
	if account == nil || userProfile == nil {
		return dto.UserProfileResponseDto{}
	}
	return dto.UserProfileResponseDto{
		ID:        account.ID,
		FullName:  userProfile.FullName,
		AvatarURL: userProfile.AvatarURL,
		IsActive:  userProfile.IsActive,
		Phone:     account.Phone,
		Email:     account.Email,
	}
}

func ToLoginResponseDto(accessToken, refreshToken string, account *model.Account, userProfile *model.UserProfile) dto.UserLoginResponseDto {
	return dto.UserLoginResponseDto{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		UserProfile:  ToUserProfileResponse(account, userProfile),
	}
}

func ToRefreshTokenResponseDto(accessToken, refreshToken string) dto.RefreshTokenResponseDto {
	return dto.RefreshTokenResponseDto{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}
}

func ToUpdateUserProfileResponseDto(userMessage string, account *model.Account, profile *model.UserProfile) dto.UpdateUserProfileResponseDto {
	return dto.UpdateUserProfileResponseDto{
		UserMessage: userMessage,
		UserProfile: ToUserProfileResponse(account, profile),
	}
}
