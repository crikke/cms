package content

// import (
// 	"fmt"
// 	"time"

// 	"github.com/google/uuid"
// 	"golang.org/x/text/language"
// )

// type Content struct {
// 	// A node are required to have a localized URLSegment for each configured locale.
// 	ID         ContentReference
// 	ParentID   uuid.UUID
// 	URLSegment string
// 	Name       string
// 	Properties []Property
// 	Created    time.Time
// 	Updated    time.Time
// }

// // A contentreference is a reference to a piece of content that is versioned & has a locale
// type ContentReference struct {
// 	ID      uuid.UUID
// 	Version *int
// 	Locale  *language.Tag
// }

// func (c ContentReference) String() string {

// 	str := c.ID.String()

// 	if c.Version != nil {
// 		str = fmt.Sprintf("%s-%d", str, *c.Version)
// 	}

// 	if c.Locale != nil {
// 		str = fmt.Sprintf("%s-%s", str, c.Locale.String())
// 	}

// 	return str
// }

// type Property struct {
// 	ID        uuid.UUID
// 	Name      string
// 	Type      string
// 	Localized bool
// 	Value     interface{}
// }
