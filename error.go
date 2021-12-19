package response

import (
	"github.com/go-openapi/errors"
	"github.com/go-openapi/swag"
)

// Problem defines the Problem JSON type defined by RFC 7807 - media type
// application/problem+json.
// It should be the expected error response for all APIs.
type Problem struct {
	// Type     string      `json:"type"`
	Title    string      `json:"title"`
	Status   int         `json:"status"`
	Detail   interface{} `json:"detail"`
	Instance string      `json:"instance"`
}

type HttpValidationError struct {
	// Error describing field validation failure
	// Required: true
	Error *string `json:"error"`

	// Indicates how the invalid field was provided
	// Required: true
	In *string `json:"in"`

	// Key of field failing validation
	// Required: true
	Key *string `json:"key"`
}

// MarshalBinary interface implementation
func (m *HttpValidationError) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *HttpValidationError) UnmarshalBinary(b []byte) error {
	var res HttpValidationError
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}

func FormatErrors(err *errors.CompositeError) []*HttpValidationError {
	valErrs := make([]*HttpValidationError, 0, len(err.Errors))
	for _, e := range err.Errors {
		switch ee := e.(type) {
		case *errors.Validation:
			valErrs = append(valErrs, &HttpValidationError{
				Key:   &ee.Name,
				In:    &ee.In,
				Error: swag.String(ee.Error()),
			})
		case *errors.CompositeError:
			valErrs = append(valErrs, FormatErrors(ee)...)
		default:
			return nil
		}
	}

	return valErrs
}
