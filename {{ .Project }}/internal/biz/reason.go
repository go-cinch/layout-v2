package biz

import (
	"context"
{{- if .Computed.enable_i18n_final }}

	"github.com/go-cinch/common/middleware/i18n"
{{- end }}
{{- if .Computed.enable_reason_proto_final }}
	"{{.Computed.module_name_final}}/api/reason"
{{- else }}
	"github.com/go-kratos/kratos/v2/errors"
{{- end }}
)
{{- if .Computed.enable_i18n_final }}

// Error constants for i18n translation keys
const (
	IdempotentTokenExpired = "idempotent.token.expired"
	TooManyRequests        = "too.many.requests"
	DataNotChange          = "data.not.change"
	DuplicateField         = "duplicate.field"
	RecordNotFound         = "record.not.found"
	InternalError          = "internal.error"
	IllegalParameter       = "illegal.parameter"
)
{{- end }}

var (
{{- if .Computed.enable_i18n_final }}
	{{- if .Computed.enable_reason_proto_final }}
	ErrIdempotentTokenExpired = func(ctx context.Context) error {
		return i18n.NewError(ctx, IdempotentTokenExpired, reason.ErrorIllegalParameter)
	}

	ErrTooManyRequests = func(ctx context.Context) error {
		return reason.ErrorTooManyRequests(i18n.FromContext(ctx).T(TooManyRequests))
	}

	ErrDataNotChange = func(ctx context.Context, args ...string) error {
		return i18n.NewError(ctx, DataNotChange, reason.ErrorIllegalParameter, args...)
	}

	ErrDuplicateField = func(ctx context.Context, args ...string) error {
		return i18n.NewError(ctx, DuplicateField, reason.ErrorIllegalParameter, args...)
	}

	ErrRecordNotFound = func(ctx context.Context, args ...string) error {
		return i18n.NewError(ctx, RecordNotFound, reason.ErrorNotFound, args...)
	}

	ErrInternal = func(ctx context.Context, args ...string) error {
		return i18n.NewError(ctx, InternalError, reason.ErrorInternal, args...)
	}

	ErrIllegalParameter = func(ctx context.Context, args ...string) error {
		return i18n.NewError(ctx, IllegalParameter, reason.ErrorIllegalParameter, args...)
	}
	{{- else }}
	ErrIdempotentTokenExpired = func(ctx context.Context) error {
		return i18n.NewError(ctx, IdempotentTokenExpired, errors.BadRequest("ILLEGAL_PARAMETER", ""))
	}

	ErrTooManyRequests = func(ctx context.Context) error {
		return errors.New(429, "TOO_MANY_REQUESTS", i18n.FromContext(ctx).T(TooManyRequests))
	}

	ErrDataNotChange = func(ctx context.Context, args ...string) error {
		return i18n.NewError(ctx, DataNotChange, errors.BadRequest("ILLEGAL_PARAMETER", ""), args...)
	}

	ErrDuplicateField = func(ctx context.Context, args ...string) error {
		return i18n.NewError(ctx, DuplicateField, errors.BadRequest("ILLEGAL_PARAMETER", ""), args...)
	}

	ErrRecordNotFound = func(ctx context.Context, args ...string) error {
		return i18n.NewError(ctx, RecordNotFound, errors.NotFound("NOT_FOUND", ""), args...)
	}

	ErrInternal = func(ctx context.Context, args ...string) error {
		return i18n.NewError(ctx, InternalError, errors.InternalServer("INTERNAL", ""), args...)
	}

	ErrIllegalParameter = func(ctx context.Context, args ...string) error {
		return i18n.NewError(ctx, IllegalParameter, errors.BadRequest("ILLEGAL_PARAMETER", ""), args...)
	}
	{{- end }}
{{- else }}
	{{- if .Computed.enable_reason_proto_final }}
	ErrIdempotentTokenExpired = func(ctx context.Context) error {
		return reason.ErrorIllegalParameter("idempotent token has expired")
	}

	ErrTooManyRequests = func(ctx context.Context) error {
		return reason.ErrorTooManyRequests("too many requests, please try again later")
	}

	ErrDataNotChange = func(ctx context.Context, args ...string) error {
		msg := "data has not changed"
		if len(args) > 0 {
			msg = args[0]
		}
		return reason.ErrorIllegalParameter(msg)
	}

	ErrDuplicateField = func(ctx context.Context, args ...string) error {
		msg := "duplicate field"
		if len(args) > 0 {
			msg = args[0]
		}
		return reason.ErrorIllegalParameter(msg)
	}

	ErrRecordNotFound = func(ctx context.Context, args ...string) error {
		msg := "record not found"
		if len(args) > 0 {
			msg = args[0]
		}
		return reason.ErrorNotFound(msg)
	}

	ErrInternal = func(ctx context.Context, args ...string) error {
		msg := "internal error"
		if len(args) > 0 {
			msg = args[0]
		}
		return reason.ErrorInternal(msg)
	}

	ErrIllegalParameter = func(ctx context.Context, args ...string) error {
		msg := "illegal parameter"
		if len(args) > 0 {
			msg = args[0]
		}
		return reason.ErrorIllegalParameter(msg)
	}
	{{- else }}
	ErrIdempotentTokenExpired = func(ctx context.Context) error {
		return errors.BadRequest("ILLEGAL_PARAMETER", "idempotent token has expired")
	}

	ErrTooManyRequests = func(ctx context.Context) error {
		return errors.New(429, "TOO_MANY_REQUESTS", "too many requests, please try again later")
	}

	ErrDataNotChange = func(ctx context.Context, args ...string) error {
		msg := "data has not changed"
		if len(args) > 0 {
			msg = args[0]
		}
		return errors.BadRequest("ILLEGAL_PARAMETER", msg)
	}

	ErrDuplicateField = func(ctx context.Context, args ...string) error {
		msg := "duplicate field"
		if len(args) > 0 {
			msg = args[0]
		}
		return errors.BadRequest("ILLEGAL_PARAMETER", msg)
	}

	ErrRecordNotFound = func(ctx context.Context, args ...string) error {
		msg := "record not found"
		if len(args) > 0 {
			msg = args[0]
		}
		return errors.NotFound("NOT_FOUND", msg)
	}

	ErrInternal = func(ctx context.Context, args ...string) error {
		msg := "internal error"
		if len(args) > 0 {
			msg = args[0]
		}
		return errors.InternalServer("INTERNAL", msg)
	}

	ErrIllegalParameter = func(ctx context.Context, args ...string) error {
		msg := "illegal parameter"
		if len(args) > 0 {
			msg = args[0]
		}
		return errors.BadRequest("ILLEGAL_PARAMETER", msg)
	}
	{{- end }}
{{- end }}
)
