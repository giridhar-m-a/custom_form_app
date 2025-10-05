package queryparams

type GoogleAuthQuery struct {
    Code string `form:"code" binding:"required"`
}