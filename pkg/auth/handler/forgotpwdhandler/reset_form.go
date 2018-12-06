package forgotpwdhandler

import (
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/skygeario/skygear-server/pkg/auth"
	"github.com/skygeario/skygear-server/pkg/auth/dependency"
	"github.com/skygeario/skygear-server/pkg/auth/dependency/forgotpwdemail"
	"github.com/skygeario/skygear-server/pkg/auth/dependency/provider/password"
	"github.com/skygeario/skygear-server/pkg/auth/dependency/userprofile"
	"github.com/skygeario/skygear-server/pkg/core/auth/authinfo"
	"github.com/skygeario/skygear-server/pkg/core/auth/authtoken"
	"github.com/skygeario/skygear-server/pkg/core/auth/authz"
	"github.com/skygeario/skygear-server/pkg/core/auth/authz/policy"
	"github.com/skygeario/skygear-server/pkg/core/db"
	"github.com/skygeario/skygear-server/pkg/core/inject"
	"github.com/skygeario/skygear-server/pkg/server/skyerr"
)

// ForgotPasswordResetFormHandlerFactory creates ForgotPasswordResetFormHandler
type ForgotPasswordResetFormHandlerFactory struct {
	Dependency auth.RequestDependencyMap
}

// NewHandler creates new ForgotPasswordResetFormHandler
func (f ForgotPasswordResetFormHandlerFactory) NewHandler(request *http.Request) http.Handler {
	h := &ForgotPasswordResetFormHandler{}
	inject.DefaultRequestInject(h, f.Dependency, request)
	return h
}

// ProvideAuthzPolicy provides authorization policy of handler
func (f ForgotPasswordResetFormHandlerFactory) ProvideAuthzPolicy() authz.Policy {
	return policy.Everybody{Allow: true}
}

type ForgotPasswordResetFormPayload struct {
	UserID          string
	Code            string
	ExpireAt        int64
	ExpireAtTime    time.Time
	NewPassword     string
	ConfirmPassword string
}

func decodeForgotPasswordResetFormRequest(request *http.Request) (payload ForgotPasswordResetFormPayload, err error) {
	if err = request.ParseForm(); err != nil {
		return
	}

	p := ForgotPasswordResetFormPayload{}
	p.UserID = request.Form.Get("user_id")
	p.Code = request.Form.Get("code")
	p.NewPassword = request.Form.Get("password")
	p.ConfirmPassword = request.Form.Get("confirm")

	expireAtStr := request.Form.Get("expire_at")
	var expireAt int64
	if expireAtStr != "" {
		if expireAt, err = strconv.ParseInt(expireAtStr, 10, 64); err != nil {
			return
		}
	}

	p.ExpireAt = expireAt
	p.ExpireAtTime = time.Unix(expireAt, 0).UTC()

	payload = p
	return
}

func (payload *ForgotPasswordResetFormPayload) Validate() error {
	if payload.UserID == "" {
		return skyerr.NewInvalidArgument("empty user_id", []string{"user_id"})
	}

	if payload.Code == "" {
		return skyerr.NewInvalidArgument("empty code", []string{"code"})
	}

	if payload.ExpireAt == 0 {
		return skyerr.NewInvalidArgument("empty expire_at", []string{"expire_at"})
	}

	return nil
}

// ForgotPasswordResetFormHandler reset user password with given code from email.
type ForgotPasswordResetFormHandler struct {
	CodeGenerator             *forgotpwdemail.CodeGenerator             `dependency:"ForgotPasswordCodeGenerator"`
	PasswordChecker           dependency.PasswordChecker                `dependency:"PasswordChecker"`
	TokenStore                authtoken.Store                           `dependency:"TokenStore"`
	AuthInfoStore             authinfo.Store                            `dependency:"AuthInfoStore"`
	PasswordAuthProvider      password.Provider                         `dependency:"PasswordAuthProvider"`
	UserProfileStore          userprofile.Store                         `dependency:"UserProfileStore"`
	ResetPasswordHTMLProvider *forgotpwdemail.ResetPasswordHTMLProvider `dependency:"ResetPasswordHTMLProvider"`
	TxContext                 db.TxContext                              `dependency:"TxContext"`
	Logger                    *logrus.Entry                             `dependency:"HandlerLogger"`
}

