package fact

import (
	commonerrors "github.com/cafo13/animal-facts/backend/common/errors"
	"github.com/pkg/errors"
)

var ErrTextTooLong = commonerrors.NewIncorrectInputError("Text too long", "text-too-long")
var ErrSourceTooLong = commonerrors.NewIncorrectInputError("Source too long", "source-too-long")

type Fact struct {
	uuid string

	text   string
	source string
}

func NewFact(uuid string, text string, source string) (*Fact, error) {
	if uuid == "" {
		return nil, errors.New("empty fact uuid")
	}
	if text == "" {
		return nil, errors.New("empty fact text")
	}
	if source == "" {
		return nil, errors.New("empty fact source")
	}
	if len(text) > 1000 {
		return nil, ErrTextTooLong
	}
	if len(source) > 1000 {
		return nil, ErrSourceTooLong
	}

	return &Fact{
		uuid:   uuid,
		text:   text,
		source: source,
	}, nil
}

func (f Fact) UUID() string {
	return f.uuid
}

func (f Fact) Text() string {
	return f.text
}

func (f Fact) Source() string {
	return f.source
}

func (f *Fact) UpdateText(text string) error {
	if len(text) > 1000 {
		return errors.WithStack(ErrTextTooLong)
	}

	f.text = text
	return nil
}

func (f *Fact) UpdateSource(source string) error {
	if len(source) > 1000 {
		return errors.WithStack(ErrTextTooLong)
	}

	f.source = source
	return nil
}
