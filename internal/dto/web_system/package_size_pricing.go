package web_system

type CreatePackageSizePricingRequestDto struct {
	ServiceID   string  `json:"service_id" binding:"required,uuid"`
	PackageSize string  `json:"package_size" binding:"required,max=100"`
	ExtraPrice  float64 `json:"extra_price" binding:"required,gte=0"`
	IsActive    bool    `json:"is_active"`
}

type UpdatePackageSizePricingRequestDto struct {
	ServiceID   string  `json:"service_id" binding:"required,uuid"`
	PackageSize string  `json:"package_size" binding:"required,max=100"`
	ExtraPrice  float64 `json:"extra_price" binding:"required,gte=0"`
	IsActive    bool    `json:"is_active"`
}

type PackageSizePricingItemDto struct {
	ID          string  `json:"id"`
	ServiceID   string  `json:"service_id"`
	PackageSize string  `json:"package_size"`
	ExtraPrice  float64 `json:"extra_price"`
	IsActive    bool    `json:"is_active"`
}