type resultTemplateContext struct {
	err         skyerr.Error
	payload     ForgotPasswordResetFormPayload
	userProfile userprofile.UserProfile
}

func (h ForgotPasswordResetFormHandler) WithTx() bool {
	return true
}

func (h ForgotPasswordResetFormHandler) prepareResultTemplateContext(r *http.Request) (ctx resultTemplateContext, err error) {
	var payload ForgotPasswordResetFormPayload
	payload, err = decodeForgotPasswordResetFormRequest(r)
	if err != nil {
		return
	}

	if err = payload.Validate(); err != nil {
		return
	}

	ctx.payload = payload

	// Get Profile
	if ctx.userProfile, err = h.UserProfileStore.GetUserProfile(payload.UserID); err != nil {
		h.Logger.WithFields(map[string]interface{}{
			"user_id": payload.UserID,
		}).WithError(err).Error("unable to get user profile")
		err = genericResetPasswordError()
		return
	}

	return
}

// HandleRequestError handle the case when the given data in the form is wrong, e.g. code, user_id, expire_at
func (h ForgotPasswordResetFormHandler) HandleRequestError(rw http.ResponseWriter, err skyerr.Error) {
	context := map[string]interface{}{
		"error": err.Message(),
	}

	url := h.ResetPasswordHTMLProvider.ErrorRedirect(context)
	if url != nil {
		rw.Header().Set("Location", url.String())
		rw.WriteHeader(http.StatusFound)
		return
	}

	html, htmlErr := h.ResetPasswordHTMLProvider.ErrorHTML(context)
	if htmlErr != nil {
		panic(htmlErr)
	}

	rw.WriteHeader(http.StatusBadRequest)
	io.WriteString(rw, html)
}

// HandleResetError handle the case when the user input data in the form is wrong, e.g. password, confirm
func (h ForgotPasswordResetFormHandler) HandleResetError(rw http.ResponseWriter, templateCtx resultTemplateContext) {
	context := map[string]interface{}{
		"error":     templateCtx.err.Message(),
		"code":      templateCtx.payload.Code,
		"user_id":   templateCtx.payload.UserID,
		"expire_at": strconv.FormatInt(templateCtx.payload.ExpireAt, 10),
	}

	url := h.ResetPasswordHTMLProvider.ErrorRedirect(context)
	if url != nil {
		rw.Header().Set("Location", url.String())
		rw.WriteHeader(http.StatusFound)
		return
	}

	context["user"] = templateCtx.userProfile.ToMap()

	// render the form again for failed post request
	html, htmlErr := h.ResetPasswordHTMLProvider.FormHTML(context)
	if htmlErr != nil {
		panic(htmlErr)
	}

	rw.WriteHeader(http.StatusOK)
	io.WriteString(rw, html)
}

func (h ForgotPasswordResetFormHandler) HandleGetForm(rw http.ResponseWriter, templateCtx resultTemplateContext) {
	context := map[string]interface{}{
		"code":      templateCtx.payload.Code,
		"user_id":   templateCtx.payload.UserID,
		"user":      templateCtx.userProfile,
		"expire_at": strconv.FormatInt(templateCtx.payload.ExpireAt, 10),
	}

	html, htmlErr := h.ResetPasswordHTMLProvider.FormHTML(context)
	if htmlErr != nil {
		panic(htmlErr)
	}

	rw.WriteHeader(http.StatusOK)
	io.WriteString(rw, html)
}

