package ginx

import (
	"gin-api-server/internal/extend/svcerrorsx"
	"gin-api-server/pkg/svcerrors"
	"net/http"
	"strconv"
	"strings"
	"unicode"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

const (
	userIDKey = "user-id"
)

// ParseParamID Parse path id
func ParseParamID(c *gin.Context, key string) uint64 {
	val := c.Param(key)
	id, err := strconv.ParseUint(val, 10, 64)
	if err != nil {
		return 0
	}
	return id
}

func GetUserID(c *gin.Context) uint64 {
	return c.GetUint64(userIDKey)
}

func SetUserID(c *gin.Context, userID uint64) {
	c.Set(userIDKey, userID)
}

func RespondData(c *gin.Context, code int, data interface{}) {
	ok, _ := strconv.ParseBool(c.Query("pretty"))
	if ok {
		c.IndentedJSON(code, data)
	} else {
		c.JSON(code, data)
	}
}

func RespondMsg(c *gin.Context, msg interface{}) {
	c.JSON(http.StatusOK, gin.H{"message": msg})
}

func RespondErr(c *gin.Context, err error) {
	var resErr *svcerrors.SvcError
	if e, ok := err.(*svcerrors.SvcError); ok {
		resErr = e
	} else if errs, ok := err.(validator.ValidationErrors); ok {
		fieldErrs := make([]error, len(errs))
		for i, v := range errs {
			key := v.Namespace()
			// 在 UserForm.name 这种情况下，需要去除 UserForm.
			if unicode.IsUpper([]rune(key)[0]) {
				keyParts := strings.SplitN(key, ".", 2)
				key = keyParts[1]
			}
			fieldErrs[i] = &svcerrors.SvcErrorItem{
				Key: key,
				Err: v.Error(),
			}
		}
		resErr = svcerrorsx.ErrInvalidParameter.WithErrors(fieldErrs...)
	} else {
		resErr = svcerrorsx.ErrServerError.WithErrors(&svcerrors.SvcErrorItem{Key: "0", Err: err.Error()})
	}
	c.JSON(resErr.Status, resErr)
}
