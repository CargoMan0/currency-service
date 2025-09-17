package handler

type registerRequest struct {
	Username string `form:"username" binding:"required"`
	Password string `form:"password" binding:"required"`
}

type loginRequest struct {
	Username string `form:"username" binding:"required"`
	Password string `form:"password" binding:"required"`
}

type currencyRequest struct {
	Currency string `form:"currency" binding:"required"`
	DateFrom string `form:"date_from" binding:"required,datetime=2006-01-02"`
	DateTo   string `form:"date_to" binding:"required,datetime=2006-01-02"`
}