func (h ForgotPasswordResetFormHandler) HandleResetSuccess(rw http.ResponseWriter, templateCtx resultTemplateContext) {
	context := map[string]interface{}{
		"code":      templateCtx.payload.Code,
		"user_id":   templateCtx.payload.UserID,
		"expire_at": strconv.FormatInt(templateCtx.payload.ExpireAt, 10),
	}

	url := h.ResetPasswordHTMLProvider.SuccessRedirect(context)
	if url != nil {
		rw.WriteHeader(http.StatusFound)
		rw.Header().Set("Location", url.String())
		return
	}

	context["user"] = templateCtx.userProfile.ToMap()

	html, htmlErr := h.ResetPasswordHTMLProvider.SuccessHTML(context)
	if htmlErr != nil {
		panic(htmlErr)
	}

	rw.WriteHeader(http.StatusOK)
	io.WriteString(rw, html)
}

func (h ForgotPasswordResetFormHandler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	if err := h.TxContext.BeginTx(); err != nil {
		h.HandleRequestError(rw, skyerr.MakeError(err))
		return
	}

	var err error
	defer func() {
		if err != nil {
			h.TxContext.RollbackTx()
		} else {
			h.TxContext.CommitTx()
		}
	}()

	var templateCtx resultTemplateContext
	if templateCtx, err = h.prepareResultTemplateContext(r); err != nil {
		h.HandleRequestError(rw, skyerr.MakeError(err))
		return
	}

	// check code expiration
	if timeNow().After(templateCtx.payload.ExpireAtTime) {
		h.Logger.Error("forgot password code expired")
		h.HandleRequestError(rw, genericResetPasswordError())
		return
	}

	authInfo := authinfo.AuthInfo{}
	if e := h.AuthInfoStore.GetAuth(templateCtx.payload.UserID, &authInfo); e != nil {
		h.Logger.WithFields(map[string]interface{}{
			"user_id": templateCtx.payload.UserID,
		}).WithError(e).Error("user not found")
		h.HandleRequestError(rw, genericResetPasswordError())
		return
	}

	// Get password auth principals
	principals, err := h.PasswordAuthProvider.GetPrincipalsByUserID(authInfo.ID)
	if err != nil {
		h.Logger.WithFields(map[string]interface{}{
			"user_id": templateCtx.payload.UserID,
		}).WithError(err).Error("unable to get password auth principals")
		h.HandleRequestError(rw, genericResetPasswordError())
		return
	}

	hashedPassword := principals[0].HashedPassword
	expectedCode := h.CodeGenerator.Generate(authInfo, templateCtx.userProfile, hashedPassword, templateCtx.payload.ExpireAtTime)
	if templateCtx.payload.Code != expectedCode {
		h.Logger.WithFields(map[string]interface{}{
			"user_id":       templateCtx.payload.UserID,
			"code":          templateCtx.payload.Code,
			"expected_code": expectedCode,
		}).Error("wrong forgot password reset password code")
		h.HandleRequestError(rw, genericResetPasswordError())
		return
	}

	if r.Method == http.MethodGet {
		h.HandleGetForm(rw, templateCtx)
		return
	}

	h.resetPassword(rw, templateCtx, principals)
}

func (h ForgotPasswordResetFormHandler) resetPassword(rw http.ResponseWriter, templateCtx resultTemplateContext, principals []*password.Principal) {
	var err error
	defer func() {
		if err != nil {
			templateCtx.err = skyerr.MakeError(err)
			h.HandleResetError(rw, templateCtx)
		} else {
			h.HandleResetSuccess(rw, templateCtx)
		}
	}()

	resetPwdCtx := password.ResetPasswordRequestContext{
		PasswordChecker:      h.PasswordChecker,
		PasswordAuthProvider: h.PasswordAuthProvider,
	}

	if templateCtx.payload.NewPassword == "" {
		err = skyerr.NewInvalidArgument("empty password", []string{"password"})
		return
	}

	if templateCtx.payload.NewPassword != templateCtx.payload.ConfirmPassword {
		err = skyerr.NewInvalidArgument("confirm password does not match the password", []string{"password", "confirm"})
		return
	}

	err = resetPwdCtx.ExecuteWithPrincipals(templateCtx.payload.NewPassword, principals)
}